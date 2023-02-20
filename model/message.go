package model

type Message struct {
	ID         int64  `json:"id"`
	ToUserId   uint   `json:"to_user_id"`
	UserId     uint   `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time" gorm:"column:createTime"` // 创建时间
}
