package handler

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Env struct {
	Forum *database.Forum
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

type Ghuser struct {
	Login string `json:"login"`
	Id    int    `json:"id"`
}

type Guser struct {
	Name  string `json:"given_name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

type Email []struct {
	Payload struct {
		Commits []struct {
			Author struct {
				Email string `json:"email"`
			} `json:"author"`
		} `json:"commits"`
	} `json:"payload"`
}

func (env *Env) CheckCookie(w http.ResponseWriter, c *http.Cookie) []string {
	co := []string{}
	if strings.Contains(c.String(), "&") {
		co = strings.Split(c.Value, "&")
	}
	if len(co) != 0 {
		if !(env.Forum.CheckSession(co[1])) {
			// Set the new token as the users `session_token` cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now(),
			})

		} else {
			return co
		}
	}
	return co
}
func (env *Env) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
	type data struct {
		Cookie      interface{}
		CurrentUser interface{}
		Posts       interface{}
	}
	var page data
	filter := r.FormValue("filter")
	c, err := r.Cookie("session_token")
	if err != nil {
		page = data{Cookie: err.Error(), Posts: env.Forum.AllPost(filter, ""), CurrentUser: env.Forum.GetUser("")}
		if err := t.Execute(w, page); err != nil {
			http.Error(w, "500 Internal error", http.StatusInternalServerError)
			return
		}
	} else {
		co := env.CheckCookie(w, c)
		content := r.FormValue("comment")
		postID := r.FormValue("postID")
		like := r.FormValue("likes")
		dislike := r.FormValue("dislike")
		likesc := r.FormValue("likesc")
		dislikec := r.FormValue("dislikec")

		if content != "" {
			env.Forum.CreateComment(co[0], postID, content)
		}
		if like != "" {
			li := strings.Split(like, "&")
			env.Forum.UpdatePostReaction(li[1], co[0], li[0])
		}
		if dislike != "" {
			disl := strings.Split(dislike, "&")
			env.Forum.UpdatePostReaction(disl[1], co[0], disl[0])
		}
		if likesc != "" {
			lic := strings.Split(likesc, "&")
			env.Forum.UpdateCommentReaction(lic[1], lic[2], co[0], lic[0])
		}
		if dislikec != "" {
			dislc := strings.Split(dislikec, "&")
			env.Forum.UpdateCommentReaction(dislc[1], dislc[2], co[0], dislc[0])
		}
		yourPost := r.FormValue("yourPost")
		yourLikedPosts := r.FormValue("yourLikedPosts")
		if yourPost == "on" {
			page = data{Cookie: c.Value, Posts: env.Forum.AllPost(filter, co[0]), CurrentUser: env.Forum.GetUser(co[0])}
		} else if yourLikedPosts == "on" {
			page = data{Cookie: c.Value, Posts: env.Forum.YourLikedPost(co[0]), CurrentUser: env.Forum.GetUser(co[0])}
		} else {
			page = data{Cookie: c.Value, Posts: env.Forum.AllPost(filter, ""), CurrentUser: env.Forum.GetUser(co[0])}
		}
		if err := t.Execute(w, page); err != nil {
			http.Error(w, "500 Internal error", http.StatusInternalServerError)
			return
		}
	}

}

func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		userName := r.FormValue("userName")
		password := r.FormValue("password")
		if userName == "" || password == "" {
			http.Error(w, "400 Bad Request.", http.StatusBadRequest)
			return
		}
		UserID, Username, SessionID, err := env.Forum.LoginUsers(userName, r.UserAgent(), GetIP(r), password)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/text")
			w.Write([]byte("0" + err.Error()))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   UserID + "&" + SessionID + "&" + Username,
			Expires: time.Now().Add(24 * time.Hour),
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte("1" + UserID + "-" + SessionID + "-" + Username))
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
	}
}
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func (env *Env) Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		userName := r.FormValue("userName")
		password := r.FormValue("password")
		email := r.FormValue("email")
		if userName == "" || password == "" || email == "" {
			http.Error(w, "400 Bad Request.", http.StatusBadRequest)
			return
		}
		_, _, _, err := env.Forum.CreateUser(userName, email, r.UserAgent(), GetIP(r), password)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/text")
			w.Write([]byte("0" + err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte("1Register successful"))
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
	}

}

func (env *Env) Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "POST":
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := strings.Split(c.Value, "&")
		err = env.Forum.RemoveSession(res[1])
		if err != nil {
			log.Fatal(err)
		}
		// Set the new token as the users `session_token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte("Logout successful"))
		fmt.Println("Logout successful")
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
		return
	}
}

