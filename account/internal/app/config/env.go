package config

import (
	"log"

	"github.com/joho/godotenv"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/utils"
)

func LoadEnv() domain.AppEnv {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	var c domain.AppEnv
	c.MySqlRootAccount = utils.GetEnv("MYSQL_ROOT_USER", "root")
	c.MySqlRootPassword = utils.GetEnv("MYSQL_ROOT_PASSWORD", "root")
	c.MySqlDataBaseName = utils.GetEnv("MYSQL_DATABASE_NAME", "giftino_account")
	c.MySqlDriver = utils.GetEnv("DATABASE_DRIVER", "mysql")
	c.MySqlHost = utils.GetEnv("MYSQL_HOST", "localhost")
	c.MysqlPort = utils.GetEnv("MYSQL_PORT", "3306")
	c.RedisHost = utils.GetEnv("REDIS_HOST", "localhost")
	c.RedisPort = utils.GetEnv("REDIS_PORT", "6379")
	c.HttpAppServerPort = utils.GetEnv("HTTP_APP_SERVER_PORT", "17070")
	c.GrpcAppServerPort = utils.GetEnv("GRPC_APP_SERVER_PORT", "localhost:17071")
	c.HttpInventoryAddress = utils.GetEnv("HTTP_INVENTORY_ADDRESS", "localhost:19090")
	c.GrpcInventoryAddress = utils.GetEnv("GRPC_INVENTORY_ADDRESS", "localhost:19091")
	c.GrpcWalletAddress = utils.GetEnv("GRPC_WALLET_ADDRESS", "localhost:48081")
	c.NatsHost = utils.GetEnv("NATS_HOST", "127.0.0.1")
	c.NatsPort = utils.GetEnv("NATS_PORT", "4222")
	c.MigrationSteps = utils.GetEnv("MIGRATION_STEPS", "7")
	c.MigrationUp = utils.GetEnv("MIGRATION_UP", "true")
	c.MigrationDown = utils.GetEnv("MIGRATION_DOWN", "false")
	return c
}
