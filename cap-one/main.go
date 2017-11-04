package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/katenicoletti/fam-api"
	"github.com/katenicoletti/fam-api/pg"
)

// Collects and stores data from Capital One hackathon API.
// docker run --name launcherdb -p 5432:5432 -e POSTGRES_USER=postgres -d postgres

// Capital One constants.
const (
	BaseURL = "https://3hkaob4gkc.execute-api.us-east-1.amazonaws.com/prod/au-hackathon"
)

// Main ...
type Main struct {
	db     *sql.DB
	client *http.Client
}

// main pull the customers listed below and their transactions and
// stores them in the database.
func main() {
	m := Main{
		db: pg.DB,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	launchers := []launcher.Launcher{
		{
			CustomerID: 100120000, // Primary 100110000
			AccountID:  100100000,
		},
		{
			CustomerID: 100220000,
			AccountID:  100200000,
		},
		{
			CustomerID: 100240000,
			AccountID:  100200000,
		},
		{
			CustomerID: 100530000,
			AccountID:  100500000,
		},
		{
			CustomerID: 100930000,
			AccountID:  100900000,
		},
	}

	for i, cust := range launchers {
		fmt.Printf("customer id: %v", cust.CustomerID)
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
		if err := m.db.QueryRow(`INSERT INTO launchers (
        customer_id
      , account_id
      , first_name
      , last_name
      , credit_limit
      , balance
      , due_date
    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
      RETURNING id`,
			cust.CustomerID,
			cust.AccountID,
			c.FirstName,
			c.LastName,
			500.00,
			0,
			time.Now(),
		).Scan(&launchers[i].ID); err != nil {
			fmt.Println(err)
		}

		if err := m.getTransactions(launchers[i].ID, launchers[i].CustomerID); err != nil {
			panic(err)
		}
		fmt.Printf("customer id: %v completed", cust.CustomerID)
	}
}

// getTransactions retrieves transactions for the provided ids.
func (m *Main) getTransactions(id int, customerID int64) error {
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

	for i, trans := range data[0].Customers[0].Transactions {
		fmt.Printf("transaction %d for customer id %v", i, customerID)
		date, _ := time.Parse("January", trans.Month)
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
			trans.Amount,
			time.Date(trans.Year, date.Month(), trans.Day, 0, 0, 0, 0, time.UTC),
		); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
