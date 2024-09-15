package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Enver interface {
	Construct()
	Get(string) string
}

type Env struct {
	isLoaded bool
}

func New() Enver {
	return &Env{
		isLoaded: false,
	}
}

var envFileNames = []string{".env.migrator", ".env"}

func (e *Env) Construct() {
	if e.isLoaded {
		return
	}
	for _, envFileName := range envFileNames {
		e.loadEnvIfExits(envFileName)
	}

	e.isLoaded = true
}

func (e *Env) loadEnvIfExits(envFileName string) {
	_, err := os.Stat(envFileName)
	if err == nil {
		err := godotenv.Load(envFileName)
		if err != nil {
		}
		return
	}
}

func (*Env) Get(name string) string {
	return os.Getenv(name)
}
