package user

type User struct {
	Id        int     `json:"-"`
	Email     string  `json:"email"`
	Password  string  `json:"-"`
	Directory *string `json:"dir"`
}
