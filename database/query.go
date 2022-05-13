package database

import (
	"errors"
	"fmt"
	"forum/password"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var Date = time.Now().Format("2006 January 02 15:04:05")

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
// Ex. Forum.Forum.ReactInPost("b081d711-aad2-4f90-acea-2f2842e28512", "b53124c2-39f0-4f10-8e02-b7244b406b86", 1)
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
// is method of forum that removes the session based on the sessionID
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

//LoginUser
//is method of forum that checks the database if the login details match the credential
// and allow them to login if their is a match credentials
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

//CheckReactInPost
func (forum *Forum) CheckReactInPost(pID, uID string) (string, int) {
	rows, err := forum.DB.Query("SELECT reactionID, postID, userID, react FROM Reaction WHERE postID = '" + pID + "' AND userID = '" + uID + "' AND commentID IS NULL")
	var reaction Reaction
	if err != nil {
		fmt.Print(err)
		return "", 0
	}
	var reactionID, postID, userID string
	var react int
	for rows.Next() {
		rows.Scan(&reactionID, &postID, &userID, &react)
		reaction = Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			UserID:     userID,
			React:      react,
		}
	}
	// fmt.Print(reactions)
	return reactionID, reaction.React
}

//UpdatePostReaction
//	Ex. Forum.Forum.UpdatePostReaction("b081d711-aad2-4f90-acea-2f2842e28512", "b53124c2-39f0-4f10-8e02-b7244b406b86", "-1")
func (forum *Forum) UpdatePostReaction(pID, uID, value string) {
	rID, v := forum.CheckReactInPost(pID, uID)
	i, _ := strconv.Atoi(value)
	if v == 0 {
		forum.ReactInPost(pID, uID, i)
	} else if v == i {
		forum.Delete("Reaction", "reactionID", rID)
	} else {
		forum.Update("Reaction", "react", value, "reactionID", rID)
	}
}

//CheckReactInPost
func (forum *Forum) CheckReactInComment(cID, uID string) (string, int) {
	rows, err := forum.DB.Query("SELECT * FROM Reaction WHERE commentID = '" + cID + "' AND userID = '" + uID + "'")
	var reaction Reaction
	if err != nil {
		fmt.Print(err)
		return "", 0
	}
	var reactionID, postID, commentID, userID string
	var react int
	for rows.Next() {
		rows.Scan(&reactionID, &postID, &commentID, &userID, &react)
		reaction = Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			CommentID:  commentID,
			UserID:     userID,
			React:      react,
		}
	}
	return reactionID, reaction.React
}

//UpdatePostReaction
//	Ex. Forum.Forum.UpdateCommentReaction("98e63b80-d1de-40c2-b4c7-988249c5c60b", "c3dd7cc7-46a1-4d57-b6fb-378f0741d461", "b53124c2-39f0-4f10-8e02-b7244b406b86", "1")
func (forum *Forum) UpdateCommentReaction(cID, pID, uID, value string) {
	rID, v := forum.CheckReactInComment(cID, uID)
	i, _ := strconv.Atoi(value)
	if v == 0 {
		i, _ := strconv.Atoi(value)
		forum.ReactInComment(pID, cID, uID, i)
	} else if v == i {
		forum.Delete("Reaction", "reactionID", rID)
	} else {
		forum.Update("Reaction", "react", value, "reactionID", rID)
	}
}

//AllPost
//is a method of forum that will return all post
func (forum *Forum) AllPost(filter string) []Post {
	rows, err := forum.DB.Query("SELECT * FROM Post ")
	var post Post
	var posts []Post
	if err != nil {
		fmt.Print(err)
		return posts
	}
	var postID, userID, title, category, dateCreated, content string
	for rows.Next() {
		rows.Scan(&postID, &userID, &title, &category, &dateCreated, &content)
		post = Post{
			PostID:       postID,
			DateCreated:  dateCreated,
			Content:      content,
			Category:     category,
			Title:        title,
			Comments:     forum.GetComments(postID),
			NumOfComment: len(forum.GetComments(postID)),
			Reaction:     forum.GetReactionsInPost(postID),
		}
		var username string
		rows2, err := forum.DB.Query("SELECT username FROM User WHERE userID = '" + userID + "'")
		if err != nil {
			fmt.Print(err)
			return posts
		}
		for rows2.Next() {
			rows2.Scan(&username)
		}
		post.UserID = username
		switch filter {
		case "go":
			if strings.Contains(category, "Go") {
				posts = append([]Post{post}, posts...)
			}
		case "javascript":
			if strings.Contains(category, "javascript") {
				posts = append([]Post{post}, posts...)
			}
		case "rust":
			if strings.Contains(category, "rust") {
				posts = append([]Post{post}, posts...)
			}
		default:
			posts = append([]Post{post}, posts...)
		}
	}

	return posts
}

