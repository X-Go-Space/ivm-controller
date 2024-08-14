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
}