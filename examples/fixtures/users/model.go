package users

import "github.com/akm/go-fixtures/examples/models"

type User struct {
	models.User
}

type Option = func(*User)

func Name(name string) Option   { return func(u *User) { u.Name = name } }
func Email(email string) Option { return func(u *User) { u.Email = email } }
