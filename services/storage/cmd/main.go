package main

import (
	"github.com/samarthasthan/21BRS1248_Backend/common/pkg/env"
)

var (
	STORAGE_GRPC_PORT         string
	STORAGE_DB_PORT           string
	STORAGE_POSTGRES_STORAGE     string
	STORAGE_POSTGRES_PASSWORD string
	STORAGE_POSTGRES_DB       string
	STORAGE_POSTGRES_HOST     string
)

func init() {
	STORAGE_GRPC_PORT = env.GetEnv("STORAGE_GRPC_PORT", "8000")
	STORAGE_DB_PORT = env.GetEnv("STORAGE_DB_PORT", "5432")
	STORAGE_POSTGRES_STORAGE = env.GetEnv("STORAGE_POSTGRES_STORAGE", "root")
	STORAGE_POSTGRES_PASSWORD = env.GetEnv("STORAGE_POSTGRES_PASSWORD", "password")
	STORAGE_POSTGRES_DB = env.GetEnv("STORAGE_POSTGRES_DB", "STORAGE-db")
	STORAGE_POSTGRES_HOST = env.GetEnv("STORAGE_POSTGRES_HOST", "localhost")
}

func main() {
}
