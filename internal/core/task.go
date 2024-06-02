package core

type Task struct {
	Id          int
	Header      string `json:"header" db:"header" validate:"required"`
	Description string `json:"description" db:"descr" validate:"required"`
	TaskDate    string `json:"task_date" db:"task_date" validate:"required"`
	Done        *bool  `json:"done" db:"done"`
}
