package handler

import (
	"database/sql"
	"fmt"
	"forum/database"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("err")
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
		UserID, Username, SessionID, _ := Forum.LoginUsers(userName, "userAgent", "ipAddress", password)

		db.Close()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte(UserID + "-" + Username + "-" + SessionID))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}

}


func Register(w http.ResponseWriter, r *http.Request ) {


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

	UserID, Username, SessionID, _:= Forum.CreateUser(userName, email, "userAgent", "ipAddress", password)

	db.Close()
		
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte(UserID + SessionID + Username ))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}


}
