package domain

type Server struct {
	Port string
}

type DataBase struct {
	Driver  string
	Name    string
	Host    string
	Port    string
	Account string
	Pass    string
	Migrate int
}

type Redis struct {
	Host string
	Port string
}

type Config struct {
	Server   Server
	DataBase DataBase
	Redis    Redis
}

type ConfigurationService interface {
	GetConfig() (*Config, error)
}
