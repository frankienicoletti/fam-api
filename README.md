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
    "first_name": "Pepillo",
    "last_name": "Bulcroft",
    "interest_rate": 12.093205759592392,
    "credit_limit": 1000,
    "balance": 502.53080000000045,
    "due_date": "2017-11-05T03:03:18.788217Z",
    "minimum_payment": 50.7595140304111
}
```

#### /launcher/{id}/transactions
Retrieves transactions for a Launcher credit card.

Accepts an id in the url and returns:
```json
[
  {
      "transaction_id": 100120001,
      "type": "charge",
      "merchant_name": "chatbooks",
      "amount": 1.3634,
      "date": "2017-10-08T00:00:00Z"
  },
  {
      "transaction_id": 100120002,
      "type": "charge",
      "merchant_name": "google *fantasy legend",
      "amount": 0.8755,
      "date": "2017-10-15T00:00:00Z"
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
    "interest_rate": 12.093205759592392,
    "total_months": 23,
    "monthly_payment": 5,
    "total_interest_cost": 23.178644372552082,
    "graph": [
        {
            "principal": 3.992232853367301,
            "interest": 1.0077671466326992,
            "balance": 96.0077671466327,
            "month": 1
        },
        {
            "principal": 4.032465264480614,
            "interest": 0.9675347355193864,
            "balance": 91.97530188215208,
            "month": 2
        }
    ]
}
```
