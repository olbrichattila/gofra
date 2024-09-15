package storage

import (
	"fmt"
	"io"
	"os"
)

func NewFileStorage() Storager {
	return &Fstorage{}
}

type Fstorage struct {
}

func (s *Fstorage) Append(key string, value string) error {
	return s.save(os.O_APPEND|os.O_CREATE|os.O_WRONLY, key, value)
}

func (s *Fstorage) Put(key string, value string) error {
	return s.save(os.O_CREATE|os.O_WRONLY|os.O_TRUNC, key, value)
}

func (s *Fstorage) Delete(key string) error {
	if _, err := os.Stat(key); err == nil {
		err = os.Remove(key)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("file does not exists")
	} else {
		return err
	}

	return nil
}

func (s *Fstorage) HasKey(key string) (bool, error) {
	if _, err := os.Stat(key); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func (s *Fstorage) Get(key string) (string, error) {
	file, err := os.Open(key)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *Fstorage) save(flag int, key string, value string) error {
	file, err := os.OpenFile(key, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(value)
	if err != nil {
		return err
	}

	return nil
}
