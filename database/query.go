package database

import (
	"database/sql"
	"fmt"
	"forum/password"
	"time"

	uuid "github.com/satori/go.uuid"
)

var Date = time.Now().Format("2022-January-01")

// start of the create query

// CreateUser
// is a method of database that add user in it.
func (forum *Forum) CreateUser(username, email, pass string) (sql.Result, error) {
	userID := uuid.NewV4()
	pass, _ = password.HashPassword(pass)
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO User (userID, username, dateCreated, email, password) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(userID, username, Date, email, pass)
	if err != nil {
		return r, err
	}
	return r, nil
}

// CreateSession
// is a method of database that add session in it based on the user login time.
func (forum *Forum) CreateSession(userID, userAgent, ipAddress string) (sql.Result, error) {
	sessionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO User (sessionID, userID, userAgent, ipAddress, Now) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(sessionID, userID, userAgent, ipAddress, Date)
	if err != nil {
		return r, err
	}
	return r, nil
}

// CreatePost
// is a method of database that add post in it.
func (forum *Forum) CreatePost(userID, content, category string) (sql.Result, error) {
	postID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Post (postID, userID, dateCreated, content, category) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(postID, userID, Date, content, category)
	if err != nil {
		return r, err
	}
	return r, nil
}

// CreateComment
// is a method of database that add comment in it.
func (forum *Forum) CreateComment(userID, postID, content string) (sql.Result, error) {
	commentID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Comment (commentID, userID, postID, dateCreated, content) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(commentID, userID, postID, Date, content)
	if err != nil {
		return r, err
	}
	return r, nil
}

// ReactInPost
// is a method of database that add reaction in the post in it.
func (forum *Forum) ReactInPost(postID, userID string, react int) (sql.Result, error) {
	reactionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Reaction (reactionID, postID, userID, react) values (?, ?, ?, ?)
	`)
	r, err := stmt.Exec(reactionID, postID, userID, react)
	if err != nil {
		return r, err
	}
	return r, nil
}

// ReactInComment
// is a method of database that add reaction in the comment in it.
func (forum *Forum) ReactInComment(postID, commentID, userID string, react int) (sql.Result, error) {
	reactionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Reaction (reactionID, postID, commentID, userID, react) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(reactionID, postID, commentID, userID, react)
	if err != nil {
		return r, err
	}
	return r, nil
}

// start of the delete query

// Delete
// is a method of the database that delete value base on table and where.
//  ex. "DELETE FROM table WHERE where = value"
func (forum *Forum) Delete(table, where, value string) error {
	dlt := "DELETE FROM " + table + " WHERE " + where
	stmt, _ := forum.DB.Prepare(dlt + "= (?)")
	r, err := stmt.Exec(value)
	if err != nil {
		return err
	}
	fmt.Print(r)
	return nil
}

// Start of the Edit query
