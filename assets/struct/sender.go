package real

type HeaderHandler struct {
	IsLogged bool
	Pseudo   string
	ID       int
	Email    string
}

type CommentHandler struct {
	Comment  []*Comment
	IsLogged bool
}
type FilterHandler struct {
	Categories []string
	IsLogged   bool
}
type User struct {
	ID     int    `json:"id"`
	Pseudo string `json:"pseudo"`
}
