package models

type User struct {
	UserId       int    `json:"id"`
	UserName     string `json:"name"`
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
	Superuser    bool   `json:"superuser"`
}

// func NewUser(name, email, password string, superuser bool) (*User, error) {
// 	if name == "" || email == "" || password == "" {
// 		return &User{}, errors.New("нельзя передавать пустые значения")
// 	}

// 	return &User{
// 		UserName:     name,
// 		UserEmail:    email,
// 		UserPassword: password,
// 		Superuser:    superuser,
// 	}, nil
// }

type UserDTO struct {
	UserId    int    `json:"id"`
	UserName  string `json:"name"`
	UserEmail string `json:"email"`
}
