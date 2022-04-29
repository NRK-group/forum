package database

import (
	"fmt"
	"forum/password"
	"time"

	uuid "github.com/satori/go.uuid"
)

var Date = time.Now().Format("2006-January-01 15:04:05")

// start of the create query

// CreateUser
// is a method of database that add user in it.
func (forum *Forum) CreateUser(username, email, userAgent, ipAddress, pass string) (string, string, error) {
	userID := uuid.NewV4()
	pass, _ = password.HashPassword(pass)
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO User (userID, username, dateCreated, email, password) values (?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(userID, username, Date, email, pass)
	if err != nil {
		return "", "", err
	}
	sessionID, _ := forum.CreateSession(userID.String(), userAgent, ipAddress)
	return userID.String(), sessionID, nil
}

// CreateSession
// is a method of database that add session in it based on the user login time.
func (forum *Forum) CreateSession(userID, userAgent, ipAddress string) (string, error) {
	sessionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Session (sessionID, userID, userAgent, ipAddress, loginTime) values (?, ?, ?, ?, ?)
	`)
	// fmt.Print(err1)
	_, err := stmt.Exec(sessionID, userID, userAgent, ipAddress, Date)
	if err != nil {
		return "", err
	}
	forum.Update("User", "sessionID", sessionID.String(), "userID", userID)
	return sessionID.String(), nil
}

// CreatePost
// is a method of database that add post in it.
func (forum *Forum) CreatePost(userID, content, category string) (string, error) {
	postID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Post (postID, userID, dateCreated, content, category) values (?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(postID, userID, Date, content, category)
	if err != nil {
		return "", err
	}
	return postID.String(), nil
}

// CreateComment
// is a method of database that add comment in it.
func (forum *Forum) CreateComment(userID, postID, content string) (string, error) {
	commentID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Comment (commentID, userID, postID, dateCreated, content) values (?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(commentID, userID, postID, Date, content)
	if err != nil {
		return "", err
	}
	return commentID.String(), nil
}

// ReactInPost
// is a method of database that add reaction in the post in it.
func (forum *Forum) ReactInPost(postID, userID string, react int) (string, error) {
	reactionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Reaction (reactionID, postID, userID, react) values (?, ?, ?, ?)
	`)
	_, err := stmt.Exec(reactionID, postID, userID, react)
	if err != nil {
		return "", err
	}
	return reactionID.String(), nil
}

// ReactInComment
// is a method of database that add reaction in the comment in it.
func (forum *Forum) ReactInComment(postID, commentID, userID string, react int) (string, error) {
	reactionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Reaction (reactionID, postID, commentID, userID, react) values (?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(reactionID, postID, commentID, userID, react)
	if err != nil {
		return "", err
	}
	return reactionID.String(), nil
}

// start of the delete query

// Delete
// is a method of the database that delete value base on table and where.
//  ex. Forum.Delete("User", "userID", "185c6549-caec-4eae-95b0-e16023432ef0")
func (forum *Forum) Delete(table, where, value string) error {
	dlt := "DELETE FROM " + table + " WHERE " + where
	stmt, err1 := forum.DB.Prepare(dlt + " = (?)")
	fmt.Print(err1)
	_, err := stmt.Exec(value)
	if err != nil {
		fmt.Println("lol")
		return err
	}
	// fmt.Print(stmt)
	return nil
}

// RemoveSession
func (forum *Forum) RemoveSession(userID, sessionID string) error {
	err := forum.Update("User", "sessionID", "", "userID", userID)
	if err != nil {
		return err
	}
	err = forum.Delete("Session", "sessionID", sessionID)
	if err != nil {
		return err
	}
	return nil
}

// Start of the Update query

// Update
//  ex. Forum.Update("User", "username", "Adriell", "userID" "7e2b4fdd-86ad-464c-a97e")
func (forum *Forum) Update(table, set, to, where, id string) error {
	update := "UPDATE " + table + " SET " + set + " = '" + to + "' WHERE " + where + " = '" + id + "'"
	stmt, _ := forum.DB.Prepare(update)
	_, err := stmt.Exec()
	if err != nil {
		return err
	}
	// fmt.Print(stmt)
	return nil
}

// Start of the Select query
