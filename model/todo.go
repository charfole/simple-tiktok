package model

// Todo Model 默认映射到当前数据库中的todos表
type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Status bool   `json:"status"`
}
