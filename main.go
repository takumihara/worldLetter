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
	ur := repository.NewUserRepositoryPG(db)
	sr := repository.NewSessionRepositoryPG(db)
	lr := repository.NewLetterRepositoryPG(db)

	// sync.Map
	//ur := repository.NewSyncMapUserRepository()
	//sr := repository.NewSyncMapSessionRepository()
	//lr := repository.NewSyncMapLetterRepository()

	uu := usecase.NewUserUsecase(ur)
	su := usecase.NewSessionUsecase(sr)
	lu := usecase.NewLetterUsecase(lr)

	// handlers
	h := newHandler(uu, su, lu)
	r := mux.NewRouter()

	//private routes
	r.HandleFunc("/create", h.jwtAuth(h.createHandler)).Methods("GET")
	r.HandleFunc("/send", h.jwtAuth(h.sendHandler)).Methods("POST")
	r.HandleFunc("/show", h.jwtAuth(h.showHandler)).Methods("GET")
	//r.HandleFunc("/submit", h.jwtAuth(h.submitHandler)).Methods("POST")

	// public routes
	r.HandleFunc("/signin", h.signinHandler).Methods("GET")
	r.HandleFunc("/signup", h.signupHandler).Methods("GET")
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
