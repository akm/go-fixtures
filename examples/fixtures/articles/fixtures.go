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
		Url("https://users.ece.utexas.edu/~adnan/pike.html"),
	)
}

func (f *Fixtures) NewGoProverbs(opts ...Option) *Article {
	return fixtures.NewWithDefaults(opts,
		Author(f.Users.RobPike()),
		Title("Go Proverbs"),
		Url("https://go-proverbs.github.io/"),
	)
}

func (f *Fixtures) NewABriefIntroduction(opts ...Option) *Article {
	return fixtures.NewWithDefaults(opts,
		Author(f.Users.KenThompson()),
		Title("A Brief Introduction"),
		Url("https://www.linfo.org/thompson.html"),
	)
}

func (f *Fixtures) FiveRules(opts ...Option) *Article  { return f.Get("FiveRules", opts...) }
func (f *Fixtures) GoProverbs(opts ...Option) *Article { return f.Get("GoProverbs", opts...) }
func (f *Fixtures) ABriefIntroduction(opts ...Option) *Article {
	return f.Get("ABriefIntroduction", opts...)
}
