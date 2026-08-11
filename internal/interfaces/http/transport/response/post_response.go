package response

import "time"

const VietNamTimeFormat = "02/01/2006 15:04:05"

var VietnamLocation = func() *time.Location {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return time.FixedZone("ICT", 7*60*60)
	}
	return location
}()

type PostResponse struct {
	ID              int64             `json:"id"`
	Content         string            `json:"content"`
	Title           string            `json:"title"`
	UserId          int64             `json:"user_id"`
	Username        string            `json:"username"`
	Tags            []string          `json:"tags"`
	Version         int               `json:"version"`
	CommentResponse []CommentResponse `json:"comments"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}

type PostWithMetaData struct {
	PostResponse
	CommentsCount int `json:"comments_count"`
}
