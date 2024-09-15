package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/sqlc"
)

type Repository struct {
	*sqlc.Queries
	Redis *redis.Client
}

func NewRepository(q *sqlc.Queries, r *redis.Client) *Repository {
	return &Repository{Queries: q, Redis: r}
}

// Upload file to the database
func (r *Repository) UploadFile(ctx context.Context, in *sqlc.UploadFileByEmailParams) error {
	err := r.Queries.UploadFileByEmail(ctx, *in)
	if err != nil {
		return err
	}
	return nil
}

// Get file from the database
func (r *Repository) GetFile(ctx context.Context, in string) (*sqlc.File, error) {
	file, err := r.Queries.GetFileByID(ctx, in)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FilesByUserRequest
func (r *Repository) FilesByUserRequest(ctx context.Context, in string) ([]sqlc.GetFilesByUserRow, error) {
	files, err := r.Queries.GetFilesByUser(ctx, in)
	if err != nil {
		return nil, err
	}
	return files, nil
}
