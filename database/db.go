package database

import (
	"fmt"
	"log"
	"time"

	"readcommend/configuration"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB(config configuration.Configuration) *sqlx.DB {
	var err error

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DBUsername,
		config.DBPassword,
		config.DBHostname,
		config.DBPort,
		config.DBDatabaseName,
		config.DBSSLMode,
	)

	for i := 0; i < 30; i++ {
		db, err = sqlx.Open("postgres", connStr)
		if err != nil {
			log.Printf("Error connecting to the database (attempt %d): %v\n", i+1, err)
			time.Sleep(time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Error pinging database (attempt %d): %v\n", i+1, err)
			db.Close()
			time.Sleep(time.Second)
			continue
		}

		log.Printf("The database is connected")
		return db
	}

	log.Fatalf("Failed to connect to the database after %d attempts", 30)
	return nil
}
