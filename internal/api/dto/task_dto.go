package dto

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,min=3"`
}

type UpdateTaskRequest struct {
	Title string `json:"title" validate:"required,min=3"`
	Done  bool   `json:"done"`
}
