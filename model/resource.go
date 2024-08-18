package model

type Resource struct {
	Id string `json:"id"`
	Name string `json:"name"`
	RedirectUrl string `json:"redirect_url"`
	Md5 string `json:"md_5"`
	UserList []User `json:"user_list" gorm:"many2many:resource_to_user;"`
}
