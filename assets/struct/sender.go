package utils

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
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct{
	Name string `json:"Name"`
	Success bool `json:"Success"`
	Message string `json:"Message"`
}