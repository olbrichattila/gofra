package db

import (
	"fmt"
	"os"

	// This blank import is necessary to have the driver
	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
}

func newMySQLConfig() DBConfiger {
	return &MySQLConfig{}
}

func (c *MySQLConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv(envdbUserName),
		os.Getenv(envdbPassword),
		os.Getenv(envdbHost),
		os.Getenv(envdbPort),
		os.Getenv(envdbDatabase),
	)
}

func (c *MySQLConfig) GetConnectionName() string {
	return DriverNameMySQL
}
