package main

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/tacomea/worldLetter/domain"
	"github.com/tacomea/worldLetter/token"
	"github.com/tacomea/worldLetter/usecase"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type handler struct {
	userUseCase    usecase.UserUseCase
	sessionUseCase usecase.SessionUseCase
	letterUseCase  usecase.LetterUseCase
}

func newHandler(uu usecase.UserUseCase, su usecase.SessionUseCase, lu usecase.LetterUseCase) *handler {
	return &handler{
		userUseCase:    uu,
		sessionUseCase: su,
		letterUseCase:  lu,
	}
}

func (h *handler) jwtAuth(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err == nil {

			sessionId, err := token.ParseToken(cookie.Value)
			if err != nil {
				log.Println(err)
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)

				err = h.sessionUseCase.Delete(sessionId)
				if err != nil {
					log.Println("session was not deleted: ", err)
				}

				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			value, err := h.sessionUseCase.Read(sessionId)
			if err != nil {
				log.Println(err)
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)

				err = h.sessionUseCase.Delete(sessionId)
				if err != nil {
					log.Println("session was not deleted: ", err)
				}

				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			context.Set(r, "email", value.Email)

			hf.ServeHTTP(w, r)

		} else {
			log.Println(err)
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		}
	}
}

func (h *handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	err := tpl.ExecuteTemplate(w, "index.html", msg)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) signinHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")

	cookie, err := r.Cookie("session")
	if err == nil {
		sessionId, err := token.ParseToken(cookie.Value)
		if err != nil {
			log.Println("cookie modified")
		} else if session, err := h.sessionUseCase.Read(sessionId); err == nil {
			msg = "Your Email: " + session.Email
		}
	}

	err = tpl.ExecuteTemplate(w, "signin.html", msg)
	if err != nil {
		log.Println("Error in WriteString: ", err)
	}
}

func (h *handler) signupHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")

	err := tpl.ExecuteTemplate(w, "signup.html", msg)
	if err != nil {
		log.Println("Error in WriteString: ", err)
	}
}

func (h *handler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		query := url.QueryEscape("You cannot when you are not logged in")
		http.Redirect(w, r, "/signin?msg="+query, http.StatusSeeOther)
		return
	}

	sessionId, err := token.ParseToken(cookie.Value)
	if err != nil {
		query := url.QueryEscape("Logout: Cookie Modified")
		http.Redirect(w, r, "/signin?msg="+query, http.StatusSeeOther)
		return
	}

	err = h.sessionUseCase.Delete(sessionId)
	if err != nil {
		log.Println("session was not deleted: ", err)
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	query := url.QueryEscape("successfully logged out")
	http.Redirect(w, r, "/signin?msg="+query, http.StatusSeeOther)
}

func (h *handler) registerHandler(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/signin?msg="+query, http.StatusSeeOther)
}

func (h *handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	user, err := h.userUseCase.Read(encodedEmail)
	if err != nil {
		query := url.QueryEscape("username doesn't exist")
		http.Redirect(w, r, "/signin?msg="+query, http.StatusSeeOther)
		return
	}
	hashedPassword := user.Password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		query := url.QueryEscape("login failed")
		http.Redirect(w, r, "/signin?msg="+query, http.StatusSeeOther)
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
			http.Redirect(w, r, "/signin?msg="+query, http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:  "session",
			Value: t,
			Path:  "/",
		}
		http.SetCookie(w, &cookie)
		query := url.QueryEscape("logged in")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
	}
}

func (h *handler) createHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	err := tpl.ExecuteTemplate(w, "create.html", msg)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) sendHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("letter")
	// there is actually no need to do this somehow
	if content == "" {
		msg := url.QueryEscape("you have to write something")
		http.Redirect(w, r, "/create?msg="+msg, http.StatusSeeOther)
		return
	}

	id := uuid.NewString()
	email := context.Get(r, "email").(string)
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	newLetter := domain.Letter{
		ID:       id,
		AuthorID: encodedEmail,
		Content:  content,
		IsSent:   false,
	}

	err := h.letterUseCase.Create(newLetter)
	if err != nil {
		log.Println(err)
		query := url.QueryEscape("sorry, internal server error")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}

	// check opt out
	optout := r.FormValue("optout")
	fmt.Println(optout)
	if optout == "on" {
		msg := url.QueryEscape("Thank you for sending a letter")
		http.Redirect(w, r, "/?msg"+msg, http.StatusSeeOther)
		return
	}

	retrievedLetter, err := h.letterUseCase.GetFirstUnsendLetter(encodedEmail)
	if err != nil {
		log.Println(err)
		msg := "sorry, internal server error"
		http.Redirect(w, r, "/?msg="+msg, http.StatusSeeOther)
	} else if retrievedLetter.Content == "" {
		err = tpl.ExecuteTemplate(w, "send.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}
	retrievedLetter.IsSent = true
	retrievedLetter.ReceiverID = encodedEmail
	err = h.letterUseCase.Update(retrievedLetter)
	if err != nil {
		log.Println(err)
	}

	err = tpl.ExecuteTemplate(w, "send.html", retrievedLetter.Content)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) letterSentHandler(w http.ResponseWriter, r *http.Request) {
	email := context.Get(r, "email").(string)
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	var contents []string

	letters, err := h.letterUseCase.GetLettersByAuthorID(encodedEmail)
	if err != nil {
		log.Println(err)
	} else if letters != "" {
		contents = strings.Split(letters, "|")
		// because it includes space in the last slice
		contents = contents[:len(contents)-1]

		for i, v := range contents {
			xb, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				log.Println(err)
			} else {
				contents[i] = string(xb)
			}
		}
	}

	err = tpl.ExecuteTemplate(w, "letterSent.html", contents)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) letterReceivedHandler(w http.ResponseWriter, r *http.Request) {
	email := context.Get(r, "email").(string)
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	var contents []string

	letters, err := h.letterUseCase.GetLettersByReceiverID(encodedEmail)
	if err != nil {
		log.Println(err)
		msg := url.QueryEscape("sorry, internal server error")
		http.Redirect(w, r, "/?"+msg, http.StatusSeeOther)
	}
	if letters != "" {
		contents = strings.Split(letters, "|")
		// because it includes space in the last slice
		contents = contents[:len(contents)-1]

		for i, v := range contents {
			xb, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				log.Println(err)
			} else {
				contents[i] = string(xb)
			}
		}
	}

	err = tpl.ExecuteTemplate(w, "letterReceived.html", contents)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) adminHandler(w http.ResponseWriter, r *http.Request) {
	letters, err := h.letterUseCase.GetAll()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = tpl.ExecuteTemplate(w, "admin.html", letters)
	if err != nil {
		log.Println("Error in WriteString: ", err)
	}
}
