package config

import (
	"fmt"
	"os"
)

func generatePostgresURL() string {
	host := os.Getenv("PG_HOSTNAME")
	port := os.Getenv("PG_PORT")
	dbname := os.Getenv("PG_DB_NAME")
	user := os.Getenv("PG_USERNAME")
	password := os.Getenv("PG_PASSWORD")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
