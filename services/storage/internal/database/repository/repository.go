package repository

import (
	"context"

	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/sqlc"
)

type Repository struct {
	*sqlc.Queries
}

func NewRepository(q *sqlc.Queries) *Repository {
	return &Repository{Queries: q}
}

// Upload file to the database
func (r *Repository) UploadFile(ctx context.Context, in *sqlc.UploadFileByEmailParams) error {
	err := r.Queries.UploadFileByEmail(ctx, *in)
	if err != nil {
		return err
	}
	return nil
}
