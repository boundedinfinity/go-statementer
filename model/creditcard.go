package model

import "github.com/boundedinfinity/rfc3339date"

type CreditCardStatement struct {
	Source         string                  `yaml:"source"`
	AccountNumber  string                  `yaml:"accountNumber"`
	OpeningBalance float32                 `yaml:"openingBalance"`
	ClosingBalance float32                 `yaml:"closingBalance"`
	OpeningDate    rfc3339date.Rfc3339Date `yaml:"openingDate"`
	ClosingDate    rfc3339date.Rfc3339Date `yaml:"closingDate"`
	Payments       []Transaction           `yaml:"payments"`
	Purchases      []Transaction           `yaml:"purchases"`
	Redemptions    []Transaction           `yaml:"redemptions"`
}

func NewCreditCardStatement() CreditCardStatement {
	return CreditCardStatement{
		Payments:    make([]Transaction, 0),
		Purchases:   make([]Transaction, 0),
		Redemptions: make([]Transaction, 0),
	}
}
