package users

import "github.com/akm/go-fixtures"

type Fixtures struct {
	*fixtures.Fixtures[User]
}

var (
	_ (fixtures.Factory[User]) = (*Fixtures)(nil)
	_ (fixtures.Getter[User])  = (*Fixtures)(nil)
)

func NewFixtures() *Fixtures {
	r := &Fixtures{}
	r.Fixtures = fixtures.NewFixtures[User](r)
	return r
}

func (f *Fixtures) NewRobPike(opts ...Option) *User {
	return fixtures.NewWithDefaults(opts,
		Name("Rob Pike"),
		Email("rob-pike@example.com"),
	)
}

func (f *Fixtures) NewKenThompson(opts ...Option) *User {
	return fixtures.NewWithDefaults(opts,
		Name("Ken Thompson"),
		Email("ken-thompson@example.com"),
	)
}

func (f *Fixtures) NewRobertGriesemer(opts ...Option) *User {
	return fixtures.NewWithDefaults(opts,
		Name("Robert Griesemer"),
		Email("robert-griesemer@example.com"),
	)
}

func (f *Fixtures) RobPike(opts ...Option) *User         { return f.Get("RobPike", opts...) }
func (f *Fixtures) KenThompson(opts ...Option) *User     { return f.Get("KenThompson", opts...) }
func (f *Fixtures) RobertGriesemer(opts ...Option) *User { return f.Get("RobertGriesemer", opts...) }
