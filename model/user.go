package model


type User struct {
	ID string `json:"id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	PwdSalt string `json:"pwd_salt"`
	Status int `json:"status"`
	UserDirectoryId string `json:"user_directory_id"`
	Path string `json:"path"`
	Mobile string `json:"mobile"`
	Email string `json:"email"`
	IsLocal int `json:"is_local"`
	Code string `json:"code" gorm:"-"` // 用来进行oauth2登录用到的
}