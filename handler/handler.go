package handler

import (
	"database/sql"
	"fmt"
	"forum/database"
	"net/http"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":

		userName := r.FormValue("userName")
		password := r.FormValue("password")

		db, _ := sql.Open("sqlite3", "./database/forum.db")
		Forum := database.CreateDatabase(db)
		UserID, Username, SessionID, _ := Forum.LoginUsers(userName, r.UserAgent(), GetIP(r), password)

		db.Close()

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   SessionID + "&" + UserID + "&" + Username,
			Expires: time.Now().Add(24 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte(UserID + "-" + Username + "-" + SessionID))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}

}
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		userName := r.FormValue("userName")
		password := r.FormValue("password")
		email := r.FormValue("email")

		db, _ := sql.Open("sqlite3", "./database/forum.db")
		Forum := database.CreateDatabase(db)

		UserID, Username, SessionID, _ := Forum.CreateUser(userName, email, r.UserAgent(), GetIP(r), password)

		db.Close()

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   SessionID + "&" + UserID + "&" + Username,
			Expires: time.Now().Add(24 * time.Hour),
		})

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte(UserID + SessionID + Username))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("err")
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

		db, _ := sql.Open("sqlite3", "./database/forum.db")
		Forum := database.CreateDatabase(db)

		Forum.RemoveSession(res[2])

		db.Close()

		// Set the new token as the users `session_token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		//w.Write([]byte(UserID + SessionID + Username))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}

}
