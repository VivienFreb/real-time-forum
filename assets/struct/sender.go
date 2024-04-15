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
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type LoginResponse struct {
	Name    string `json:"Name"`
	Success bool   `json:"Success"`
	Message string `json:"Message"`
}

type ArrayConcept struct {
	Name string   `json:"Name"`
	Ray  []string `json:"Ray"`
}

type AllMyFellas struct {
	Name  string `json:"Name"`
	Users []User `json:"Users"`
}

type Update struct {
	Name   string `json:"Name"`
	Status string `json:"Status"`
}

type NewStatus struct {
	Name   string   `json:"Name"`
	Checks []Update `json:"Checks"`
}

type MessageInner struct {
	Speaker  string `json:"Speaker"`
	Listener string `json:"Listener"`
	Content  string `json:"Content"`
}

type MessageOuter struct {
	Name  string      `json:"Name"`
	Chats []MessageInner `json:"Chats"`
}