package storage

import (
	"fmt"

	"github.com/olbrichattila/gofra/pkg/app/db"

	builder "github.com/olbrichattila/gosqlbuilder/pkg"
)

func NewDatabaseStorage(s string, db db.DBer, sqlBuilder builder.Builder) Storager {
	return &DbStore{
		tablename:  s,
		db:         db,
		sqlBuilder: sqlBuilder,
	}
}

type DbStore struct {
	tablename  string
	db         db.DBer
	sqlBuilder builder.Builder
}

func (s *DbStore) Append(key string, value string) error {
	sql, err := s.sqlBuilder.Insert(s.tablename).Fields("name", "message").Values(key, value).AsSQL()
	if err != nil {
		return err
	}

	_, err = s.db.Execute(sql, s.sqlBuilder.GetParams()...)
	return err

}

func (s *DbStore) Put(key string, value string) error {
	err := s.Delete(key)
	if err != nil {
		return err
	}
	return s.Append(key, value)
}

func (s *DbStore) Delete(key string) error {
	sql, err := s.sqlBuilder.Delete(s.tablename).Where("name", "=", key).AsSQL()
	if err != nil {
		return err
	}

	_, err = s.db.Execute(sql, s.sqlBuilder.GetParams()...)
	return err
}

func (s *DbStore) HasKey(key string) (bool, error) {

	sql, err := s.sqlBuilder.Select(s.tablename).RawFields("count(*) as db").Where("name", "=", key).AsSQL()
	if err != nil {
		return false, err
	}

	row, err := s.db.QueryOne(sql, s.sqlBuilder.GetParams()...)
	if err != nil {
		return false, err
	}

	if message, ok := row["db"]; ok {
		if message.(int64) > 0 {
			return true, nil
		}

		return false, nil
	}
	return false, fmt.Errorf("unknown session fetch error, field DB missing from response")
}

func (s *DbStore) Get(key string) (string, error) {
	sql, err := s.sqlBuilder.Select(s.tablename).Fields("message").Where("name", "=", key).AsSQL()
	if err != nil {
		return "", err
	}

	row, err := s.db.QueryOne(sql, s.sqlBuilder.GetParams()...)
	if err != nil {
		return "", err
	}
	if message, ok := row["message"]; ok {
		return message.(string), nil
	}
	return "", fmt.Errorf("unknown session fetch error, field message missing from response")
}
