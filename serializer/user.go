package serializer

import "chatDemo/model"

type User struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	NickName string `json:"nickname"`
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
	CreateAt int64  `json:"create_at"`
}

func BindUser(user model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Status:   user.Status,
		Avatar:   user.AvatarURL(),
		CreateAt: user.CreatedAt.Unix(),
	}
}

func BindUsers(items []model.User) (users []User) {
	for _, item := range items {
		user := BindUser(item)
		users = append(users, user)
	}
	return users
}
