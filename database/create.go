package database

import (
	"database/sql"
	"forum/password"
	"time"

	uuid "github.com/satori/go.uuid"
)

var Now = time.Now().Format("2022-January-01")

func (forum *Forum) AddUser(username, email, pass string) (sql.Result, error) {
	userID := uuid.NewV4()
	pass, _ = password.HashPassword(pass)
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO User (userID, username, dateCreated, email, password) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(userID, username, Now, email, pass)
	if err != nil {
		return r, err
	}
	return r, nil
}

// func (forum *Forum) Delete(value string) error {
// 	// dlt := "DELETE FROM " + table + " WHERE " + where + " = " + value
// 	stmt, _ := forum.DB.Prepare(`
// 		DELETE FROM User WHERE (username)
// 	`)
// 	_, err := stmt.Exec()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (forum *Forum) CreateSession(userID, userAgent, ipAddress string) (sql.Result, error) {
	sessionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO User (sessionID, userID, userAgent, ipAddress, Now) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(sessionID, userID, userAgent, ipAddress, Now)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (forum *Forum) CreatePost(userID, content, category string) (sql.Result, error) {
	postID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Post (postID, userID, dateCreated, content, category) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(postID, userID, Now, content, category)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (forum *Forum) AddComment(userID, postID, content string) (sql.Result, error) {
	commentID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Comment (commentID, userID, postID, dateCreated, content) values (?, ?, ?, ?, ?)
	`)
	r, err := stmt.Exec(commentID, userID, postID, Now, content)
	if err != nil {
		return r, err
	}
	return r, nil
}

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
