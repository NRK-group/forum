package handler

import (
	"fmt"
	"forum/database"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Env struct {
	Forum *database.Forum
}

func (env *Env) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
	c, err := r.Cookie("session_token")
	co := []string{}
	if strings.Contains(c.String(), "&") {
		co = strings.Split(c.Value, "&")
	}
	// fmt.Println(co)
	if len(co) != 0 {
		if !(env.Forum.CheckSession(co[1])) {
			// Set the new token as the users `session_token` cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now(),
			})
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/text")
		}
	}
	if err != nil {
		// a, _ := fmt.Fprintf(w, "err")
		t.Execute(w, err.Error())
	} else {
		t.Execute(w, c.Value)
	}
}

func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		userName := r.FormValue("userName")
		password := r.FormValue("password")
		UserID, Username, SessionID, _ := env.Forum.LoginUsers(userName, r.UserAgent(), GetIP(r), password)

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   UserID + "&" + SessionID + "&" + Username,
			Expires: time.Now().Add(24 * time.Hour),
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte(UserID + "-" + SessionID + "-" + Username))
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
	}
}
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func (env *Env) Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		userName := r.FormValue("userName")
		password := r.FormValue("password")
		email := r.FormValue("email")
		if userName == "" || password == "" || email == "" {
			http.Error(w, "400 Bad Request.", http.StatusBadRequest)
			return
		}
		_, _, _, err := env.Forum.CreateUser(userName, email, r.UserAgent(), GetIP(r), password)
		// fmt.Println(err.Error())
		w.Write([]byte(err.Error()))
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
	}

}

func (env *Env) Logout(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := strings.Split(c.Value, "&")
		fmt.Println(res)
		err = env.Forum.RemoveSession(res[1])
		if err != nil {
			log.Fatal(err)
		}
		// Set the new token as the users `session_token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		//w.Write([]byte(UserID + SessionID + Username))
		fmt.Println("Logout successful")
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
		return
	}
}
