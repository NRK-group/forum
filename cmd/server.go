package main

import (
	"fmt"
	"forum/handler"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func PostProcessTrial(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
	t.Execute(w, "Hello")
}

func main() {

	http.HandleFunc("/", PostProcessTrial)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	cssPath := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", cssPath)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
