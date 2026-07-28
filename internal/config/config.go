package config

type Config struct {
	Addr string
	Env  string

	Database DatabaseConfig
}

type DatabaseConfig struct {
	Addr               string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        string
}
