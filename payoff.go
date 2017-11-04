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
}

// Calculate calculates either the total months or monthly payment for the payoff.
func (p *Payoff) Calculate() error {
	if p.TotalMonths > 0 && p.MonthlyPayment > 0 {
		return errors.New("only one of total months or monthly payment may be submitted")
	}

	mpr := p.InterestRate / 1200                    // monthly percentage rate
	if p.TotalMonths > 0 && p.MonthlyPayment == 0 { // Calculate by number of months
		p.MonthlyPayment = p.Balance * mpr / math.Pow(1-1/(1+mpr), float64(p.TotalMonths)) * -1 * (1 + mpr)
	} else if p.TotalMonths == 0 && p.MonthlyPayment > 0 { // Calculate by monthly payment
		p.TotalMonths = int(math.Ceil(math.Log(-1*-p.MonthlyPayment/p.Balance*mpr) / math.Log(1+mpr)))
	}

	p.TotalInterestCost = p.Balance * mpr * float64(p.TotalMonths)
	return nil
}
