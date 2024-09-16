package config

type Config struct {
	Host             string
	ConnectionString string
	IsTest           bool
}

func New(
	host string,
	connStr string,
	isTest bool,
) *Config {
	return &Config{
		Host:             host,
		ConnectionString: connStr,
		IsTest:           isTest,
	}
}
