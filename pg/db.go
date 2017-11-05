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

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS launchers (
  id serial primary key,
  customer_id bigint NOT NULL UNIQUE,
  account_id bigint NOT NULL,
  first_name text NOT NULL,
  last_name text NOT NULL,
  interest_rate float NOT NULL,
  credit_limit bigint NOT NULL,
  balance float default 0,
  due_date timestamp with time zone NOT NULL,
  minimum_payment float default 0,
  reward_balance float default 0,
  created timestamp with time zone default current_timestamp,
  modified timestamp with time zone default NULL
);`); err != nil {
		panic(err)
	} else if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS transactions (
    id bigint primary key,
    launchers_id int NOT NULL,
    type text NOT NULL,
    merchant text NOT NULL,
    amount float NOT NULL,
    purchase_date timestamp with time zone,
    CONSTRAINT launchers_fk
      FOREIGN KEY(launchers_id) REFERENCES launchers
      ON DELETE CASCADE
);`); err != nil {
		panic(err)
	}

	fmt.Println("set up database")
}
