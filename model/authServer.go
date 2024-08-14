package model

type Header struct {
	From string `json:"from"`
	TO string  `json:"to"`
}

type AuthConfig struct {
	Headers [][2]string `json:"headers"`
	Params [][2]string `json:"params"`
	Body    string    `json:"body"`
}

type AuthServer struct {
	Id string `json:"id"`
	Name string `json:"name"`
	UserDirectoryId string `json:"user_directory_id"`
	AuthConfig AuthConfig  `json:"auth_config" gorm:"-"`
	AuthConfigJson string `json:"-"`
}