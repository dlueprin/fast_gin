package model

type UserModel struct {
	Model           //继承
	Username string `gorm:"size:16 json:username"`
	Nickname string `gorm:"size:32 json:nickname"`
	Password string `gorm:"size:64 json:password"`
	RoleID   int8   `json:"roleID"` //简单模式 1 管理员 2 普通用户，int8已经是go最小的整数类型了，-128到127
}
