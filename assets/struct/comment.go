package real

import "fmt"

type Comment struct {
	ID      int
	Message string
	Pseudo      string
	User_id     int
	SubComments []*Comment // store sub comments
}

// Update Information to a post
func (c *Comment) Update(
	message string,
	likes int,
	dislikes int,
) {
	c.Message = message
}

func GetComment(id int, comments []*Comment) *Comment {
	for _, c := range comments {
		fmt.Println( id ,"==",c.ID)
		if c.ID == id {
			fmt.Println( c.ID)
			return c
		}
	}
	return nil
}