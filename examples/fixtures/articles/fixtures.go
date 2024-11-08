package articles

import (
	"github.com/akm/go-fixtures"
	"github.com/akm/go-fixtures/examples/fixtures/users"
)

type Fixtures struct {
	*fixtures.Fixtures[Article]
	Users *users.Fixtures
}

var (
	_ (fixtures.Factory[Article]) = (*Fixtures)(nil)
	_ (fixtures.Getter[Article])  = (*Fixtures)(nil)
)

func NewFixtures() *Fixtures {
	r := &Fixtures{
		Users: users.NewFixtures(),
	}
	r.Fixtures = fixtures.NewFixtures[Article](r)
	return r
}

func (f *Fixtures) NewFiveRules(opts ...Option) *Article {
	return fixtures.NewWithDefaults(opts,
		Author(f.Users.RobPike()),
		Title("Rob Pike's 5 Rules of Programming"),
		Body("Rule 1. You can't tell where a program is going to spend its time..."),
	)
}

func (f *Fixtures) NewGoProverbs(opts ...Option) *Article {
	return fixtures.NewWithDefaults(opts,
		Author(f.Users.RobPike()),
		Title("Go Proverbs"),
		Body("Simple, Poetic, Pithy..."),
	)
}

func (f *Fixtures) NewABriefIntroduction(opts ...Option) *Article {
	return fixtures.NewWithDefaults(opts,
		Author(f.Users.KenThompson()),
		Title("A Brief Introduction"),
		Body("Kenneth Lane Thompson was the principal inventor of UNIX..."),
	)
}

func (f *Fixtures) FiveRules(opts ...Option) *Article  { return f.Get("FiveRules", opts...) }
func (f *Fixtures) GoProverbs(opts ...Option) *Article { return f.Get("GoProverbs", opts...) }
func (f *Fixtures) ABriefIntroduction(opts ...Option) *Article {
	return f.Get("ABriefIntroduction", opts...)
}
