package model

type FriendUser struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	FollowCount     uint   `json:"follow_count"`
	FollowerCount   uint   `json:"follower_count"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	TotalFavorited  uint   `json:"total_favorited"`
	FavoriteCount   uint   `json:"favorite_count"`
	IsFollow        bool   `json:"is_follow"`
	Message         string `json:"message"` //聊天信息
	MsgType         int8   `json:"msgType"` //message信息的类型，0=>请求用户接受信息，1=>当前请求用户发送的信息
}
