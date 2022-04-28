package database

import "database/sql"

type Forum struct {
	DB *sql.DB
}

type User struct {
	UserID      string
	Username    string
	Email       string
	DateCreated string
	Password    string
}
type Session struct {
	SessionID string
	UserID    string
	UserAgent string
	IpAddress string
	LoginTime string
}
type Post struct {
	PostID      string
	UserID      string
	DateCreated string
	Content     string
	Category    string
}

type Comment struct {
	CommentID   string
	UserID      string
	PostID      string
	DateCreated string
	Content     string
}

type Reaction struct {
	ReactionID string
	PostID     string
	CommentID  string
	UserID     string
	React      int
}
