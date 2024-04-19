package utils

import (
	"database/sql"
	"fmt"

	s "real/assets/struct"

	_ "github.com/mattn/go-sqlite3"
)

func InsertUser(db *sql.DB, username string, email string, password string, confirm string) {
	if password != confirm {
		fmt.Println("Mots de passe différents.")
	}
	_, err := db.Exec("INSERT INTO Users (username, email, password, status) VALUES (?, ?, ?,?)", username, email, password, "inactive")
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
	rows, err := db.Query("SELECT id, title, description, user_name FROM Posts")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %v", err)
	}
	defer rows.Close()
	var posts []s.Post
	for rows.Next() {
		post := s.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.User_name)
		if err != nil {
			return nil, fmt.Errorf("echec lors du scan des colonnes: %v", err)
		}
		comments, err := GetComments(db, post.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch comments: %v", err)
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de l'itération %v", err)
	}
	return posts, nil
}

func GetComments(db *sql.DB, postTitle string) ([]s.Comment, error) {
	rows, err := db.Query("SELECT id, user_name, content FROM Comments WHERE parent = ?", postTitle)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve comments: %v", err)
	}
	defer rows.Close()

	var comments []s.Comment
	for rows.Next() {
		comment := s.Comment{}
		err := rows.Scan(&comment.ID, &comment.Username, &comment.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment rows: %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during comment iteration: %v", err)
	}
	return comments, nil
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
		if err := rows.Scan(&statue.Name, &statue.Status); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		statusList = append(statusList, statue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %v", err)
	}
	ForcedActive(db, currentUser)
	return statusList, nil
}

func ForcedActive(db *sql.DB, username string) {
	var status string
	err := db.QueryRow("SELECT status FROM Users WHERE username = ?", username).Scan(&status)
	if err != nil {
		fmt.Println("Error querying user status:", err)
		return
	}
	if status != "active" {
		_, updateErr := db.Exec("UPDATE Users SET status = 'active' WHERE username = ?", username)
		if updateErr != nil {
			fmt.Println("Error updating user status:", updateErr)
			return
		}
	}
}

func Activation(db *sql.DB, currentUser string) error {
	_, err := db.Exec("UPDATE Users SET status = 'active' WHERE username = ?", currentUser)
	if err != nil {
		return fmt.Errorf("failed to activate user: %v", err)
	}
	fmt.Println(currentUser, "is now active") // Log successful activation
	return nil
}

func Deactivation(db *sql.DB) error {
	_, err := db.Exec("UPDATE Users SET status = 'inactive'")
	if err != nil {
		return fmt.Errorf("failed to deactivate users: %v", err)
	}
	return nil
}

func GetDiscussion(db *sql.DB, currentUser string, otherUser string) ([]s.MessageInner, error) {
	rows, err := db.Query("SELECT Speaker, Listener, Content FROM Disscussions WHERE (Speaker = ? AND Listener = ?) OR (Speaker = ? AND Listener = ?)", currentUser, otherUser, otherUser, currentUser)
	if err != nil {
		return nil, fmt.Errorf("impossible de récupérer les discussions: %v", err)
	}
	defer rows.Close()

	var discussionList []s.MessageInner
	for rows.Next() {
		var message s.MessageInner
		if err := rows.Scan(&message.Speaker, &message.Listener, &message.Content); err != nil {
			return nil, fmt.Errorf("impossible de scan une discussion: %v", err)
		}
		discussionList = append(discussionList, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur pendant l'itération des convs: %v", err)
	}
	return discussionList, nil
}

func NewMessage(db *sql.DB, me string, them string, message string) {
	_, err := db.Exec("INSERT INTO Disscussions (Speaker, Listener, Content) VALUES (?, ?, ?)", me, them, message)
	if err != nil {
		fmt.Println("Impossible d'enregistrer le message.")
	}
}

func Delog(db *sql.DB, username string) {
	_, err := db.Exec("UPDATE Users SET status = 'inactive' WHERE username = ?", username)
	if err != nil {
		fmt.Printf("Impossible de delog %s.\n", username)
	}
}

func CreatePost(db *sql.DB, username string, subject string, content string) {
	_, err := db.Exec("INSERT INTO Posts (title, description, user_name) VALUES (?,?,?)", subject, content, username)
	if err != nil {
		fmt.Println("Impossible d'ajouter le post à la database.")
	}
}