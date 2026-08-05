package request

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   *string  `json:"title" validate:"omitempty,min=1,max=500"`
	Content *string  `json:"content" validate:"omitempty,min=1,max=50000"`
	Tags    []string `json:"tags" validate:"max=20,dive,min=1,max=50"`
}
