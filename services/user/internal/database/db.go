package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	zipkinsql "github.com/openzipkin-contrib/zipkin-go-sql"
	"github.com/openzipkin/zipkin-go"
	"github.com/redis/go-redis/v9"
	"github.com/samarthasthan/21BRS1248_Backend/services/user/internal/database/sqlc"
)

type Database interface {
	Connect(string, string) error
	Close() error
	RegisterZipkin(*zipkin.Tracer) string
}

type Postgres struct {
	Queries *sqlc.Queries
	DB      *sql.DB
}

type Redis struct {
	*redis.Client
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func NewRedis() *Redis {
	return &Redis{}
}

func (s *Postgres) Connect(driverName string, addr string) error {
	db, err := sql.Open(driverName, addr)
	if err != nil {
		return err
	}
	s.DB = db
	s.Queries = sqlc.New(db)
	return nil
}

func (s *Postgres) RegisterZipkin(tracer *zipkin.Tracer) string {
	// Register our zipkinsql wrapper for the provided MySQL driver.
	driverName, err := zipkinsql.Register("postgres", tracer, zipkinsql.WithAllTraceOptions())
	if err != nil {
		log.Fatalf("unable to register zipkin driver: %v\n", err)
	}
	return driverName
}

func (s *Postgres) Close() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Connect(addr string) error {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	r.Client = rdb
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Close() error {
	err := r.Close()
	if err != nil {
		return err
	}
	return nil
}
