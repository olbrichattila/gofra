package db

type DBConfiger interface {
	GetConnectionName() string
	GetConnectionString() string
}
