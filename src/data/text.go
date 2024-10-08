package data

type RemoveTaskRequest struct {
	TaskId int `json:"task_id"`
}

type TaskDetailsResponse struct {
	TaskId       int    `json:"taskId"`
	Title        string `json:"title"`
	Contents     string `json:"contents"`
	Status       bool   `json:"status"`
	LastModified string `json:"lastModified"`
	CreatedOn    string `json:"createdOn"`
}

type NewTaskResponse struct {
	TaskId       int    `json:"taskId"`
	LastModified string `json:"lastModified"`
	CreatedOn    string `json:"createdOn"`
}

type NewTaskRequest struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Contents string `json:"contents"`
}

type UpdateTaskContentsRequest struct {
	TaskId   int    `json:"taskId"`
	Title    string `json:"title"`
	Contents string `json:"contents"`
	Status   bool   `json:"status"`
}
