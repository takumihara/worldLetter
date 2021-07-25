package main

import (
	"github.com/tacomea/worldLetter/database"
	"github.com/tacomea/worldLetter/repository"
	"github.com/tacomea/worldLetter/usecase"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// Postgres
	db := database.NewPostgresDB()
	ur := repository.NewUserRepositoryMySQL(db)
	sr := repository.NewSessionRepositoryMySQL(db)

	// sync.Map
	//ur := repository.NewSyncMapUserRepository()
	//sr := repository.NewSyncMapSessionRepository()

	uu := usecase.NewUserUsecase(ur)
	su := usecase.NewSessionUsecase(sr)

	h := newHandler(uu, su)
	http.HandleFunc("/", h.indexHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
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
