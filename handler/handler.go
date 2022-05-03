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
		userId, userN, _ := Forum.LoginUsers(userName, "userAgent", "ipAddress", password)

		db.Close()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte(userId + "-" + userN))

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}

}
