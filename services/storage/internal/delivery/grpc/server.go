package grpc

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/kafka"
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"github.com/samarthasthan/21BRS1248_Backend/services/storage/internal/database/repository"
)

type StorageService struct {
	proto_go.UnimplementedFileServiceServer
	repo *repository.Repository
	log  *logger.Logger
	p    *kafka.Producer
	c    *kafka.Consumer
}

func NewStorageService(log *logger.Logger, repo *repository.Repository, p *kafka.Producer, c *kafka.Consumer) *StorageService {
	return &StorageService{log: log, repo: repo, p: p, c: c}
}
