import os
import re
import inflect

p = inflect.engine()

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\repository"

for filename in os.listdir(directory):
    if not filename.endswith(".go"):
        continue
    if filename in ["user_repository.go", "patient_repository.go"]:
        continue
        
    filepath = os.path.join(directory, filename)
    with open(filepath, 'r') as f:
        content = f.read()

    match = re.search(r'type\s+([A-Za-z]+)Repository\s+interface', content)
    if not match:
        continue
    
    model_name = match.group(1)
    
    # Simple pluralization
    if model_name.endswith('y'):
        collection_name = model_name[:-1].lower() + 'ies'
    elif model_name.endswith('s'):
        collection_name = model_name.lower()
    else:
        # e.g., MedicalRecord -> medicalrecords? actually firestore conventionally is just plural
        collection_name = p.plural(model_name).lower()
        
    struct_name = model_name[0].lower() + model_name[1:] + "Repository"

    # We will just write a standard CRUD template and replace the whole file.
    # We must scan the interface first to see what methods exist.
    interface_match = re.search(r'type\s+' + model_name + r'Repository\s+interface\s+\{([^}]+)\}', content)
    if not interface_match:
        continue
        
    interface_body = interface_match.group(1)
    
    methods = []
    for line in interface_body.split('\n'):
        line = line.strip()
        if not line: continue
        methods.append(line)

    new_content = f"""package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type {model_name}Repository interface {{
"""
    for method in methods:
        # Convert uint to string
        method = re.sub(r'\buint\b', 'string', method)
        new_content += f"\t{method}\n"
        
    new_content += f"""}}

type {struct_name} struct {{
	db *firestore.Client
}}

func New{model_name}Repository(db *firestore.Client) {model_name}Repository {{
	return &{struct_name}{{db}}
}}
"""
    # Now generate methods based on what's in the interface
    for method in methods:
        method = re.sub(r'\buint\b', 'string', method)
        method_name = method.split('(')[0].strip()
        
        if method_name == "FindAll":
            new_content += f"""
func (r *{struct_name}) FindAll(offset, limit int) ([]models.{model_name}, int64, error) {{
	ctx := context.Background()
	var items []models.{model_name}
	var total int64

	aggQuery := r.db.Collection("{collection_name}").Where("DeletedAt", "==", nil).NewAggregationQuery().WithCount("all")
	res, err := aggQuery.Get(ctx)
	if err == nil {{
		total = res["all"].(*firestore.AggregationResult).Value.(*firestore.CountResult).Value
	}}

	iter := r.db.Collection("{collection_name}").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {{
		doc, err := iter.Next()
		if err == iterator.Done {{
			break
		}}
		if err != nil {{
			return nil, 0, err
		}}
		var item models.{model_name}
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}}

	return items, total, nil
}}
"""
        elif method_name == "FindByID":
            new_content += f"""
func (r *{struct_name}) FindByID(id string) (*models.{model_name}, error) {{
	ctx := context.Background()
	doc, err := r.db.Collection("{collection_name}").Doc(id).Get(ctx)
	if err != nil {{
		return nil, err
	}}
	var item models.{model_name}
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}}
"""
        elif method_name == "Create":
            new_content += f"""
func (r *{struct_name}) Create(item *models.{model_name}) error {{
	ctx := context.Background()
	ref := r.db.Collection("{collection_name}").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}}
"""
        elif method_name == "Update":
            new_content += f"""
func (r *{struct_name}) Update(item *models.{model_name}) error {{
	ctx := context.Background()
	_, err := r.db.Collection("{collection_name}").Doc(item.ID).Set(ctx, item)
	return err
}}
"""
        elif method_name == "Delete":
            new_content += f"""
func (r *{struct_name}) Delete(id string) error {{
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("{collection_name}").Doc(id).Update(ctx, []firestore.Update{{
		{{Path: "DeletedAt", Value: &now}},
	}})
	return err
}}
"""
        elif method_name == "Restore":
            new_content += f"""
func (r *{struct_name}) Restore(id string) error {{
	ctx := context.Background()
	_, err := r.db.Collection("{collection_name}").Doc(id).Update(ctx, []firestore.Update{{
		{{Path: "DeletedAt", Value: nil}},
	}})
	return err
}}
"""
        elif method_name.startswith("FindBy"):
            # Generic finder (e.g. FindByPatientID)
            field = method_name.replace("FindBy", "")
            field_param = field[0].lower() + field[1:]
            
            # Check if it returns array or single
            returns_array = "[]" in method
            
            if returns_array:
                new_content += f"""
func (r *{struct_name}) {method_name}({field_param} string) ([]models.{model_name}, error) {{
	ctx := context.Background()
	var items []models.{model_name}
	iter := r.db.Collection("{collection_name}").Where("{field}", "==", {field_param}).Documents(ctx)
	for {{
		doc, err := iter.Next()
		if err == iterator.Done {{
			break
		}}
		if err != nil {{
			return nil, err
		}}
		var item models.{model_name}
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}}
	return items, nil
}}
"""
            else:
                new_content += f"""
func (r *{struct_name}) {method_name}({field_param} string) (*models.{model_name}, error) {{
	ctx := context.Background()
	iter := r.db.Collection("{collection_name}").Where("{field}", "==", {field_param}).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {{
		return nil, nil
	}}
	if err != nil {{
		return nil, err
	}}
	var item models.{model_name}
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}}
"""

    with open(filepath, 'w') as f:
        f.write(new_content)

print("Repos completely rewritten.")

