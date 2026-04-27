package domain


type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Done   bool   `json:"done,omitempty"`
	UserID int    `json:"UserID"`
}
