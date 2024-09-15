package db

type MemoryConfig struct {
}

func newMemoryDBConfig() DBConfiger {
	return &MemoryConfig{}
}

func (c *MemoryConfig) GetConnectionString() string {
	return ":memory:"
}

func (c *MemoryConfig) GetConnectionName() string {
	return DriverNameSqLite
}
