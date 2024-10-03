package handler

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/argon2"
	"johtotimes.com/src/assert"
	"johtotimes.com/src/database"
	"johtotimes.com/src/model"
	"johtotimes.com/src/templates"
)

var user *model.User

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ok = matchCredentials(username, password)

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func CookieAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		log.Println("AdminHandler: Checking cookie")
		cookie, err := req.Cookie("user")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Println("AdminHandler: No cookie")
				w.WriteHeader(http.StatusBadRequest)
				// w.Write([]byte("No cookie"))
			} else {
				log.Println("AdminHandler: Error reading cookie")
				w.WriteHeader(http.StatusInternalServerError)
				// w.Write([]byte("Error reading cookie"))
			}
			LoginPageHandler(w, req)
			return
		}

		db := database.Connect()
		defer db.Close()
		user, err = db.Users.GetByEmail(cookie.Value)
		assert.NoError(err, "Error getting user by email")

		next.ServeHTTP(w, req)

	})
}

func LoginPageHandler(w http.ResponseWriter, req *http.Request) {
	loginPage := templates.LoginTemplate()
	render(loginPage, isHTMX(req), "Login", w)
}

func LoginRequestHandler(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("email")
	password := req.FormValue("password")
	log.Printf("Username: %s, Password: %s", username, password)

	if !matchCredentials(username, password) {
		// w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		w.Write([]byte("Login failed"))
		return
	}

	cookie := http.Cookie{
		Name:     "user",
		Value:    username,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &cookie)

	// w.Write([]byte("Login successful"))
	CookieAuth(AdminHandler)(w, req)
}

func AdminEditorHandler(w http.ResponseWriter, req *http.Request) {
	postID, err := strconv.ParseInt(req.PathValue("postID"), 10, 64)
	assert.NoError(err, "AdminHandler: Error parsing postID")

	db := database.Connect()
	defer db.Close()
	post, err := db.Posts.GetByID(postID)
	assert.NoError(err, "AdminHandler: Error getting post by ID")
	log.Printf("AdminHandler: Post: %+v", post)

	postEditor := templates.PostEditorTemplate(*post)
	render(postEditor, isHTMX(req), "Post Editor", w)
}

func AdminHandler(w http.ResponseWriter, req *http.Request) {
	db := database.Connect()
	defer db.Close()

	posts := db.Posts.GetPage(model.ISSUE, 1, 10)
	log.Printf("Found %d posts", len(posts))
	log.Printf("User: %+v", user)
	adminPage := templates.AdminTemplate(user.Name, posts)
	render(adminPage, isHTMX(req), "Admin", w)
}

func matchCredentials(username, password string) bool {
	db := database.Connect()
	defer db.Close()
	user, err := db.Users.GetByEmail(username)
	if err != nil {
		log.Printf("Error getting user by email: %s", err)
		return false
	}

	passwordArgon := argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32)

	hashUsername := sha256.Sum256([]byte(username))
	hashPassword := sha256.Sum256([]byte(passwordArgon))

	expectedUsername := sha256.Sum256([]byte(user.Email))
	expectedPassword := sha256.Sum256([]byte(user.Password))

	usernameMatch := (subtle.ConstantTimeCompare(hashUsername[:], expectedUsername[:]) == 1)
	passwordMatch := (subtle.ConstantTimeCompare(hashPassword[:], expectedPassword[:]) == 1)
	return usernameMatch && passwordMatch
}
