package utils

import (
	"database/sql"
	"net/http"
	s "real/assets/struct"
)

func FilterTemplate(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Categories := SelectAllCategory(db)
	// filterHandler := s.FilterHandler{
	// 	Categories: Categories,
	// 	IsLogged:   SessionExists(r, "Connexion"),
	// }
	// tmpl := template.Must(template.ParseFiles("asset/html/filter.tmpl"))
	// err := tmpl.Execute(w, filterHandler)
	// if err != nil {
	// 	PrintError("FilterTemplate error:", err, -1)
	// }
}
func GetComment(comments []*s.Comment, id int) *s.Comment {
	for _, comment := range comments {
		if comment.ID == id {
			return comment
		}
	}
	return nil
}
