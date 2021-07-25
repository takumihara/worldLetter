package database

import (
	"database/sql"
	"github.com/tacomea/worldLetter/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func NewPostgresDB() *gorm.DB {
	url := getDatabaseUrl()
	//max := getMaxConnection()

	sqlDB, err := sql.Open("postgres", url)
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

	return connection
}

func getDatabaseUrl() string {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		panic("Environment variable 'DATABASE_URL' not defined")
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
