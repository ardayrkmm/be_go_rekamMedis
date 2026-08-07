import os
import re

directory = r"e:\Bisnis\RekamMedis\backend_go_firebase\internal\repository"

def fix_aggregation():
    for filename in os.listdir(directory):
        if not filename.endswith(".go"): continue
        filepath = os.path.join(directory, filename)
        with open(filepath, 'r') as f:
            content = f.read()
            
        # Fix NewAggregationQuery
        # Replace:
        # aggQuery := r.db.Collection("...").Where("DeletedAt", "==", nil).NewAggregationQuery().WithCount("all")
        # With:
        # q := r.db.Collection("...").Where("DeletedAt", "==", nil)
        # aggQuery := q.NewAggregationQuery().WithCount("all")
        
        pattern = re.compile(r'aggQuery := (r\.db\.Collection\([^)]+\)\.Where\([^)]+\))\.NewAggregationQuery\(\)\.WithCount\("all"\)')
        content = pattern.sub(r'q := \1\n\taggQuery := q.NewAggregationQuery().WithCount("all")', content)
        
        # Some repos don't have Where("DeletedAt") for aggregation maybe?
        pattern2 = re.compile(r'aggQuery := (r\.db\.Collection\([^)]+\))\.NewAggregationQuery\(\)\.WithCount\("all"\)')
        content = pattern2.sub(r'q := \1\n\taggQuery := q.NewAggregationQuery().WithCount("all")', content)

        with open(filepath, 'w') as f:
            f.write(content)

def add_missing_methods():
    # Auth
    auth_file = os.path.join(directory, "auth_repository.go")
    if os.path.exists(auth_file):
        with open(auth_file, 'a') as f:
            f.write("""
func (r *authRepository) CreateBlocklist(blocklist *models.JwtBlocklist) error {
	ctx := context.Background()
	_, err := r.db.Collection("jwt_blocklists").Doc(blocklist.Token).Set(ctx, blocklist)
	return err
}

func (r *authRepository) CreateResetToken(token *models.PasswordResetToken) error {
	ctx := context.Background()
	_, err := r.db.Collection("password_reset_tokens").NewDoc().Set(ctx, token)
	return err
}

func (r *authRepository) GetResetToken(email, token string) (*models.PasswordResetToken, error) {
	ctx := context.Background()
	iter := r.db.Collection("password_reset_tokens").Where("Email", "==", email).Where("Token", "==", token).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var res models.PasswordResetToken
	doc.DataTo(&res)
	return &res, nil
}

func (r *authRepository) DeleteResetToken(email string) error {
	ctx := context.Background()
	iter := r.db.Collection("password_reset_tokens").Where("Email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(ctx)
	}
	return nil
}
""")

    # Master
    master_file = os.path.join(directory, "master_repository.go")
    if os.path.exists(master_file):
        with open(master_file, 'a') as f:
            f.write("""
func (r *masterRepository) GetPatientCategories() ([]models.PatientCategory, error) {
	ctx := context.Background()
	var items []models.PatientCategory
	iter := r.db.Collection("patientcategories").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.PatientCategory
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *masterRepository) GetGenders() ([]models.Gender, error) {
	ctx := context.Background()
	var items []models.Gender
	iter := r.db.Collection("genders").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.Gender
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}
""")

    # Notification
    notif_file = os.path.join(directory, "notification_repository.go")
    if os.path.exists(notif_file):
        with open(notif_file, 'a') as f:
            f.write("""
func (r *notificationRepository) FindUnread() ([]models.Notification, error) {
	ctx := context.Background()
	var items []models.Notification
	iter := r.db.Collection("notifications").Where("IsRead", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.Notification
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *notificationRepository) MarkAllAsRead() error {
	ctx := context.Background()
	iter := r.db.Collection("notifications").Where("IsRead", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return err }
		doc.Ref.Update(ctx, []firestore.Update{{Path: "IsRead", Value: true}})
	}
	return nil
}
""")

fix_aggregation()
add_missing_methods()
print("Fixed.")
