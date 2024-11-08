package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/akm/go-fixtures"
	"github.com/akm/go-fixtures/examples/fixtures/articles"
	"github.com/akm/go-fixtures/examples/fixtures/users"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ResultItem struct {
	AuthorName string
	Count      int
}

func logic(db *sql.DB) ([]*ResultItem, error) {
	// author の name ごとに articles の件数を集計する
	rows, err := db.Query(`
SELECT
  u.name,
  COUNT(a.id) AS article_count
FROM
  users u
  LEFT JOIN articles a ON u.id = a.author_id
GROUP BY
  u.name
ORDER BY
  u.name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*ResultItem{}

	for rows.Next() {
		var name string
		var count int
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		result = append(result, &ResultItem{AuthorName: name, Count: count})
	}
	return result, nil
}

func TestExample(t *testing.T) {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/go-fixtures-db?parseTime=true")
	if err != nil {
		t.Fatalf("unable to open database: %v", err)
	}
	defer db.Close()

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{SlowThreshold: time.Second, LogLevel: logger.Info},
	)

	fx := fixtures.NewDB(
		mysql.New(mysql.Config{Conn: db}),
		&gorm.Config{Logger: gormLogger},
	)(t)
	fx.DeleteFromTable(t, &articles.Article{})
	fx.DeleteFromTable(t, &users.User{})

	fxArticles := articles.NewFixtures()
	fx.Create(t,
		fxArticles.FiveRules(),
		fxArticles.GoProverbs(),
		fxArticles.ABriefIntroduction(),
	)
	fxUsers := fxArticles.Users
	if fxUsers.RobPike().ID == 0 {
		fx.Create(t, fxUsers.RobPike())
	}
	if fxUsers.KenThompson().ID == 0 {
		fx.Create(t, fxUsers.KenThompson())
	}
	if fxUsers.RobertGriesemer().ID == 0 {
		fx.Create(t, fxUsers.RobertGriesemer())
	}

	result, err := logic(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("unexpected result length: %v", result)
	}
	for i, r := range result {
		t.Logf("%d: AuthorName: %s, Count: %d", i, r.AuthorName, r.Count)
	}

	if result[0].AuthorName != "Ken Thompson" {
		t.Fatalf("unexpected result[0].AuthorName: %+v", *result[0])
	}
	if result[0].Count != 1 {
		t.Fatalf("unexpected result[0].Count: %+v", *result[0])
	}

	if result[1].AuthorName != "Rob Pike" {
		t.Fatalf("unexpected result[1].AuthorName: %+v", *result[1])
	}
	if result[1].Count != 2 {
		t.Fatalf("unexpected result[1].Count: %+v", *result[1])
	}

	if result[2].AuthorName != "Robert Griesemer" {
		t.Fatalf("unexpected result[2].AuthorName: %+v", *result[2])
	}
	if result[2].Count != 0 {
		t.Fatalf("unexpected result[2].Count: %+v", *result[2])
	}
}
