package models

type Project struct {
	ID int64 `json:"-"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type Task struct {
	ID int64 `json:"-"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreationDate string `json:"creationDate"`
}

type Subtask struct {
	ID int64 `json:"-"`
	Title string `json:"title"`
	Description string `json:"description"`
	CreationDate string `json:"creationDate"`
}
