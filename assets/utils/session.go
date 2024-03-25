package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	s "real/assets/struct"
)

func SessionExists(req *http.Request, CookieName string) bool {
	cookie, err := req.Cookie(CookieName)
	if err == http.ErrNoCookie {
		return false
	} else if err != nil {
		log.Println(err)
		return false
	}
	if _, exists := s.Sessions[cookie.Value]; !exists {
		return false
	}
	return true
}

// Delete a cookie from the Session, WITH LogoutHandler
func ClearSession(w http.ResponseWriter, CookieName string) {
	Cookie := &http.Cookie{
		Name:  CookieName,
		Value: "",
		// Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, Cookie)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		NewPassword := r.FormValue("NewPassword")
		fmt.Println("new password", NewPassword)
		// UPDATE Users SET `password` = "nope" WHERE `id` = Cookie
	}
}

func Connect(w http.ResponseWriter, r *http.Request, db *sql.DB, pseudo string, password string) (bool, string, string) {
	var Strid_user string
	var count int = 0
	CorrectPassword, _ := ComparePassword(db, pseudo, password)

	if CorrectPassword {
		rows, err := db.Query("SELECT id FROM Users WHERE pseudo= ?;", pseudo)
		if err != nil {
			PrintError("Request Select for connection failed", err, -1)
		}

		for rows.Next() {
			var user s.User
			rows.Scan(&user.ID)

			if user.ID != 0 {
				Strid_user = strconv.Itoa(user.ID)
			}
			count++
		}
	}

	if count == 0 {
		fmt.Println("unconnected")
		return true, "Identifiant ou mot de passe incorrect!", ""
	} else if count == 1 {
		fmt.Println("User connected: " + pseudo)
		s.Sessions[Strid_user] = true
		Cookie := &http.Cookie{
			Name:   "Connexion",
			Value:  Strid_user,
			MaxAge: 0,
		}
		r.AddCookie(Cookie)
		http.SetCookie(w, Cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		return false, "", ""
	}
	return true, "Identifiant ou mot de passe incorrect!", ""
}
