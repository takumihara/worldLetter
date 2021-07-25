package main

import (
	"github.com/gorilla/mux"
	"github.com/tacomea/worldLetter/database"
	"github.com/tacomea/worldLetter/repository"
	"github.com/tacomea/worldLetter/usecase"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tpl *template.Template

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tpl = template.Must(template.ParseGlob("templates/*html"))
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

	// handlers
	h := newHandler(uu, su)
	r := mux.NewRouter()

	//private routes
	//r.HandleFunc("/", h.jwtAuth(h.indexHandler)).Methods("GET")
	//r.HandleFunc("/edit", h.jwtAuth(h.editHandler)).Methods("GET")
	//r.HandleFunc("/submit", h.jwtAuth(h.submitHandler)).Methods("POST")

	// public routes
	r.HandleFunc("/", h.indexHandler).Methods("GET")
	r.HandleFunc("/enter", h.enterHandler).Methods("GET")
	r.HandleFunc("/register", h.registerHandler).Methods("POST")
	r.HandleFunc("/login", h.loginHandler).Methods("POST")
	r.HandleFunc("/logout", h.logoutHandler).Methods("POST")

	r.Handle("/favicon.ico", http.NotFoundHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatalln(http.ListenAndServe(":"+port, r))
}
