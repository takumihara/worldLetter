package main

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	"github.com/tacomea/worldLetter/domain"
	"github.com/tacomea/worldLetter/token"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type handler struct {
	userUseCase    domain.UserUseCase
	sessionUseCase domain.SessionUseCase
	letterUseCase  domain.LetterUseCase
}

func newHandler(uu domain.UserUseCase, su domain.SessionUseCase, lu domain.LetterUseCase) *handler {
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

func (h *handler) indexHandler(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
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

func (h *handler) createHandler(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) sendHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("letter")
	id := uuid.NewString()
	email := context.Get(r, "email").(string)
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	letter := domain.Letter{
		ID:       id,
		AuthorID: encodedEmail,
		Content:  content,
		IsSent:   false,
	}

	err := h.letterUseCase.Create(letter)
	if err != nil {
		log.Println(err)
		query := url.QueryEscape("sorry, internal server error")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}

	user, err := h.userUseCase.Read(encodedEmail)
	if err != nil {
		log.Println(err)
		query := url.QueryEscape("sorry, internal server error")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}

	err = h.userUseCase.Update(user)
	if err != nil {
		log.Println(err)
		query := url.QueryEscape("sorry, internal server error")
		http.Redirect(w, r, "/?msg="+query, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/show", http.StatusSeeOther)
}

func (h *handler) showHandler(w http.ResponseWriter, r *http.Request) {
	email := context.Get(r, "email").(string)
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	letter, err := h.letterUseCase.GetFirstUnsendLetter(encodedEmail)
	if err != nil {
		log.Println(err)
	} else if letter.Content == "" {
		letter.Content = "sorry, letter was not retrieved. Your letter might be out first one!"
	} else {
		letter.IsSent = true
		letter.ReceiverID = encodedEmail
		err := h.letterUseCase.Update(letter)
		if err != nil {
			log.Println(err)
		}
	}

	err = tpl.ExecuteTemplate(w, "show.html", letter.Content)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) letterSentHandler(w http.ResponseWriter, r *http.Request) {
	email := context.Get(r, "email").(string)
	encodedEmail := base64.StdEncoding.EncodeToString([]byte(email))

	contents := make([]string, 0, 5)

	letters, err := h.letterUseCase.GetLettersByAuthorID(encodedEmail)
	if err != nil {
		log.Println(err)
	} else if letters == "" {
		content := "I guess you haven't sent any letter yet."
		contents = append(contents, content)
	} else {
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

	contents := make([]string, 0, 5)

	letters, err := h.letterUseCase.GetLettersByReceiverID(encodedEmail)
	if err != nil {
		log.Println(err)
		msg := url.QueryEscape("sorry, internal server error")
		http.Redirect(w, r, "/?"+msg, http.StatusSeeOther)
	}
	if letters == "" {
		content := "I guess you haven't received any letter yet."
		contents = append(contents, content)
	} else {
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
