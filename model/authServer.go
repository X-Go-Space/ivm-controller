package model

type SuccessCondition struct {
	ResponseFiled string `json:"response_filed"`
	ResponseCondition string `json:"response_condition"`
	ResponseResult string `json:"response_result"`
}

type AuthConfig struct {
	RequestType string `json:"request_type"`
	BaseUrl string `json:"base_url"`
	Headers [][2]string `json:"headers"`
	Params [][2]string `json:"params"`
	Body    string    `json:"body"`
	SuccessCondition SuccessCondition `json:"success_condition"`
}

type AuthServer struct {
	Id string `json:"id"`
	Name string `json:"name"`
	UserDirectoryId string `json:"user_directory_id"`
	AuthConfig []AuthConfig  `json:"auth_config" gorm:"-"`
	AuthConfigJson string `json:"-"`
}