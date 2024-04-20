package models

type User struct {
	UserId       int    `json:"id"`
	UserName     string `json:"name"`
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
	Superuser    bool   `json:"superuser"`
}

type UserDTO struct {
	UserId    int    `json:"id"`
	UserName  string `json:"name"`
	UserEmail string `json:"email"`
	Superuser bool   `json:"superuser"`
}
