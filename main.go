package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/mgutz/dat.v1"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

// global database (pooling provided by SQL driver)
var DB *runner.DB

func init() {
	// create a normal database connection through database/sql
	db, err := sql.Open("postgres", "dbname=dat_test user=dat password=!test host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(db)

	// set to reasonable values for production
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 10 * time.Millisecond

	DB = runner.NewDB(db, "postgres")
}

type Post struct {
	ID        int64        `db:"id"`
	Title     string       `db:"title"`
	Body      string       `db:"body"`
	UserID    int64        `db:"user_id"`
	State     string       `db:"state"`
	UpdatedAt dat.NullTime `db:"updated_at"`
	CreatedAt dat.NullTime `db:"created_at"`
}

func main() {
	var post Post
	err := DB.
		Select("id, title").
		From("posts").
		Where("id = $1", 13).
		QueryStruct(&post)
	fmt.Println("Title", post.Title)
}
