package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/katenicoletti/fam-api"
	"github.com/katenicoletti/fam-api/pg"
)

// Main ...
type Main struct {
	db     *sql.DB
	client *http.Client
}

func main() {
	m := Main{
		db: pg.DB,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/launcher/{id}", m.handleGetLauncher).Methods("GET")
	router.HandleFunc("/launcher/{id}/transactions", m.handleGetTransactions).Methods("GET")
	router.HandleFunc("/payoff", m.handlePayoff).Methods("POST")

	log.Fatal(http.ListenAndServe(":6000", router))
}

func (m *Main) handleGetLauncher(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var l launcher.Launcher
	if err := m.db.QueryRow(`SELECT
      id
    , first_name
    , last_name
    , credit_limit
    , balance
    , due_date
    , minimum_payment
  FROM launchers
  WHERE id = $1`,
		params["id"],
	).Scan(
		&l.ID,
		&l.FirstName,
		&l.LastName,
		&l.Balance,
		&l.CreditLimit,
		&l.DueDate,
		&l.MinimumPayment,
	); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(l)
}

func (m *Main) handleGetTransactions(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	rows, err := m.db.Query(`SELECT
      id
    , type
    , merchant
    , amount
    , purchase_date
		FROM transactions
		WHERE launchers_id = $1`,
		params["id"],
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var transactions []launcher.Transaction
	for rows.Next() {
		var t launcher.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.Type,
			&t.Merchant,
			&t.Amount,
			&t.Date,
		); err != nil {
			panic(err)
		}
		transactions = append(transactions, t)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(transactions)
}

func (m *Main) handlePayoff(w http.ResponseWriter, req *http.Request) {
	var p launcher.Payoff
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		panic(err)
	}

	if err := p.Calculate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(p)
}
