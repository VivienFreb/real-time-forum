// main.go
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			content TEXT,
			created_at TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}

func main() {
	initDB()
	defer db.Close()

	// http.HandleFunc("/", indexHandler)
	// http.Handle("/", http.FileServer(http.Dir('.')))
	http.Handle("/static/", http.StripPrefix("/static/",http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", ultimateHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/css")
        http.ServeFile(w, r, "styles.css")
    })
	// http.HandleFunc("/create", createHandler)
	// http.HandleFunc("/login", loginHandler)


	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ultimateHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodGet:
		http.ServeFile(w, r, "index.html")
	case http.MethodPost:
		// path := r.URL.Path
		// switch path{
		// case "/register":
		// 	registerHandler(w, r)
		// case "/login":
		// 	loginHandler(w,r)
		// }
		switch r.FormValue("formName"){
		case "register":
			registerHandler(w,r)
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting the parse!")
	//Parse form data???
	// Vérifier que tout se déroule sans encombre d'un POV technique?
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Ici je récupère les valeurs pour les mettre dans des variables
	username := r.Form.Get("username")
	email := r.Form.Get("mail")
	password := r.Form.Get("password")
	confirmPassword := r.Form.Get("confirm_password")

	log.Printf("Rec  eived form data - Username: %s, Email: %s, Password: %s, ConfirmPasword: %s", username, email, password, confirmPassword)

	if username == "" || email == "" || password == "" || confirmPassword == "" {
		http.Error(w, "Invalid Form Data", http.StatusBadRequest)
		return
	}

	// Verif du mail avec du RegEx
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if match, _ := regexp.MatchString(emailRegex, email); !match {
		http.Error(w, "Invalid Email Format", http.StatusBadRequest)
		return
	}
	// fmt.Println("Email templated!")

	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w,r,"/",http.StatusFound)
	fmt.Println("New user registered.")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Mauvaise méthode!", http.StatusMethodNotAllowed)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		http.Error(w, "Invalid Form Data", http.StatusBadRequest)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	_, err := db.Exec("INSERT INTO posts (title, content, created_at) VALUES (?, ?, ?)", title, content, time.Now())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
