package model

import (
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/boundedinfinity/rfc3339date"
)

type Config struct {
	InputPath   string                `yaml:"inputPath"`
	OutputPath  string                `yaml:"outputPath"`
	WorkPath    string                `yaml:"workPath"`
	SumExt      string                `yaml:"sumExt"`
	InputExt    string                `yaml:"inputExt"`
	Debug       bool                  `yaml:"debug"`
	IgnorePaths []string              `yaml:"ignorePaths"`
	Statements  []StatementDescriptor `yaml:"statements"`
}

func (t *Config) FindAccount(id string) (StatementDescriptor, bool) {
	return slicer.FindFn(t.Statements, func(s StatementDescriptor) bool {
		return s.Account == id || s.Name == id
	})
}

type StatementDescriptor struct {
	Name      string `yaml:"name"`
	Account   string `yaml:"account"`
	Processor string `yaml:"processor"`
}

type CheckingStatement struct {
	Source         string        `yaml:"source"`
	AccountNumber  string        `yaml:"accountNumber"`
	OpeningBalance float32       `yaml:"openingBalance"`
	ClosingBalance float32       `yaml:"closingBalance"`
	OpeningDate    float32       `yaml:"openingDate"`
	ClosingDate    float32       `yaml:"closingDate"`
	Deposits       []Transaction `yaml:"deposits"`
	Checks         []Transaction `yaml:"checks"`
	Withdrawals    []Transaction `yaml:"withdrawals"`
}

type ChaseCreditCardStatement struct {
	Source          string        `yaml:"source"`
	AccountNumber   string        `yaml:"accountNumber"`
	PreviousBalance float32       `yaml:"previousBalance"`
	NewBalance      float32       `yaml:"newBalance"`
	OpeningDate     float32       `yaml:"openingDate"`
	ClosingDate     float32       `yaml:"closingDate"`
	Payments        []Transaction `yaml:"payments"`
	Purchases       []Transaction `yaml:"purchases"`
	Redemptions     []Transaction `yaml:"redemptions"`
}

type Transaction struct {
	Date   rfc3339date.Rfc3339Date `yaml:"date"`
	Memo   string                  `yaml:"memo"`
	Amount float32                 `yaml:"amount"`
}
