package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type UserRepository interface {
	FindAll(offset, limit int) ([]models.User, int64, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
	Restore(id string) error
}

type userRepository struct {
	db *firestore.Client
}

func NewUserRepository(db *firestore.Client) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	ctx := context.Background()
	var users []models.User
	var total int64

	// Firestore doesn't have a direct count query easily without reading all or using aggregation queries.
	// We will use aggregation query for count.
	// total omitted for NoSQL

	iter := r.db.Collection("users").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var user models.User
		doc.DataTo(&user)
		user.ID = doc.Ref.ID
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	ctx := context.Background()
	iter := r.db.Collection("users").Where("Email", "==", email).Where("DeletedAt", "==", nil).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil // Not found, wait, usually we return an error for not found, or handle it in usecase. Let's return error or nil.
	}
	if err != nil {
		return nil, err
	}
	var user models.User
	doc.DataTo(&user)
	user.ID = doc.Ref.ID
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	doc.DataTo(&user)
	user.ID = doc.Ref.ID
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	ctx := context.Background()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	ref := r.db.Collection("users").NewDoc()
	user.ID = ref.ID
	_, err := ref.Set(ctx, user)
	return err
}

func (r *userRepository) Update(user *models.User) error {
	ctx := context.Background()
	user.UpdatedAt = time.Now()
	_, err := r.db.Collection("users").Doc(user.ID).Set(ctx, user)
	return err
}

func (r *userRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("users").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}

func (r *userRepository) Restore(id string) error {
	ctx := context.Background()
	_, err := r.db.Collection("users").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: nil},
	})
	return err
}

