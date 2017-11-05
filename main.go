package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/katenicoletti/fam-api/launcher"
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

	m.seedDatabase()

	router := mux.NewRouter()
	router.HandleFunc("/launchers/{id}", m.handleGetLauncher).Methods("GET")
	router.HandleFunc("/launchers/{id}/transactions", m.handleGetTransactions).Methods("GET")
	router.HandleFunc("/launchers/{id}/payoff", m.handlePayoff).Methods("POST")

	log.Fatal(http.ListenAndServe(":6000", router))
}

func (m *Main) handleGetLauncher(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var l launcher.Launcher
	if err := m.db.QueryRow(`SELECT
      id
    , first_name
    , last_name
		, interest_rate
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
		&l.InterestRate,
		&l.CreditLimit,
		&l.Balance,
		&l.DueDate,
		&l.MinimumPayment,
	); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (m *Main) handlePayoff(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var p launcher.Payoff
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		panic(err)
	}

	var balance, interestRate float64
	if err := m.db.QueryRow(`SELECT
			  balance
			, interest_rate
		FROM launchers
		WHERE id = $1`,
		params["id"],
	).Scan(
		&balance,
		&interestRate,
	); err != nil {
		log.Fatal(err)
	}

	if p.Balance == 0 {
		p.Balance = balance
	}
	p.InterestRate = interestRate

	if err := p.Calculate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Collects and stores data from Capital One hackathon API.
// docker run --name launcherdb -p 5432:5432 -e POSTGRES_USER=postgres -d postgres

// Capital One constants.
const (
	BaseURL = "https://3hkaob4gkc.execute-api.us-east-1.amazonaws.com/prod/au-hackathon"
)

// main pull the customers listed below and their transactions and
// stores them in the database.
func (m *Main) seedDatabase() {
	launchers := []launcher.Launcher{
		{CustomerID: 100120000, AccountID: 100100000},
		{CustomerID: 100220000, AccountID: 100200000},
		{CustomerID: 100240000, AccountID: 100200000},
		{CustomerID: 100530000, AccountID: 100500000},
		{CustomerID: 100930000, AccountID: 100900000},
	}

	var testID int
	if err := m.db.QueryRow(`SELECT id FROM launchers LIMIT 1`).Scan(&testID); err != nil {
		panic(err)
	} else if testID != 0 {
		return // do not run script if db is already seeded
	}

	for i, cust := range launchers {
		var jsonStr = []byte(fmt.Sprintf(`{"customer_id": %v}`, cust.CustomerID))
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", BaseURL, "customers"), bytes.NewBuffer(jsonStr))
		if err != nil {
			panic(err)
		}
		resp, err := m.client.Do(req)
		if err != nil {
			panic(err)
		} else if resp.StatusCode != http.StatusOK {
			panic(fmt.Sprintf("status code: %#v", resp))
		}
		defer resp.Body.Close()

		var data []struct {
			Customers []launcher.Customer
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			panic(err)
		} else if data[0].Customers[0].IsPrimary {
			fmt.Printf("skipped primary customer id: %v", data[0].Customers[0].ID)
		}

		c := data[0].Customers[0]
		interestRate := rand.Float64() * 20
		if err := m.db.QueryRow(`INSERT INTO launchers (
        customer_id
      , account_id
      , first_name
      , last_name
			, interest_rate
      , credit_limit
      , due_date
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
      RETURNING id`,
			cust.CustomerID,
			cust.AccountID,
			c.FirstName,
			c.LastName,
			interestRate,
			1000.00,
			time.Now(),
		).Scan(&launchers[i].ID); err != nil {
			fmt.Println(err)
		}

		if err := m.getTransactions(launchers[i].ID, launchers[i].CustomerID, interestRate); err != nil {
			panic(err)
		}
	}
	fmt.Println("capital one data seeded in database")
}

// getTransactions retrieves transactions for the provided ids.
func (m *Main) getTransactions(id int, customerID int64, interestRate float64) error {
	var jsonStr = []byte(fmt.Sprintf(`{"customer_id": %v}`, customerID))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", BaseURL, "transactions"), bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	resp, err := m.client.Do(req)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("status code: %#v", resp))
	}
	defer resp.Body.Close()

	var data []struct {
		Customers []struct {
			CustomerID   int64                  `json:"customer_id"`
			Transactions []launcher.Transaction `json:"transactions"`
		} `json:"customers"`
		AccountID int64 `json:"account_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	var totalRewards, totalBalance float64
	for _, trans := range data[0].Customers[0].Transactions {
		// date, _ := time.Parse("January", trans.Month)
		if _, err := m.db.Exec(`INSERT INTO transactions (
	      id
	    , launchers_id
	    , type
	    , merchant
	    , amount
	    , purchase_date
	  ) VALUES ($1, $2, $3, $4, $5, $6)`,
			trans.ID,
			id,
			launcher.TransactionTypeCharge,
			trans.Merchant,
			trans.Amount/100,
			time.Date(trans.Year, 10, trans.Day, 0, 0, 0, 0, time.UTC),
		); err != nil {
			fmt.Println(err)
		}
		totalRewards += trans.RewardsEarned
		totalBalance += trans.Amount / 100
	}

	minPayment := (totalBalance * 0.1) * (1 + interestRate/1200)
	if _, err := m.db.Exec(`UPDATE launchers SET reward_balance = $1, balance = $2, minimum_payment = $3 WHERE id = $4`, totalRewards, totalBalance, minPayment, id); err != nil {
		fmt.Println(err)
	}

	return nil
}
