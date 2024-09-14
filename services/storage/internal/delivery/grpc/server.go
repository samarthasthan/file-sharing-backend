package grpc

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/repository"
)

type StorageService struct {
	proto_go.UnimplementedFileServiceServer
	repo *repository.Repository
	log  *logger.Logger
}

func NewStorageService(log *logger.Logger, repo *repository.Repository) *StorageService {
	return &StorageService{log: log, repo: repo}
}
