package response

type PostResponse struct {
	ID              int64             `json:"id"`
	Content         string            `json:"content"`
	Title           string            `json:"title"`
	UserId          int64             `json:"user_id"`
	Tags            []string          `json:"tags"`
	Version         int               `json:"version"`
	CommentResponse []CommentResponse `json:"comments"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}
