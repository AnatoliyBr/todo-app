package entity

type Task struct {
	TaskID    int     `json:"task_id"`
	TaskTitle string  `json:"task_title"`
	Details   string  `json:"details"`
	Deadline  TimeISO `json:"deadline"`
	Done      bool    `json:"done"`
	ListID    int     `json:"list_id"`
}
