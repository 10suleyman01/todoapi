package todo

type Todo struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	UserId string `json:"user_id"`
}

type CreateTodoDto struct {
	Title string `json:"title"`
}

type DeleteTodoDto struct {
	TodoId string `json:"todo_id"`
}
