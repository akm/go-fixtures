# github.com/akm/go-fixtures 

- module name: `github.com/akm/go-fixtures`
- package name: `fixtures`

## Overview

The `fixtures` package simplifies database operations in tests using models based on the Functional Option Pattern. It integrates with [GORM](https://gorm.io/) for object manipulation but remains agnostic to the specific models you use. Test fixtures can be defined directly in Go code.

`fixtures` package supports database operations in test with Functional Option Pattern based models.
`fixtures` uses [GORM](https://gorm.io/) to manipulate objects but it doesn't matter your model.
You can define your test fixtures in Go code.

### ORM support

#### GORM

`fixtures` works seamlessly with models designed for [GORM](https://gorm.io/), making it straightforward to create and manage fixtures.


#### sqlc

Even if you use sqlc to generate models, fixtures is compatible. The structure of sqlc-generated models is similar to GORM models, making the integration easy. For example:

- [Original sqlc model](https://github.com/akm/svelte-connect-todo/blob/9e77fa3d0ae1777ab23c2ba4753ff4bd541f44ed/backends/biz/models/models.go#L63-L69)
- [Fixture model](https://github.com/akm/svelte-connect-todo/blob/9e77fa3d0ae1777ab23c2ba4753ff4bd541f44ed/backends/biz/fixtures/tasks/model.go)
- [Defined fixtures](https://github.com/akm/svelte-connect-todo/blob/9e77fa3d0ae1777ab23c2ba4753ff4bd541f44ed/backends/biz/fixtures/tasks/fixtures.go)

#### Other ORMs

Even if you're using another ORM or none at all, fixtures allows you to define plain fixture models. Feel free to reach out if you encounter any issues.

## Install

```sh
go get github.com/akm/go-fixtures
```

## How to Create Models and Fixtures

1. Define fixture models.
1. Add functional `Option` types for the models.
1. Create a `Fixtures` struct for each model.
1. Use the defined fixtures in your tests.

### 1. Define Fixture Models

If you already have a model, you can embed it. Otherwise, define a new model or skip this step if it's unnecessary.

#### Original Model

```golang
type Article struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	AuthorID  int
	Title     string
	Url       string
}
```

#### Fixture Model

```golang
type Article struct {
	models.Article

	Author *users.User
}
```

For examples, see:

- [examples/models/user.go](./examples/models/user.go)
- [examples/fixtures/users/model.go](./examples/fixtures/users/model.go)
- [examples/models/article.go](./examples/models/article.go)
- [examples/fixtures/articles/model.go](examples/fixtures/articles/model.go)

### 2. Add Functional `Option` Types

Define Option types and corresponding functions using the Functional Option Pattern. Refer to [Rob Pike's article](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html) and [Dave Cheney's article](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).

```golang
type Option = func(*Article)

func Title(v string) Option { return func(a *Article) { a.Title = v } }
func Url(v string) Option   { return func(a *Article) { a.Url = v } }
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
```

You can instantiate models using `fixtures.New[T any]` or `fixtures.NewWithDefaults[T any]`.

### 3. Make Fixtures for each fixture model

#### Define the `Fixtures` Struct

```golang
type Fixtures struct {
	*fixtures.Fixtures[Article]
	Users *users.Fixtures
}
```

Embed `*fixtures.Fixtures` to inherit its methods `New` and `Get`.
Include related fixture types as fields if fixture model refers other models.

#### Add a Constructor Function

```golang
func NewFixtures() *Fixtures {
	r := &Fixtures{
		Users: users.NewFixtures(),
	}
	r.Fixtures = fixtures.NewFixtures[Article](r)
	return r
}
```

#### Add `NewXXX` Methods for each Fixture.

Define methods to create specific fixture instances:

```
func (f *Fixtures) NewFiveRules(opts ...Option) *Article {
	return fixtures.NewWithDefaults(opts,
		Author(f.Users.RobPike()),
		Title("Rob Pike's 5 Rules of Programming"),
		Url("https://users.ece.utexas.edu/~adnan/pike.html"),
	)
}
```

These methods must:

- Be variadic to accept zero or more `Option` s.
- Return a pointer to the model.


#### Add Fixture Retrieval Methods

```
func (f *Fixtures) FiveRules(opts ...Option) *Article  { return f.Get("FiveRules", opts...) }
```

This is a fixed form. You can change the method name, but the first argument of `f.Get` must match the method name.


For examples, see:
- [examples/fixtures/users/fixtures.go](./examples/fixtures/users/fixtures.go)
- [examples/fixtures/articles/fixtures.go](./examples/fixtures/articles/fixtures.go)

### 4. Use Fixtures

#### Setup `fixtures.DB`

`fixtures` uses [GORM](https://gorm.io).

To set up `fixtures.DB`, you need to provide a `gorm.Dialector` and optionally one or more `gorm.Option`s. If you use a specific database (e.g., MySQL), make sure to install the corresponding GORM driver. For instance, for MySQL, run:

```sh
go get gorm.io/driver/mysql
```

Here’s an example of setting up `fixtures.DB`:

```golang
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{SlowThreshold: time.Second, LogLevel: logger.Info},
	)

	fx := fixtures.NewDB(
		mysql.New(mysql.Config{Conn: db}),
		&gorm.Config{Logger: gormLogger},
	)(t)
```

#### Deleting Records from a Table

To delete all records from specific tables, use `fx.DeleteFromTable`. Pass instances of your model with zero-value fields as arguments:

```golang
fx.DeleteFromTable(t, &articles.Article{}, &users.User{})
```


#### Inserting Fixture Data

To insert fixture data, create an instance of your `Fixtures` struct and call `fx.Create`, passing the fixture instances you want to add to the database.


```golang
	fxArticles := articles.NewFixtures()
	fx.Create(t,
		fxArticles.FiveRules(),
		fxArticles.GoProverbs(),
		fxArticles.ABriefIntroduction(),
	)
```

#### Additional Details

For more information on setting up and using GORM, see [Connecting to a Database](https://gorm.io/docs/connecting_to_the_database.html). You can also refer to [examples/tests/example_test.go](./examples/tests/example_test.go) for a complete implementation.

## Test

You can view the test results in GitHub Actions.[GitHub Actions](https://github.com/akm/go-fixtures/actions).

## License

MIT License
