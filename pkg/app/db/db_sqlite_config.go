package db

import (
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/mattn/go-sqlite3"
)

func newSqliteConfig() DBConfiger {
	return &SqLiteConfig{}
}

type SqLiteConfig struct {
}

func (c *SqLiteConfig) GetConnectionString() string {
	return os.Getenv(envdbDatabase)
}

func (c *SqLiteConfig) GetConnectionName() string {
	return DriverNameSqLite
}
