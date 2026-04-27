package dto

type CreateTaskRequest struct {
	Title  string `json:"title" validate:"required,min=3"`
	UserID string `json:"userId" validate:"required"`
}

type UpdateTaskRequest struct {
	Title string `json:"title" validate:"required,min=3"`
	Done  bool   `json:"done"`
}
