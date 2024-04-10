package utils

import (
	"database/sql"
	"fmt"

	s "real/assets/struct"

	_ "github.com/mattn/go-sqlite3"
)

// Get All post from database

// Get One post from database, thanks to Connexion Cookie --> id_user

// Insert a New Post on the database
func InsertPost(db *sql.DB, title string, text string, user int, date string, img string) (int, error) {
	var postId int = -1 // initialized to -1 for error
	_, err := db.Exec("INSERT INTO Posts(title,description,date,user_id,imgUrl) VALUES (?,?,?,?,?);", title, text, date, user, img)
	if err != nil {
		return postId, err
	}
	// execute query to get all posts from database
	rows, err := db.Query("SELECT id FROM Posts WHERE  date = ? AND user_id = ?", date, user)
	if err != nil {
		fmt.Errorf("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		err = rows.Scan(&postId)
		if err != nil {
			fmt.Errorf("Error Row:", err, 2)
		}
	}
	return postId, nil
}

// Insert a Category to a post on the database
func InsertPostCategory(db *sql.DB, post int, category string) error {
	var catId int = -1 // initialized to -1 for error
	// execute query to get all categories from database
	rows, err := db.Query("SELECT id FROM Categories WHERE  name = ? ", category)
	if err != nil {
		fmt.Errorf("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all categories rows get from database
	for rows.Next() {
		err = rows.Scan(&catId)
		if err != nil {
			fmt.Errorf("Error Row:", err, 2)
			return err
		}
	}
	if catId != -1 {
		_, err := db.Exec("INSERT INTO Categorized(post_id,category_id) VALUES (?,?);", post, catId)
		if err != nil {
			return err
		}
	}
	return nil
}

// GET all comment linked to a post
func SelectAllComments(db *sql.DB, post_id int) []*s.Comment {
	var Comments []*s.Comment

	rows, err := db.Query("SELECT c.id, u.id,u.Pseudo, c.message, c.date, c.comment_id FROM Comments c INNER JOIN Posts p ON c.post_id = p.id INNER JOIN Users u ON c.user_id = u.id WHERE post_id = ?;", post_id)
	if err != nil {
		fmt.Errorf("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		var comment s.Comment
		var reply sql.NullInt64 // Use sql.NullInt64 to handle NULL values in the database
		err = rows.Scan(&comment.ID, &comment.User_id, &comment.Pseudo, &comment.Message, &reply)
		if err != nil {
			fmt.Errorf("Error Row:", err, 2)
		}
		// Check if reply is not NULL
		if reply.Valid {
			commentID := int(reply.Int64)
			tComment := GetComment(Comments, commentID)
			tComment.SubComments = append(tComment.SubComments, &comment)
		} else {
			Comments = append(Comments, &comment)
		}
	}
	return Comments
}

// Find the pseudo of a User with the given id (Cookie)
func SelectPseudo(db *sql.DB, user_id int) string {
	// Variables declaration
	var pseudo string
	// execute query to get all posts from database
	rows, err := db.Query("SELECT pseudo FROM Users WHERE id = ?", user_id)
	if err != nil {
		fmt.Errorf("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&pseudo)
		if err != nil {
			fmt.Errorf("Error Row:", err, 2)
		}
	}
	return pseudo
}

// Find the pseudo of a User with the given id (Cookie)
func SelectProfilUser(db *sql.DB, user_id int) (string, string) {
	// Variables declaration
	var pseudo string
	var email string
	// execute query to get all posts from database
	rows, err := db.Query("SELECT pseudo, email FROM Users WHERE id = ?", user_id)
	if err != nil {
		fmt.Errorf("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&pseudo, &email)
		if err != nil {
			fmt.Errorf("Error Row:", err, 2)
		}
	}
	return pseudo, email
}

func SelectIdPost(db *sql.DB, post_id int) int {
	// Variables declaration
	var id_pseudo int
	// execute query to get all posts from database
	rows, err := db.Query("SELECT user_id FROM Posts WHERE id = ?", post_id)
	if err != nil {
		fmt.Errorf("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&id_pseudo)
		if err != nil {
			fmt.Errorf("Error Row:", err, 2)
		}
	}
	return id_pseudo
}

func InsertComment(db *sql.DB, Comment string, date string, id_user int, id_post int) {
	db.Exec("INSERT INTO Comments(message,date,user_id, post_id) VALUES (?,?,?,?)", Comment, date, id_user, id_post)
}

func InsertReply(db *sql.DB, Comment string, date string, id_user int, id_post int, id_comment int) {
	db.Exec("INSERT INTO Comments(message,date,user_id, post_id,comment_id) VALUES (?,?,?,?,?)", Comment, date, id_user, id_post, id_comment)
}

func InsertUser(db *sql.DB, username string, email string, password string, confirm string) {
	if password != confirm {
		fmt.Println("Mots de passe différents.")
	}
	_, err := db.Exec("INSERT INTO Users (username, email, password) VALUES (?, ?, ?)", username, email, password)
	if err != nil {
		fmt.Println("Impossible d'enregistrer l'utilisateur")
	}
}

func GetUserByUsername(db *sql.DB, username string) (*s.User, error) {
	query := "SELECT id, username, email, password FROM users WHERE username = ?"
	row := db.QueryRow(query, username)

	user := &s.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("impossible de trouver l'user par son nom:%v", err)
	}
	return user, nil
}

func GetPosts(db *sql.DB) ([]s.Post, error) {
	_, err := db.Exec("SELECT * FROM Posts")
	if err != nil {
		fmt.Println("Erreur pour chopper les posts.")
	}
	rows, err := db.Query("SELECT id, title, description, user_id FROM Posts")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %v", err)
	}
	defer rows.Close()

	var posts []s.Post
	for rows.Next() {
		post := s.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.User_id)
		if err != nil {
			return nil, fmt.Errorf("echec lors du scan des colonnes: %v", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération %v", err)
	}
	return posts, nil
}

func Reboot(db *sql.DB) {
	_, err := db.Exec("UPDATE Users SET status = 'inactive'")
	if err != nil {
		fmt.Println("Impossible de reset les membres actifs.")
	}
}

func GetFriends(db *sql.DB, currentUser string) ([]s.User, error) {
	rows, err := db.Query("SELECT username, id, email, status FROM Users WHERE username != ?", currentUser)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %v", err)
	}
	defer rows.Close()

	var userList []s.User
	for rows.Next() {
		var user s.User
		if err := rows.Scan(&user.Username, &user.ID, &user.Email, &user.Status); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		userList = append(userList, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %v", err)
	}

	return userList, nil
}

func GetStatus(db *sql.DB, currentUser string) ([]s.Update, error) {
	rows, err := db.Query("SELECT username, status FROM Users where username != ?", currentUser)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve status: %v", err)
	}
	defer rows.Close()
	var statusList []s.Update
	for rows.Next() {
		var statue s.Update
		if err := rows.Scan(&statue.Name, statue.Status); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		statusList = append(statusList, statue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %v", err)
	}

	return statusList, nil
}
