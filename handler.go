package main

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/tacomea/worldLetter/domain"
	"github.com/tacomea/worldLetter/token"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type handler struct {
	userUseCase    domain.UserUseCase
	sessionUseCase domain.SessionUseCase
}

func newHandler(uu domain.UserUseCase, su domain.SessionUseCase) *handler {
	return &handler{
		userUseCase:    uu,
		sessionUseCase: su,
	}
}

func (h *handler) jwtAuth(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err == nil {

			sessionId, err := token.ParseToken(cookie.Value)
			if err != nil {

				cookie.MaxAge = -1
				http.SetCookie(w, cookie)

				err = h.sessionUseCase.Delete(sessionId)
				if err != nil {
					log.Println("session was not deleted: ", err)
				}

				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			value, err := h.sessionUseCase.Read(cookie.Value)
			if err != nil {

				cookie.MaxAge = -1
				http.SetCookie(w, cookie)

				err = h.sessionUseCase.Delete(sessionId)
				if err != nil {
					log.Println("session was not deleted: ", err)
				}

				http.Redirect(w, r, "/auth", http.StatusSeeOther)
				return
			}

			context.Set(r, "email", value.Email)

			hf.ServeHTTP(w, r)

		} else {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
		}
	}
}

func (h *handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")

	html, err := template.ParseFiles("templates/index.gohtml")

	cookie, err := r.Cookie("session")
	if err == nil {
		sessionId, err := token.ParseToken(cookie.Value)
		if err != nil {
			log.Println("cookie modified")
		} else if session, err := h.sessionUseCase.Read(sessionId); err == nil{
			msg = "Your Email: " + session.Email
		}
	}

	err = html.Execute(w, msg)
	if err != nil {
		log.Println("Error in WriteString: ", err)
	}
}

func (h *handler) enterHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")

	html, err := template.ParseFiles("templates/enter.gohtml")
	if err != nil {
		_, err := w.Write([]byte("500: internal server error"))
		if err != nil {
			log.Println("Error in WriteString: ", err)
		}
		return
	}

	err = html.Execute(w, msg)
	if err != nil {
		log.Println("Error in WriteString: ", err)
	}
}

func (h *handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		query := url.QueryEscape("You cannot when you are not logged in")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}

	sessionId, err := token.ParseToken(cookie.Value)
	if err != nil {
		query := url.QueryEscape("Logout: Cookie Modified")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}

	err = h.sessionUseCase.Delete(sessionId)
	if err != nil {
		log.Println("session was not deleted: ", err)
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	query := url.QueryEscape("successfully logged out")
	http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
}

func (h *handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	encodedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error while hashing", err)
	}

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))
	err = h.userUseCase.Create(domain.User{
		Email:    encodedEmail,
		Password: encodedPassword,
	})

	query := url.QueryEscape("account successfully created")
	http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
}


func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	user, err := h.userUseCase.Read(encodedEmail)
	if err != nil {
		query := url.QueryEscape("username doesn't exist")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}
	hashedPassword := user.Password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		query := url.QueryEscape("login failed")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
	} else {
		sessionId := uuid.NewString()
		err := h.sessionUseCase.Create(domain.Session{
			ID:    sessionId,
			Email: email,
		})
		t, err := token.CreateToken(sessionId)
		if err != nil {
			log.Println("Error in createToken(): ", err)
			query := url.QueryEscape("Server Error, Try Again")
			http.Redirect(w, r, "/?msg="+query, http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "session",
			Value: t,
		}
		http.SetCookie(w, &cookie)
		query := url.QueryEscape("logged in")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
	}
}