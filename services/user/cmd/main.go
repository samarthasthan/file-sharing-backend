package main

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/pkg/env"
)

var (
	USER_GRPC_PORT         string
	USER_DB_PORT           string
	USER_POSTGRES_USER     string
	USER_POSTGRES_PASSWORD string
	USER_POSTGRES_DB       string
	USER_POSTGRES_HOST     string
)

func init() {
	USER_GRPC_PORT = env.GetEnv("USER_GRPC_PORT", "8000")
	USER_DB_PORT = env.GetEnv("USER_DB_PORT", "5432")
	USER_POSTGRES_USER = env.GetEnv("USER_POSTGRES_USER", "root")
	USER_POSTGRES_PASSWORD = env.GetEnv("USER_POSTGRES_PASSWORD", "password")
	USER_POSTGRES_DB = env.GetEnv("USER_POSTGRES_DB", "user-db")
	USER_POSTGRES_HOST = env.GetEnv("USER_POSTGRES_HOST", "localhost")
}

func main() {
}
