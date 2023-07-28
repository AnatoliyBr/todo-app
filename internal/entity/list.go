package entity

type List struct {
	ListID    int    `json:"list_id"`
	ListTitle string `json:"list_title"`
	UserID    int    `json:"user_id"`
}
