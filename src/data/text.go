package data

type TaskId struct {
	TaskId int `json:"task_id"`
}

type NewTask struct {
	Username string `json:"username"`
	Contents string `json:"contents"`
}

type UpdateTask struct {
	TaskId   int    `json:"task_id"`
	Contents string `json:"contents"`
}