func (env *Env) Post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	c, err := r.Cookie("session_token")
	co := env.CheckCookie(w, c)
	if err == nil {
		switch r.Method {
		case "POST":
			categories := r.FormValue("categories")
			title := r.FormValue("title")
			post := r.FormValue("post")
			postID, _ := env.Forum.CreatePost(co[0], post, categories, title)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/text")
			w.Write([]byte(postID))
		default:
			http.Error(w, "400 Bad Request.", http.StatusBadRequest)
		}
	}
}

func (env *Env) Comment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	c, err := r.Cookie("session_token")
	co := env.CheckCookie(w, c)
	if err == nil {
		switch r.Method {
		case "POST":
			content := r.FormValue("comment")
			postID := r.FormValue("postID")
			commentID, _ := env.Forum.CreateComment(co[0], postID, content)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/text")
			w.Write([]byte(commentID))
		default:
			http.Error(w, "400 Bad Request.", http.StatusBadRequest)
		}
	}
}

func (env *Env) Redirected(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/callback" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	const clientID = "c298bb52526f90357763"
	const clientSecret = "66afb6c6f8ce0259a92799823210c2bfd2625e58"

fmt.Println(os.Getenv("GSecret"))
	httpClient := http.Client{}

	// First, we need to get the value of the `code` query param
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := r.FormValue("code")

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token

	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("Authorization", "token "+t.AccessToken)

	res, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()

	var user Ghuser

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	req, err = http.NewRequest("GET", "https://api.github.com/users/"+user.Login+"/events/public", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	res, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()

	var email Email

	if err := json.NewDecoder(res.Body).Decode(&email); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	UserID, Username, SessionID, err4 := env.Forum.OauthSigninOrRegister(user.Login, email[0].Payload.Commits[0].Author.Email, r.UserAgent(), GetIP(r), strconv.Itoa(user.Id))


	if err4 != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte("0" + err.Error()))
		return
	}

	w.Header().Set("Location", "/")

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   UserID + "&" + SessionID + "&" + Username,
		Expires: time.Now().Add(24 * time.Hour),
	})

		w.Header().Set("Location", "/?access="+UserID + "&" + SessionID + "&" + Username)
		w.WriteHeader(http.StatusFound)
}

func (env *Env) Redirected2(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/callback/2" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	const clientID = "1025088139209-li5k87h94rdp8cm1m81turmkucs7c2c0.apps.googleusercontent.com"
	const clientSecret = "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"
	httpClient := http.Client{}

	// First, we need to get the value of the `code` query param
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	code := r.FormValue("code")

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	reqURL := fmt.Sprintf("https://oauth2.googleapis.com/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=http://localhost:8800/login/callback/2", clientID, clientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// We set this header since we want the response
	// as JSON
	req.Header.Set("accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var g OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token

	req, err = http.NewRequest("GET", "https://www.googleapis.com/userinfo/v2/me", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+g.AccessToken)

	res, err = httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()

	var user Guser

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	UserID, Username, SessionID, err4 := env.Forum.OauthSigninOrRegister(user.Name, user.Email, r.UserAgent(), GetIP(r), user.ID)


	if err4 != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte("0" + err.Error()))
		return
	}
	

	w.Header().Set("Location", "/")

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   UserID + "&" + SessionID + "&" + Username,
		Expires: time.Now().Add(24 * time.Hour),
	})

		w.Header().Set("Location", "/?access="+UserID + "&" + SessionID + "&" + Username)
		w.WriteHeader(http.StatusFound)
}
