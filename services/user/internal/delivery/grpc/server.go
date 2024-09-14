package grpc

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/logger"
	"github.com/samarthasthan/21BRS1248_Backend/common/proto_go"
	"github.com/samarthasthan/21BRS1248_Backend/services/user/internal/database/repository"
)

type UserService struct {
	proto_go.UnimplementedUserServiceServer
	repo *repository.Repository
	log  *logger.Logger
}

func NewUserService(repo *repository.Repository, log *logger.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}
