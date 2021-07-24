package main

import (
	"github.com/tacomea/worldLetter/repository"
	"github.com/tacomea/worldLetter/usecase"
	"log"
	"net/http"
	"os"
)

//func connect() *gorm.DB {
//	connection, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
//	if err != nil {
//		log.Panic("Could not connect to the database")
//	}
//
//	connection.AutoMigrate(&domain.User{})
//	connection.AutoMigrate(&domain.Session{})
//
//	return connection
//}

func main() {
	// Postgres
	//db := connect()
	//ur := repository.NewUserRepositoryMySQL(db)
	//sr := repository.NewSessionRepositoryMySQL(db)

	// sync.Map
	ur := repository.NewSyncMapUserRepository()
	sr := repository.NewSyncMapSessionRepository()

	uu := usecase.NewUserUsecase(ur)
	su := usecase.NewSessionUsecase(sr)

	h := newHandler(uu, su)
	http.HandleFunc("/", h.indexHandler)
	http.HandleFunc("/enter", h.enterHandler)
	http.HandleFunc("/logout", h.logoutHandler)
	http.HandleFunc("/register", h.registerHandler)
	http.HandleFunc("/login", h.loginHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
