package db

import (
	"fmt"
	"os"
)

func NewDBFactory() DBFactoryer {
	return &DBconf{}
}

type DBFactoryer interface {
	GetConnectionConfig() (DBConfiger, error)
}

const (
	DriverNameFirebird = "firebirdsql"
	DriverNameSqLite   = "sqlite3"
	DriverNameMySQL    = "mysql"
	DriverNamePostgres = "postgres"

	DbConnectionTypeSqLite   = "sqlite"
	DbConnectionTypeMySQL    = "mysql"
	DbConnectionTypePgSQL    = "pgsql"
	DbConnectionTypeFirebird = "firebird"
	DbConnectionTypeMemory   = "memory"

	EnvdbConnection = "DB_CONNECTION"
	envdbUserName   = "DB_USERNAME"
	envdbPassword   = "DB_PASSWORD"
	envdbHost       = "DB_HOST"
	envdbPort       = "DB_PORT"
	envdbDatabase   = "DB_DATABASE"
	envdbSSLMode    = "DB_SSLMODE"
)

type DBconf struct {
	config DBConfiger
}

func (c *DBconf) GetConnectionConfig() (DBConfiger, error) {
	dbConnection := os.Getenv(EnvdbConnection)

	switch dbConnection {
	case DbConnectionTypeSqLite:
		return newSqliteConfig(), nil
	case DbConnectionTypeMySQL:
		return newMySQLConfig(), nil
	case DbConnectionTypePgSQL:
		return newPgsqlConfig(), nil
	case DbConnectionTypeFirebird:
		return newFirebirdConfig(), nil
	case DbConnectionTypeMemory:
		return newMemoryDBConfig(), nil
	default:
		return nil, fmt.Errorf("invalid DB_CONNECTION %s", dbConnection)
	}
}
