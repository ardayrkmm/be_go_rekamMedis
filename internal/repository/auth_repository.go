package repository

import (
	"context"
	"backend_go/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type AuthRepository interface {
	CreateResetToken(token *models.PasswordResetToken) error
	GetResetToken(email, token string) (*models.PasswordResetToken, error)
	DeleteResetToken(email string) error
	CreateBlocklist(blocklist *models.JwtBlocklist) error
}

type authRepository struct {
	db *firestore.Client
}

func NewAuthRepository(db *firestore.Client) AuthRepository {
	return &authRepository{db}
}

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