//YourPost
//is a method of forum that return all the post base on the userID
func (forum *Forum) YourPost(filter, uID string) []Post {
	rows, err := forum.DB.Query("SELECT * FROM Post WHERE userID = '" + uID + "'")
	var post Post
	var posts []Post
	if err != nil {
		fmt.Print(err)
		return posts
	}
	var postID, userID, title, category, dateCreated, content string
	for rows.Next() {
		rows.Scan(&postID, &userID, &title, &category, &dateCreated, &content)
		post = Post{
			PostID:      postID,
			DateCreated: dateCreated,
			Content:     content,
			Category:    category,
			Title:       title,
		}
		if filter == "old" {
			posts = append(posts, post)
		} else {
			posts = append([]Post{post}, posts...)
		}

	}
	return posts
}

//Get Comments
//is a method of forum that return all the comment with that specific postID
func (forum *Forum) GetComments(pID string) []Comment {
	rows, err := forum.DB.Query("SELECT * FROM Comment WHERE postID = '" + pID + "'")
	var comment Comment
	var comments []Comment
	if err != nil {
		fmt.Print(err)
		return comments
	}
	var commentID, postID, userID, dateCreated, content string
	for rows.Next() {
		rows.Scan(&commentID, &postID, &userID, &dateCreated, &content)
		comment = Comment{
			CommentID:   commentID,
			PostID:      postID,
			UserID:      userID,
			DateCreated: dateCreated,
			Content:     content,
			Reaction:    forum.GetReactionsInComment(commentID),
		}
		var username string
		rows2, err := forum.DB.Query("SELECT username FROM User WHERE userID = '" + userID + "'")
		if err != nil {
			fmt.Print(err)
			return comments
		}
		for rows2.Next() {
			rows2.Scan(&username)
		}
		comment.UserID = username
		comments = append([]Comment{comment}, comments...)
	}
	return comments
}

func (forum *Forum) GetReactionsInPost(pID string) Reaction {
	rows, err := forum.DB.Query("SELECT reactionID, postID, userID, react FROM Reaction WHERE postID = '" + pID + "' AND commentID IS NULL")
	var reaction Reaction
	// var reactions []Reaction
	if err != nil {
		fmt.Print(err)
		return reaction
	}
	var reactionID, postID, userID string
	var react, likes, dislikes int
	for rows.Next() {
		rows.Scan(&reactionID, &postID, &userID, &react)
		reaction = Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			UserID:     userID,
			React:      react,
		}
		if react == 1 {
			likes++
		}
		if react == -1 {
			dislikes++
		}
		reaction.Likes = likes
		reaction.Dislikes = dislikes
		// reactions = append(reactions, reaction)
	}
	return reaction
}

func (forum *Forum) GetReactionsInComment(cID string) Reaction {
	rows, err := forum.DB.Query("SELECT * FROM Reaction WHERE commentID = '" + cID + "'")
	var reaction Reaction
	if err != nil {
		fmt.Print(err)
		return reaction
	}
	var reactionID, postID, commentID, userID string
	var react, likes, dislikes int
	for rows.Next() {
		rows.Scan(&reactionID, &postID, &commentID, &userID, &react)
		reaction = Reaction{
			ReactionID: reactionID,
			PostID:     postID,
			CommentID:  commentID,
			UserID:     userID,
			React:      react,
		}
		if react == 1 {
			likes++
		}
		if react == -1 {
			dislikes++
		}
		reaction.Likes = likes
		reaction.Dislikes = dislikes
	}
	return reaction
}
