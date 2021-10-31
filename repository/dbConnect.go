package repository

import (
	"github.com/RobotApocalypse/configuration"
	"github.com/RobotApocalypse/constants"

	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/sirupsen/logrus"
)

type DBConnect interface {
	ConnectDB() *sql.DB
	Migrate()
}

type dbConnect struct {
	cfg configuration.Config
	log *logrus.Entry
}

func (d dbConnect) ConnectDB() *sql.DB {
	var err error
	crsqlQuery := fmt.Sprintf("%s:%s@(%s:%s)/", d.cfg.GetString(constants.DBUserName), d.cfg.GetOsEnvString(d.cfg.GetString(constants.DBPassword)),
		d.cfg.GetString(constants.DBHost), d.cfg.GetString(constants.DBPort))
	crdb, err := sql.Open(d.cfg.GetString(constants.Database),crsqlQuery)
	if err != nil {
		d.log.Fatalf("error: connecting db - %s", err.Error())
	}

	defer crdb.Close()

	_,err = crdb.Exec("CREATE DATABASE IF NOT EXISTS "+d.cfg.GetString(constants.DBName)+" DEFAULT COLLATE utf8_general_ci")
	if err != nil {
		d.log.Fatalf("error: executing db init query - %s", err.Error())
	}

	_,err = crdb.Exec("USE "+ d.cfg.GetString(constants.DBName))
	if err != nil {
		d.log.Fatalf("error: executing use db query - %s", err.Error())
	}

	sqlQuery := fmt.Sprintf("%s:%s@(%s:%s)/%s",d.cfg.GetString(constants.DBUserName), d.cfg.GetOsEnvString(d.cfg.GetString(constants.DBPassword)),
		d.cfg.GetString(constants.DBHost) ,d.cfg.GetString(constants.DBPort), d.cfg.GetString(constants.DBName))

	db, err := sql.Open(d.cfg.GetString(constants.Database),sqlQuery)
	if err != nil {
		d.log.Fatalf("error: exec schema migration table - %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		d.log.Fatalf("error: pinging database - %s", err.Error())
	}

	db.SetMaxOpenConns(d.cfg.GetInt(constants.DBMaxOpenConnection))
	db.SetMaxIdleConns(d.cfg.GetInt(constants.DBMaxIdleConnection))
	db.SetConnMaxLifetime(time.Duration(d.cfg.GetInt(constants.DBMaxLifeTime)) * time.Minute)
	return db
}

/**
 * @Description: Migrate updates the schema and table init in database listed in the file in directory  /migrations/mysql
 */
func (d dbConnect) Migrate() {
	sqlQuery := fmt.Sprintf("%s:%s@(%s:%s)/%s",d.cfg.GetString(constants.DBUserName), d.cfg.GetOsEnvString(d.cfg.GetString(constants.DBPassword)),
		d.cfg.GetString(constants.DBHost) ,d.cfg.GetString(constants.DBPort), d.cfg.GetString(constants.DBName))

	db, err := sql.Open("mysql", sqlQuery)
	if err != nil {
		d.log.Fatalf("error: migrate conn open - %s", err.Error())
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		d.log.Fatalf("error: migrate (WithInstance) - %s", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/mysql",
		"mysql",
		driver)
	if err != nil {
		d.log.Fatalf("error: migrate NewWithDatabaseInstance  - %s", err.Error())
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange{
		d.log.Fatalf("error: migration up - %s", err.Error())
	}
}


func NewDBConnect(cfg configuration.Config, log *logrus.Entry) DBConnect {
	return &dbConnect{
		cfg: cfg,
		log: log,
	}
}

