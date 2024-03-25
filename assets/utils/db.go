package utils

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	s "real/assets/struct"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Get All post from database
func SelectAllPosts(db *sql.DB) []*s.Post {
	var Date time.Time
	// Variables declaration
	var Posts []*s.Post
	// execute query to get all posts from database
	rows, err := db.Query("SELECT * FROM Posts ORDER BY date DESC")
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		var post s.Post
		err = rows.Scan(&post.ID, &post.Title, &post.ImgUrl, &post.Description, &Date, &post.User_id)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
		post.Date = Date.Format("02-01-2006 15:04:06")
		// Get Post Reaction
		SelectCountReaction(db, &post, nil)
		// fmt.Println("a", post.Name, post.Likes, post.Dislikes)
		// Get Pseudo Creator of post
		post.Pseudo = strings.Title(SelectPseudo(db, post.User_id))

		Posts = append(Posts, &post)

	}
	return Posts
}

// Get One post from database, thanks to Connexion Cookie --> id_user
func SelectPost(db *sql.DB, id int) *s.Post {
	// Variables declaration
	var post s.Post
	var Date time.Time
	// execute query to get all posts from database
	rows, err := db.Query("SELECT * FROM Posts WHERE id = ?", id)
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Title, &post.ImgUrl, &post.Description, &Date, &post.User_id)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}

		post.Date = Date.Format("02-01-2006 15:04:06")

		// Get Post Reaction
		SelectCountReaction(db, &post, nil)
		// fmt.Println("a", post.Name, post.Likes, post.Dislikes)
		// Get Pseudo Creator of post
		post.Pseudo = strings.Title(SelectPseudo(db, post.User_id))
	}
	return &post
}

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
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		err = rows.Scan(&postId)
		if err != nil {
			PrintError("Error Row:", err, 2)
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
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all categories rows get from database
	for rows.Next() {
		err = rows.Scan(&catId)
		if err != nil {
			PrintError("Error Row:", err, 2)
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
	var Date time.Time
	var Comments []*s.Comment

	rows, err := db.Query("SELECT c.id, u.id,u.Pseudo, c.message, c.date, c.comment_id FROM Comments c INNER JOIN Posts p ON c.post_id = p.id INNER JOIN Users u ON c.user_id = u.id WHERE post_id = ?;", post_id)
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		var comment s.Comment
		var reply sql.NullInt64 // Use sql.NullInt64 to handle NULL values in the database
		err = rows.Scan(&comment.ID, &comment.User_id, &comment.Pseudo, &comment.Message, &Date, &reply)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
		comment.Date = Date.Format("02-01-2006 15:04:06") // Problème dans la récupération, Date ne marche pas
		SelectCountReaction(db, nil, &comment)
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

// Inform how many like and dislike, a post contain
func SelectCountReaction(db *sql.DB, post *s.Post, comment *s.Comment) {
	var like int
	var dislike int
	var rows *sql.Rows
	var err error
	if post != nil {
		rows, err = db.Query("SELECT COUNT(CASE WHEN type = \"like\" THEN 1 END) AS NBlike,COUNT(CASE WHEN type = \"dislike\" THEN 1 END) AS NBdislike FROM Reactions INNER JOIN Posts ON id = post_id WHERE post_id = ?;", post.ID)
	} else {
		rows, err = db.Query("SELECT COUNT(CASE WHEN type = \"like\" THEN 1 END) AS NBlike,COUNT(CASE WHEN type = \"dislike\" THEN 1 END) AS NBdislike FROM Reactions r INNER JOIN Comments c ON c.id = r.comment_id WHERE c.id = ?;", comment.ID)
	}
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {
		err = rows.Scan(&like, &dislike)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
		// fmt.Println("b", post.Name, post.Likes, post.Dislikes)
		if post != nil {
			post.Update(post.Title, post.Description, like, dislike)
		} else {
			comment.Update(comment.Message, like, dislike)
		}
	}
}

// Find the pseudo of a User with the given id (Cookie)
func SelectPseudo(db *sql.DB, user_id int) string {
	// Variables declaration
	var pseudo string
	// execute query to get all posts from database
	rows, err := db.Query("SELECT pseudo FROM Users WHERE id = ?", user_id)
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&pseudo)
		if err != nil {
			PrintError("Error Row:", err, 2)
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
		PrintError("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&pseudo, &email)
		if err != nil {
			PrintError("Error Row:", err, 2)
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
		PrintError("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&id_pseudo)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
	}
	return id_pseudo
}

// Verify Password with bcrypt algorithm CompareHashAndPassword([]byte(Passwordindb), []byte(tested password))
func ComparePassword(db *sql.DB, pseudo string, password string) (bool, error) {
	var dbpassword string

	passwordDB := "SELECT password FROM Users WHERE pseudo= ?"
	err := db.QueryRow(passwordDB, pseudo).Scan(&dbpassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Aucun utilisateur trouvé avec ce nom d'utilisateur.
			fmt.Println("Aucun utilisateur trouvé avec ce nom d'utilisateur.")
			return false, nil
		}
		// Une erreur s'est produite lors de la requête à la base de données.
		fmt.Println("Une erreur s'est produite lors de la requête à la base de données")
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {
		// Le mot de passe ne correspond pas.
		fmt.Println("Le mot de passe ne correspond pas")
		return false, nil
	}

	// L'utilisateur existe et le mot de passe est correct.
	return true, nil
}

func InsertComment(db *sql.DB, Comment string, date string, id_user int, id_post int) {
	db.Exec("INSERT INTO Comments(message,date,user_id, post_id) VALUES (?,?,?,?)", Comment, date, id_user, id_post)
}

func InsertReply(db *sql.DB, Comment string, date string, id_user int, id_post int, id_comment int) {
	db.Exec("INSERT INTO Comments(message,date,user_id, post_id,comment_id) VALUES (?,?,?,?,?)", Comment, date, id_user, id_post, id_comment)
}

func SelectAllCategory(db *sql.DB) []string {
	var query string = "SELECT name FROM Categories"
	var name string
	var names []string
	rows, err := db.Query(query)
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&name)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
		names = append(names, name)
	}
	return names
}

// Get AllPosts from database with filter

func SelectAllPostFiltered(db *sql.DB, userID int, category string, created bool, react bool, commented bool, userPost string, byLike int, byDislike int, byComment int, byDate int) []*s.Post {
	// Variables declaration
	var Posts []*s.Post
	var Date time.Time
	var asWhere bool = false
	// SELECT
	var query string = "SELECT DISTINCT p.id, p.title, p.imgUrl, p.description, p.date, u.pseudo,"
	// count like
	query = fmt.Sprintf("%s (SELECT COUNT(CASE WHEN type = \"like\" THEN 1 END)", query)
	query = fmt.Sprintf("%s FROM Reactions INNER JOIN Posts ON id = post_id WHERE post_id = p.id) AS NBlike,", query)
	// count dislike
	query = fmt.Sprintf("%s (SELECT COUNT(CASE WHEN type = \"dislike\" THEN 1 END)", query)
	query = fmt.Sprintf("%s FROM Reactions INNER JOIN Posts ON id = post_id WHERE post_id = p.id) AS NBdislike,", query)
	// count comments
	query = fmt.Sprintf("%s (SELECT COUNT(*)", query)
	query = fmt.Sprintf("%s FROM Comments WHERE post_id = p.id) AS NBcomment", query)
	// end of sub request
	query = fmt.Sprintf("%s FROM Posts p", query)
	// JOIN
	query = fmt.Sprintf("%s INNER JOIN Users u ON p.user_id = u.id", query)
	if react {
		query = fmt.Sprintf("%s INNER JOIN Reactions r ON p.id = r.post_id", query)
	}
	if commented {
		query = fmt.Sprintf("%s INNER JOIN Comments co ON p.id = co.post_id", query)
	}
	if category != "" {
		query = fmt.Sprintf("%s LEFT JOIN Categorized c ON c.post_id = p.id", query)
		query = fmt.Sprintf("%s LEFT JOIN Categories ca ON ca.id = c.category_id", query)
	}
	// WHERE/ORDERBY
	if category != "" {
		if !asWhere {
			asWhere = true
			query = fmt.Sprintf("%s WHERE", query)
		}
		query = fmt.Sprintf("%s ca.name = '%s'", query, category)
	}
	if created && userID != 0 {
		if !asWhere {
			asWhere = true
			query = fmt.Sprintf("%s WHERE", query)
		} else {
			query = fmt.Sprintf("%s AND", query)
		}
		query = fmt.Sprintf("%s p.user_id = '%d'", query, userID)
	}
	if react && userID != 0 {
		if !asWhere {
			asWhere = true
			query = fmt.Sprintf("%s WHERE", query)
		} else {
			query = fmt.Sprintf("%s AND", query)
		}
		query = fmt.Sprintf("%s r.user_id = '%d'", query, userID)
	}
	if commented && userID != 0 {
		if !asWhere {
			asWhere = true
			query = fmt.Sprintf("%s WHERE", query)
		} else {
			query = fmt.Sprintf("%s AND", query)
		}
		query = fmt.Sprintf("%s co.user_id = '%d'", query, userID)
	}
	query = fmt.Sprintf("%s ORDER BY", query)
	if byLike != 0 {
		query = fmt.Sprintf("%s NBlike", query)
		if byLike == 1 {
			query = fmt.Sprintf("%s ASC,", query)
		} else {
			query = fmt.Sprintf("%s DESC,", query)
		}
	}
	if byDislike != 0 {
		query = fmt.Sprintf("%s NBdislike", query)
		if byDislike == 1 {
			query = fmt.Sprintf("%s ASC,", query)
		} else {
			query = fmt.Sprintf("%s DESC,", query)
		}
	}
	if byComment != 0 {
		query = fmt.Sprintf("%s NBcomment", query)
		if byComment == 1 {
			query = fmt.Sprintf("%s ASC,", query)
		} else {
			query = fmt.Sprintf("%s DESC,", query)
		}
	}
	if byDate != 0 {
		query = fmt.Sprintf("%s p.date", query)
		if byDate == 1 {
			query = fmt.Sprintf("%s ASC,", query)
		} else {
			query = fmt.Sprintf("%s DESC,", query)
		}
	}
	query = fmt.Sprintf("%s p.Title;", query)

	// execute query to get all posts from database
	rows, err := db.Query(query)
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	// Read all posts rows get from database
	for rows.Next() {

		var post s.Post
		err = rows.Scan(&post.ID, &post.Title, &post.ImgUrl, &post.Description, &Date, &post.Pseudo, &post.Likes, &post.Dislikes, &post.Comments)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
		post.Date = Date.Format("02-01-2006 15:04:06")
		Posts = append(Posts, &post)

	}
	return Posts
}

func InsertCategories(db *sql.DB, user_id int, categories string) (int, error) {
	var idpost int
	rows, err := db.Query("SELECT p.id FROM Posts p INNER JOIN Users u on u.id = p.user_id WHERE u.id = ? ORDER BY p.date DESC LIMIT 1;", user_id)
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	for rows.Next() {

		err = rows.Scan(&idpost)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
	}
	// fmt.Println("données: id du post", idpost, "Categorie", categories)
	_, err = db.Exec("INSERT INTO Categorized(post_id, category_id) VALUES (?,(SELECT id FROM Categories WHERE name = ?));", idpost, categories)
	if err != nil {
		// Ajouter une fenettre modale dans le html
		return idpost, err
	}
	return idpost, nil
}

func HadReacted(db *sql.DB, tReaction string, user_id int, post_id any, comment_id any) bool {
	var countR int
	var rows *sql.Rows
	var err error
	if post_id != nil {
		rows, err = db.Query("SELECT COUNT(*) FROM Reactions WHERE post_id = ? AND user_id = ? AND type = ?", post_id, user_id, tReaction)
	} else {
		rows, err = db.Query("SELECT COUNT(*) FROM Reactions WHERE comment_id = ? AND user_id = ? AND type = ?", comment_id, user_id, tReaction)
	}
	if err != nil {
		PrintError("Error Query:", err, 2)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&countR)
		if err != nil {
			PrintError("Error Row:", err, 2)
		}
		if countR == 0 {
			return false
		} else {
			return true
		}
	}
	return false
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
