package db

import (
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/lib/pq"
)

type PgSQLConfig struct {
}

func newPgsqlConfig() DBConfiger {
	return &PgSQLConfig{}
}

func (c *PgSQLConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv(envdbUserName),
		os.Getenv(envdbPassword),
		os.Getenv(envdbHost),
		os.Getenv(envdbPort),
		os.Getenv(envdbDatabase),
		os.Getenv(envdbSSLMode),
	)
}

func (c *PgSQLConfig) GetConnectionName() string {
	return DriverNamePostgres
}
