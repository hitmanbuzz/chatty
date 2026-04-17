package main

type new_user struct {
	username string
}

func create_user(username string) *new_user {
	return &new_user{
		username: username,
	}
}

type User struct {
	UserID   *uint
	GroupID  *uint
	Username string
	IsOnline bool
}

func NewUser(username string) *User {
	return &User{
		UserID:   nil,
		GroupID:  nil,
		Username: username,
		IsOnline: false,
	}
}

func (u *User) CreateUserData() {

}
