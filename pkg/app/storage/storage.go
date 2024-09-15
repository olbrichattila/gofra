package storage

type Storager interface {
	Append(string, string) error
	Put(string, string) error
	Delete(string) error
	Get(string) (string, error)
	HasKey(string) (bool, error)
}
