package todo_dto

type CreateTodoRequest struct {
	TodoName        string `json:"todo_name" binding:"required,min=3,max=20"`
	TodoDescription string `json:"todo_description" binding:"required,min=5,max=255"`
}

type GetTodoRequest struct {
	TodoID string `json:"todo_id" binding:"required,uuid"`
}
