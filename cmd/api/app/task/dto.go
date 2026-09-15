package api

type CreateTaskRequest struct {
	Title  string `json:"title"`
	UserID string `json:"userId"`
}

type UpdateTaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
