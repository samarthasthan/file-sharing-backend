package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	zipkinsql "github.com/openzipkin-contrib/zipkin-go-sql"
	"github.com/openzipkin/zipkin-go"
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

func NewPostgres() *Postgres {
	return &Postgres{}
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
	driverName, err := zipkinsql.Register("mysql", tracer, zipkinsql.WithAllTraceOptions())
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
