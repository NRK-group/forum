package database

import (
	"database/sql"
)

func initUser(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "User" (
		"userID"	TEXT NOT NULL,
		"username"	CHARACTER(20) UNIQUE NOT NULL,
		"email"	TEXT UNIQUE NOT NULL,
		"dateCreated" TEXT NOT NULL,
		"password"	TEXT NOT NULL,
		"sessionID" TEXT,
		PRIMARY KEY("userID")
		FOREIGN KEY ("sessionID")
			REFERENCES "Session" ("sessionID")
		CHECK (length("username") >= 3 AND length("username") <= 20 )
		CHECK (("email") LIKE '%_@__%.__%')
		CHECk (length("password") >= 8)
	);
	`)
	stmt.Exec()
}
func initSession(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Session" (
		"sessionID"	TEXT UNIQUE NOT NULL,
		"userAgent" TEXT NOT NULL,
		"ipAddress" TEXT NOT NULL,
		"loginTime" TEXT NOT NULL,
		"userID" TEXT NOT NULL,
		PRIMARY KEY("sessionID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}
func initPost(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Post" (
		"postID"	TEXT UNIQUE NOT NULL,
		"userID"	TEXT NOT NULL,
		"title"     TEXT NOT NULL,
		"category"	TEXT NOT NULL,
		"dateCreated" TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		PRIMARY KEY("postID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}
func initComment(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Comment" (
		"commentID" TEXT UNIQUE NOT NULL,
		"postID"	TEXT NOT NULL,
		"userID"	TEXT NOT NULL,
		"dateCreated" TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		PRIMARY KEY("commentID")
		FOREIGN KEY ("postID")
			REFERENCES "Post" ("postID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}
func initReaction(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Reaction" (
		"reactionID" TEXT NOT NULL,
		"postID"	TEXT NOT NULL,
		"commentID" TEXT,
		"userID"	TEXT NOT NULL,
		"react"	int,
		PRIMARY KEY("reactionID")
		FOREIGN KEY ("postID")
			REFERENCES "Post" ("postID")
		FOREIGN KEY ("commentID")
			REFERENCES "Comment" ("commentID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func CreateDatabase(db *sql.DB) *Forum {
	initUser(db)
	initSession(db)
	initPost(db)
	initComment(db)
	initReaction(db)
	return &Forum{
		DB: db,
	}
}
