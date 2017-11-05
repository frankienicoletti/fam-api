package launcher

import (
	"errors"
	"math"
)

// Payoff defines the calculations for repayment.
type Payoff struct {
	// Account balance or hypothetical amount of a future purchase.
	Balance float64 `json:"balance"`

	// Interest rate out of 100 (%18 => 18.00).
	InterestRate float64 `json:"interest_rate"`

	// Total months to pay off balance.
	TotalMonths int `json:"total_months"`

	// Monthly payment to pay off balance.
	MonthlyPayment float64 `json:"monthly_payment"`

	// Total interest paid over the life of the balance.
	TotalInterestCost float64 `json:"total_interest_cost"`

	// Graph is an array of GraphData structs used for creating charts on the front end.
	Graph []GraphData `json:"graph"`
}

// GraphData represents a row of graph data.
// (Principal + interest = previous month balance)
type GraphData struct {
	Principal float64 `json:"principal"`
	Interest  float64 `json:"interest"`
	Balance   float64 `json:"balance"`
	Month     int     `json:"month"`
}

// Calculate calculates either the total months or monthly payment for the payoff.
func (p *Payoff) Calculate() error {
	if p.TotalMonths > 0 && p.MonthlyPayment > 0 {
		return errors.New("only one of total months or monthly payment may be submitted")
	}

	mpr := p.InterestRate / 1200                    // monthly percentage rate
	if p.TotalMonths > 0 && p.MonthlyPayment == 0 { // Calculate by number of months
		if p.InterestRate <= 0 {
			return errors.New("interest rate must be greater than zero to calculate monthly payment")
		}
		// Calculate monthly payment.
		p.MonthlyPayment = p.Balance * mpr / (1 - 1/math.Pow(1+mpr, float64(p.TotalMonths)))
	}

	// Build graph.
	currentPrinciple := p.Balance
	month := 1
	for currentPrinciple > 0 {
		var gd GraphData
		gd.Interest = currentPrinciple * mpr          // current monthly interest
		gd.Principal = p.MonthlyPayment - gd.Interest // principal paid each month
		gd.Balance = currentPrinciple - gd.Principal  // balance after payment
		gd.Month = month
		currentPrinciple = gd.Balance

		// Last payment
		if gd.Balance < 0 {
			gd.Balance = 0
			gd.Principal = p.Graph[len(p.Graph)-1].Balance
		}

		p.Graph = append(p.Graph, gd)
		month++
	}
	if p.TotalMonths == 0 {
		p.TotalMonths = len(p.Graph)
	}

	p.TotalInterestCost = p.Balance * mpr * float64(p.TotalMonths)
	return nil
}
