package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"github.com/samarthasthan/21BRS1248_Backend/services/user/internal/database/sqlc"
)

type Repository struct {
	*sqlc.Queries
	Redis *redis.Client
}

func NewRepository(q *sqlc.Queries, reids *redis.Client) *Repository {
	return &Repository{Queries: q, Redis: reids}
}

// RegisterUser registers a new user
func (r *Repository) RegisterUser(ctx context.Context, in *proto_go.RegisterRequest, uuid string) error {
	err := r.Queries.RegisterUser(ctx, sqlc.RegisterUserParams{
		Userid:    uuid,
		Firstname: in.FirstName,
		Lastname:  in.LastName,
		Email:     in.Email,
		Password:  in.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

// Login logs in a user
func (r *Repository) GetPasswordByEmail(ctx context.Context, in *proto_go.LoginRequest) (string, error) {
	email, err := r.Queries.GetPasswordByEmail(ctx, in.Email)
	if err != nil {
		return "", err
	}
	return email, nil
}
