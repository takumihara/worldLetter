package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewPostgresDB() *gorm.DB {
	url := getDatabaseUrl()
	//max := getMaxConnection()

	sqlDB, err :=sql.Open("postgres", url)
	if err != nil {
		log.Fatalln(err)
	}
	connection, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	//conn.DB().SetMaxIdleConns(max / 5)
	//conn.DB().SetMaxOpenConns(max)

	connection.AutoMigrate(&domain.User{})
	connection.AutoMigrate(&domain.Session{})
	connection.AutoMigrate(&domain.Letter{})

	return connection
}

func getDatabaseUrl() string {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		pgUser := os.Getenv("POSTGRES_USER")
		pgPass := os.Getenv("POSTGRES_PASSWORD")
		pgHost := "postgresql"
		pgPort := "5432"
		pgDB := os.Getenv("POSTGRES_DB")
		//panic("Environment variable 'DATABASE_URL' not defined")
		url = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", pgUser, pgPass, pgHost, pgPort, pgDB)
		//log.Println("environment variable 'DATABASE_URL' not defined: ", url)
		//url = "user=pg-user dbname=world-letter password=password sslmode=disable port=5433 host=localhost"
	}

	return url
}

//func getMaxConnection() int {
//	env := os.Getenv("DATABASE_MAX_CONNECTIONS")
//	if env == "" {
//		return defaultMaxConnections
//	}
//
//	max, _ := strconv.Atoi(env)
//	return max
//}
