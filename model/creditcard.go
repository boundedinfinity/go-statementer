package model

import "github.com/boundedinfinity/rfc3339date"

type CreditCardStatement struct {
	Source         string                  `yaml:"source"`
	Account        string                  `yaml:"account"`
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

type CreditCardStatementRaw struct {
	Source         string           `yaml:"source"`
	Account        string           `yaml:"account"`
	OpeningBalance string           `yaml:"openingBalance"`
	ClosingBalance string           `yaml:"closingBalance"`
	OpeningDate    string           `yaml:"openingDate"`
	ClosingDate    string           `yaml:"closingDate"`
	Payments       []TransactionRaw `yaml:"payments"`
	Purchases      []TransactionRaw `yaml:"purchases"`
	Redemptions    []TransactionRaw `yaml:"redemptions"`
}

func NewCreditCardStatementRaw() CreditCardStatementRaw {
	return CreditCardStatementRaw{
		Payments:    make([]TransactionRaw, 0),
		Purchases:   make([]TransactionRaw, 0),
		Redemptions: make([]TransactionRaw, 0),
	}
}
