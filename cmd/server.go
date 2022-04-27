package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Forum struct {
	DB *sql.DB
}

type User struct {
	UserID   string
	Username string
	Email    string
	Password string
}
type Post struct {
	PostID   string
	UserID   string
	Content  string
	Category string
}

type Comment struct {
	CommentID string
	UserID    string
	PostID    string
	Content   string
}

type Reaction struct {
	ReactionID string
	PostID     string
	CommentID  string
	UserID     string
	React      int
}

func initUser(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "User" (
		"userID"	TEXT NOT NULL,
		"username"	TEXT,
		"email"	TEXT,
		"password"	TEXT,
		PRIMARY KEY("userID")
	);
	`)
	stmt.Exec()
}
func initPost(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Post" (
		"postID"	TEXT NOT NULL,
		"userID"	TEXT UNIQUE,
		"category"	TEXT,
		"content"	TEXT,
		PRIMARY KEY("postID")
	);
	`)
	stmt.Exec()
}
func initComment(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Comment" (
		"commentID" TEXT NOT NULL,
		"postID"	TEXT,
		"userID"	TEXT,
		"content"	TEXT,
		PRIMARY KEY("commentID")
	);
	`)
	stmt.Exec()
}
func initReaction(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Reaction" (
		"reactionID" TEXT,
		"postID"	TEXT,
		"commentID" TEXT,
		"userID"	TEXT,
		"react"	int,
		PRIMARY KEY("reactionID")
	);
	`)
	stmt.Exec()
}

func CreateDatabase(db *sql.DB) *Forum {
	initUser(db)
	initPost(db)
	initComment(db)
	initReaction(db)
	return &Forum{
		DB: db,
	}
}
func (forum *Forum) AddUser(user User) {
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO User (userID, username, email, password) values (?, ?, ?, ?)
	`)
	stmt.Exec(user.UserID, user.Username, user.Email, user.Password)
}
func (forum *Forum) CreatePost(post Post) {
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Post (postID, userID, content, category) values (?, ?, ?, ?)
	`)
	stmt.Exec(post.PostID, post.UserID, post.Content, post.Category)
}

func (forum *Forum) AddComment(comment Comment) {
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Comment (commentID, userID, postID, content) values (?, ?, ?, ?)
	`)
	stmt.Exec(comment.CommentID, comment.UserID, comment.PostID, comment.Content)
}

func (forum *Forum) ReactInPost(react Reaction) {
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Reaction (reactionID, postID, userID, react) values (?, ?, ?, ?)
	`)
	stmt.Exec(react.ReactionID, react.PostID, react.UserID, react.React)
}
func (forum *Forum) ReactInComment(react Reaction) {
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Reaction (reactionID, postID, commentID, userID, react) values (?, ?, ?, ?, ?)
	`)
	stmt.Exec(react.ReactionID, react.PostID, react.CommentID, react.UserID, react.React)
}

//HashPassword converts the password into hash values
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

//CheckPasswordHash checks if the password is match
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func PostProcessTrial(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
	}
	t.Execute(w, "Hello")
}

func main() {
	/*
	db, _ := sql.Open("sqlite3", "./forum.db")
	user := CreateDatabase(db)
	password, _ := HashPassword("hello")
	user.AddUser(User{
		UserID:   "lol1",
		Username: "Adriell",
		Email:    "hello@gmail.com",
		Password: password,
	})
	user.CreatePost(Post{
		PostID:   "1",
		UserID:   "lol1",
		Content:  "Hello, World",
		Category: "Programming",
	})
	user.AddComment(Comment{
		CommentID: "1",
		UserID:    "lol1",
		PostID:    "1",
		Content:   "HELLO, WORLD",
	})
	user.ReactInPost(Reaction{
		ReactionID: "1",
		UserID:     "lol1",
		PostID:     "1",
		React:      1,
	})
	user.ReactInComment(Reaction{
		ReactionID: "2",
		UserID:     "lol1",
		PostID:     "1",
		CommentID:  "1",
		React:      0,
	})
	*/
	http.HandleFunc("/", PostProcessTrial)
	cssPath := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", cssPath)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
