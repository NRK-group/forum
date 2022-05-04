package main

import (
	"fmt"
	"forum/handler"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}

	c, err := r.Cookie("session_token")
if err != nil {
	t.Execute(w, "Hello")
} else {
t.Execute(w, c.Value)
}
}

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/logout", handler.Logout)
	cssPath := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", cssPath)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
