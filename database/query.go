package database

import (
	"errors"
	"fmt"
	"forum/password"
	"time"

	uuid "github.com/satori/go.uuid"
)

var Date = time.Now().Format("2006-January-01 15:04:05")

// start of the create query

// CreateUser
// is a method of database that add user in it.
func (forum *Forum) CreateUser(username, email, userAgent, ipAddress, pass string) (string, string, string, error) {
	userID := uuid.NewV4()
	pass, _ = password.HashPassword(pass)
	stmt, err := forum.DB.Prepare(`
		INSERT INTO User (userID, username, dateCreated, email, password) values (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return "", "", "", err
	}
	_, err = stmt.Exec(userID, username, Date, email, pass)
	if err != nil {
		return "", "", "", err
	}
	// sessionID, err := forum.CreateSession(userID.String(), userAgent, ipAddress)
	// if err != nil {
	// 	return "", "", "", err
	// }
	forum.Update("User", "sessionID", "", "userID", userID.String())
	return userID.String(), username, "", nil
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
func (forum *Forum) CreatePost(userID, content, category, title string) (string, error) {
	postID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Post (postID, userID, dateCreated, content, category, title) values (?, ?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(postID, userID, Date, content, category, title)
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
	stmt, err := forum.DB.Prepare(dlt + " = (?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(value)
	if err != nil {
		return err
	}
	// fmt.Print(stmt)
	return nil
}

// RemoveSession
func (forum *Forum) RemoveSession(sessionID string) error {
	err := forum.Update("User", "sessionID", "", "sessionID", sessionID)
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

func (forum *Forum) LoginUsers(userName, userAgent, ipAddress, pas string) (string, string, string, error) {
	var users User
	rows, err := forum.DB.Query("SELECT * FROM User WHERE username = '" + userName + "'")
	if err != nil {
		return "", "", "", err
	}
	var userID, sessionID, username, email, dateCreated, pass string
	for rows.Next() {
		// fmt.Println(rows, "lol")s
		rows.Scan(&userID, &username, &email, &dateCreated, &pass, &sessionID)
		users = User{
			UserID:      userID,
			SessionID:   sessionID,
			Username:    username,
			Email:       email,
			DateCreated: dateCreated,
			Password:    pass,
		}
	}
	if users.Username == "" {
		return "", "", "", errors.New("user not found")
	}
	if !(password.CheckPasswordHash(pas, users.Password)) {
		return "", "", "", errors.New("password not macth")
	}
	if users.SessionID != "" {
		forum.RemoveSession(users.SessionID)
	}
	sess, err := forum.CreateSession(users.UserID, userAgent, ipAddress)
	if err != nil {
		return "", "", "", err
	}
	users.SessionID = sess
	return users.UserID, users.Username, users.SessionID, nil
}

//CheckSession
//is a method of forum that checks the session in the User table and checks if it is match with the sessionID inputed
func (forum *Forum) CheckSession(sessionId string) bool {
	rows, err := forum.DB.Query("SELECT sessionID FROM User WHERE sessionID = '" + sessionId + "'")
	if err != nil {
		fmt.Print(err)
		return false
	}
	var session string
	for rows.Next() {
		rows.Scan(&session)
	}
	return session != ""
}

//AllPost
//is a method of forum that will return all post
func (forum *Forum) AllPost(userName string)[]Post {
	rows, err := forum.DB.Query("SELECT * FROM Post ")
	var post Post
	var posts []Post

	if err != nil {
		fmt.Print(err)
		return posts
	}

	for rows.Next() {
		var col1, col2, col3, col4, col5, col6 string
		rows.Scan(&col1, &col2, & col3, & col4, & col5, &col6)
		post=Post{PostID:      col1,
			UserID:      userName,
			DateCreated: col4,
			Content:     col6,
			Category:    col4,
			Title: col3}
		posts =	append(posts, post)
	}
	return posts
}
