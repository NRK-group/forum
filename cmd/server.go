package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	cssPath := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", cssPath)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
