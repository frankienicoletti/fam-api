package launcher_test

import (
	"testing"

	"github.com/katenicoletti/fam-api"
)

func TestPayoff_MonthlyPayment(t *testing.T) {
	t.Skip()
	p := launcher.Payoff{
		Balance:        100,
		MonthlyPayment: 10,
		InterestRate:   0.01,
	}
	if err := p.Calculate(); err != nil {
		t.Fatal(err)
	} else if p.TotalMonths != 11 {
		t.Fatalf("unexpected total months: %d", p.TotalMonths)
	}
}

func TestPayoff_MonthlyPayment_ZeroInterest(t *testing.T) {
	t.Skip()
	p := launcher.Payoff{
		Balance:        100,
		MonthlyPayment: 10,
		InterestRate:   0,
	}
	if err := p.Calculate(); err != nil {
		t.Fatal(err)
	} else if p.TotalMonths != 10 {
		t.Fatalf("unexpected total months: %d", p.TotalMonths)
	}
}

func TestPayoff_TotalMonths(t *testing.T) {
	p := launcher.Payoff{
		Balance:      100,
		TotalMonths:  10,
		InterestRate: 0.01,
	}
	if err := p.Calculate(); err != nil {
		t.Fatal(err)
	} else if p.MonthlyPayment != 10.000458339017815 {
		t.Fatalf("%#v", p)
	}
}
