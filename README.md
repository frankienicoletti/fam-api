# Launcher
Capital One {FAM} Hackathon backend

### Run

To start a postgres docker container:
```
$ docker run --name launcherdb -p 5432:5432 -e POSTGRES_USER=postgres -d postgres
```

```
$ glide install
$ docker build -t fam-api .
$ docker run -p 4000:6000 fam-api
```

API will run on `localhost:4000`

### API

#### /launcher/{id}
Retrieves a Launcher credit card account.

Accepts an id in the url and returns:
```json
{
    "id": 1,
    "customer_id": 0,
    "account_id": 0,
    "first_name": "Pepillo",
    "last_name": "Bulcroft",
    "credit_limit": 0,
    "balance": 500,
    "due_date": "2017-11-04T20:53:38.550706Z",
    "minimum_payment": 0
}
```

#### /launcher/{id}/transactions
Retrieves transactions for a Launcher credit card.

Accepts an id in the url and returns:
```json
[
  {
      "transaction_id": 100120001,
      "launcher_id": 0,
      "customer_id": 0,
      "type": "charge",
      "merchant_name": "chatbooks",
      "amount": 136.34,
      "date": "2017-09-08T00:00:00Z"
  },
  {
      "transaction_id": 100120002,
      "launcher_id": 0,
      "customer_id": 0,
      "type": "charge",
      "merchant_name": "google *fantasy legend",
      "amount": 87.55,
      "date": "2017-03-15T00:00:00Z"
  }
]
```

#### /payoff
Calculates the payoff of a balance.

Accepts the following body; either total number of months or monthly payment may be provided:
```json
{
	"balance": 100,
	"interest_rate": 10,
	"total_months": 2
}

{
	"balance": 100,
	"interest_rate": 10,
	"monthly_payment": 50.00
}
```

Returns:
```json
{
    "balance": 100,
    "interest_rate": 10,
    "total_months": 2,
    "monthly_payment": 55.00,
    "total_interest_cost": 1.6666666666666667
}
```
