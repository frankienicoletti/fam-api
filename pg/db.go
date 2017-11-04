package pg

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// DB ...
var DB *sql.DB

func init() {
	// Get env vars
	port := "5432"     // os.Getenv("PGPORT")
	user := "postgres" // os.Getenv("PGUSER")
	// password := os.Getenv("PGPASSWORD")
	name := "postgres" // TODO rename
	host := os.Getenv("PGHOST")
	if host == "" {
		host = "localhost"
	}

	// Connect to db.
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, name))
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DB = db
	fmt.Println("successfully connected to database")

	// TODO can add create scripts here
}
