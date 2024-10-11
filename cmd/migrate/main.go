package main

import (
	"log"
	"os"

	"github.com/Nasa28/CommerceCore/config"
	"github.com/Nasa28/CommerceCore/db"
	mysqlConfig "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

)

func main() {
	db, err := db.NewMysqlDatabase(mysqlConfig.Config{
		User:                 config.Env.DbUser,
		Passwd:               config.Env.DBPassword,
		Addr:                 config.Env.DBAddress,
		DBName:               config.Env.DBName,
		AllowNativePasswords: true,
		Net:                  "tcp",
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations", "mysql", driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
