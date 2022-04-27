package main

import (
	"database/sql"
	"fmt"
	"forum/database"
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
	db, _ := sql.Open("sqlite3", "./database/forum.db")
	Forum := database.CreateDatabase(db)
	r, err := Forum.AddUser("sadasd", "hasd@gmail.com", "hello1235")
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(r)
	}
	// err1 := Forum.Delete("sad")
	// if err1 != nil {
	// 	fmt.Print(err1)
	// }
	defer db.Close()
	http.HandleFunc("/", PostProcessTrial)
	cssPath := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", cssPath)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
