package domain

type AppEnv struct {
	NatsHost             string
	NatsPort             string
	MySqlHost            string
	MysqlPort            string
	RedisHost            string
	RedisPort            string
	MySqlDriver          string
	MySqlDataBaseName    string
	MySqlRootAccount     string
	MySqlRootPassword    string
	MigrationSteps       string
	MigrationDown        string
	MigrationUp          string
	HttpAppServerPort    string
	GrpcAppServerPort    string
	HttpInventoryAddress string
	GrpcInventoryAddress string
	GrpcWalletAddress    string
}
