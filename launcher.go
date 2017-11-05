package launcher

import "time"

// Internal indicates that the value is calculated by this API, not Capital One.

// Launcher defines a "Launcher" kid's credit card.
type Launcher struct {
	ID             int       `json:"id"`
	CustomerID     int64     `json:"customer_id"`
	AccountID      int64     `json:"account_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CreditLimit    int       `json:"credit_limit"`    // Internal
	Balance        float64   `json:"balance"`         // Internal
	DueDate        time.Time `json:"due_date"`        // Internal
	MinimumPayment float64   `json:"minimum_payment"` // Internal
}

// Customer defines an account primary or authorized user.
type Customer struct {
	ID        int64  `json:"customers>customer_id"`
	AccountID int64  `json:"account_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsPrimary bool   `json:"is_primary"`
}

// Account defines a credit card account.
type Account struct {
	ID              int64   `json:"account_id"`
	CardType        string  `json:"card_type"`
	RewardsEarned   float64 `json:"total_rewards_earned"`
	RewardsUsed     float64 `json:"total_rewards_used"`
	CreditLimit     float64 `json:"credit_limit"`
	Balance         float64 `json:"balance"`
	AuthorizedUsers []Customer
	PrimaryUser     Customer
}

// Payment defines a credit card payment.
type Payment struct {
	ID               int64
	CustomerID       int64 // Internal
	Month            int
	Year             int
	Amount           float64
	BalanceRemaining float64
}

// Transaction type constants.
const (
	TransactionTypeCharge   = "charge"
	TransactionTypePayment  = "payment"
	TransactionTypeInterest = "interest" // Internal
	TransactionTypeLateFee  = "late_fee" // Internal
)

// Transaction defines a credit card transaction.
type Transaction struct {
	ID            int64     `json:"transaction_id"`
	LauncherID    int       `json:"launcher_id"` // Internal
	CustomerID    int64     `json:"customer_id"`
	Type          string    `json:"type"` // Internal
	Merchant      string    `json:"merchant_name"`
	Amount        float64   `json:"amount"`
	Day           int       `json:"day,omitempty"`
	Month         string    `json:"month,omitempty"`
	Year          int       `json:"year,omitempty"`
	Date          time.Time `json:"date"`
	RewardsEarned float64   `json:"rewards_earned"`
}
