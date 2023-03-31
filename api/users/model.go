package users

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

func (u *User) CleanPassword() {
	u.Password = ""
}

type CreateDtoUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInDtoUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
