package model

import (
	"github.com/boundedinfinity/rfc3339date"
)

type CheckingStatement struct {
	Account            string                  `yaml:"accountNumber"`
	OpeningBalance     float32                 `yaml:"openingBalance"`
	ClosingBalance     float32                 `yaml:"closingBalance"`
	OpeningDate        rfc3339date.Rfc3339Date `yaml:"openingDate"`
	ClosingDate        rfc3339date.Rfc3339Date `yaml:"closingDate"`
	DepositsBalance    float32                 `yaml:"depositsBalance"`
	Deposits           []Transaction           `yaml:"deposits"`
	ChecksBalance      float32                 `yaml:"checksBalance"`
	Checks             []Transaction           `yaml:"checks"`
	WithdrawalsBalance float32                 `yaml:"withdrawalsBalance"`
	Withdrawals        []Transaction           `yaml:"withdrawals"`
	AtmDebit           []Transaction           `yaml:"atmDebit"`
}

func NewCheckingStatement() CheckingStatement {
	return CheckingStatement{
		Deposits:    make([]Transaction, 0),
		Checks:      make([]Transaction, 0),
		Withdrawals: make([]Transaction, 0),
	}
}
