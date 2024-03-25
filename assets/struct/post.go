package real

type Post struct {
	ID          int
	Title       string
	ImgUrl      string
	Description string
	Date        string
	Likes       int
	Dislikes    int
	Comments    int
	User_id     int
	Pseudo      string
}

type MyPost struct{
	ID int
	Title string
	Description string
	Date string
	Comments int
}

// Builder
// Create a new Post
func NewPost(
	id int,
	title string,
	imgUrl string,
	description string,
	// date time.Time,
	date string,
	likes int,
	dislikes int,
	NBcomments int,
	user_id int,
	pseudo int,
) Post {
	return Post{
		ID:          id,
		Title:       title,
		ImgUrl:      imgUrl,
		Description: description,
		Date:        date,
		Likes:       likes,
		Dislikes:    dislikes,
		User_id:     user_id,
	}
}

// Update Information to a post
func (p *Post) Update(
	title string,
	description string,
	likes int,
	dislikes int,
) {
	p.Title = title
	p.Description = description
	p.Likes = likes
	p.Dislikes = dislikes
}
