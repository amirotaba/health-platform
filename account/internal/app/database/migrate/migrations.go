package migrate

import (
	"embed"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
)

var (
	//go:embed migrations/*
	migrations embed.FS
)

const migrationDirectory = "migrations"

func Up(conf domain.AppEnv) error {
	// todo use switch case for when we need handle multiple sql driver
	//conf := config.LoadEnv()
	//log.Info("storage: applying MysqlQL data migrations")
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		conf.MySqlRootAccount, conf.MySqlRootPassword, conf.MySqlHost,
		conf.MysqlPort, conf.MySqlDataBaseName)
	db, err := sqlx.Open(conf.MySqlDriver, dns)
	//db, err := sqlx.Open("mysql", "root:root@(172.19.0.2:3306)/test?multiStatements=true&parseTime=true")
	if err != nil {
		log.Println(fmt.Errorf("storage: connect mysql driver error: %w", err))
		return err
	}

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Println(fmt.Errorf("storage: migrate mysql driver error: %w", err))
		return err
	}

	src, err := httpfs.New(http.FS(migrations), migrationDirectory)
	if err != nil {
		log.Println(fmt.Errorf("new httpfs error: %w", err))
		return err
	}

	m, err := migrate.NewWithInstance("httpfs", src, conf.MySqlDriver, driver)
	if err != nil {
		log.Println(fmt.Errorf("storage: new migrate instance error: %w", err))
		return err
	}

	oldVersion, _, _ := m.Version()

	steps, err := strconv.Atoi(conf.MigrationSteps)
	log.Println(steps)
	if err = m.Steps(steps); err != nil && err != migrate.ErrNoChange {
		log.Println(fmt.Errorf("storage: migrate up error: %w", err))
		return err
	}

	//if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//	log.Println(fmt.Errorf("storage: migrate up error: %w", err))
	//	return err
	//}

	newVersion, _, _ := m.Version()
	if oldVersion != newVersion {
		log.WithFields(log.Fields{
			"from_version": oldVersion,
			"to_version":   newVersion,
		}).Info("storage: MysqlQL data migrations applied")

	}

	return nil
}

func Down(conf domain.AppEnv) error {
	// todo use switch case for when we need handle multiple sql driver
	log.Info("storage: applying MysqlQL data migrations")
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true&character-set-server=utf8mb4",
		conf.MySqlRootAccount, conf.MySqlRootPassword, conf.MySqlHost,
		conf.MysqlPort, conf.MySqlDataBaseName)
	db, err := sqlx.Open(conf.MySqlDriver, dns)
	//db, err := sqlx.Open("mysql", "root:root@(172.19.0.2:3306)/test?multiStatements=true&parseTime=true")
	if err != nil {
		log.Println(fmt.Errorf("storage: connect mysql driver error: %w", err))
		return err
	}

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Println(fmt.Errorf("storage: migrate mysql driver error: %w", err))
		return err
	}

	src, err := httpfs.New(http.FS(migrations), migrationDirectory)
	if err != nil {
		log.Println(fmt.Errorf("new httpfs error: %w", err))
		return err
	}

	m, err := migrate.NewWithInstance("httpfs", src, conf.MySqlDriver, driver)
	if err != nil {
		log.Println(fmt.Errorf("storage: new migrate instance error: %w", err))
		return err
	}

	oldVersion, _, _ := m.Version()
	//if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	//	log.Println(fmt.Errorf("storage: migrate up error: %w", err))
	//	return err
	//}

	//steps, err := strconv.Atoi(config.MigrationStep)
	if err = m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Println(fmt.Errorf("storage: migrate up error: %w", err))
		return err
	}

	newVersion, _, _ := m.Version()
	if oldVersion != newVersion {
		log.WithFields(log.Fields{
			"from_version": oldVersion,
			"to_version":   newVersion,
		}).Info("storage: MysqlQL data migrations applied")

	}

	return nil
}
