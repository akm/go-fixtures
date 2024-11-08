package articles

import (
	"github.com/akm/go-fixtures/examples/fixtures/users"
	"github.com/akm/go-fixtures/examples/models"
)

type Article struct {
	models.Article

	Author *users.User
}

type Option = func(*Article)

func Title(v string) Option { return func(a *Article) { a.Title = v } }
func Body(v string) Option  { return func(a *Article) { a.Body = v } }
func AuthorID(v int) Option { return func(a *Article) { a.AuthorID = v } }
func Author(v *users.User) Option {
	return func(a *Article) {
		a.Author = v
		if v != nil {
			a.AuthorID = v.ID
		} else {
			a.AuthorID = 0
		}
	}
}
