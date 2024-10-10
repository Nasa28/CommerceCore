package main

import (
	"database/sql"
	"log"

	"github.com/Nasa28/CommerceCore/cmd/api"
	"github.com/Nasa28/CommerceCore/config"
	"github.com/Nasa28/CommerceCore/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
log.Println(config.Env.DBAddress)
	db, err := db.NewMysqlDatabase(mysql.Config{
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

	initDB(db)
	server := api.NewAPIServer(config.Env.Port, db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initDB(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")
}
