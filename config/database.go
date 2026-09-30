package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase() (*pgxpool.Pool, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)

	pool, err := pgxpool.New(
		context.Background(),
		dsn,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuat koneksi database: %w",
			err,
		)
	}

	// Test koneksi ke database
	err = pool.Ping(context.Background())

	if err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"gagal terhubung ke PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
