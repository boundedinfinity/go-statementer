package model

import "github.com/boundedinfinity/rfc3339date"

type MatcherConfig struct {
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
}

type CheckingStatement struct {
	Source             string                  `yaml:"source"`
	AccountNumber      string                  `yaml:"accountNumber"`
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
}

type CheckingStatementRaw struct {
	Source             string           `yaml:"source"`
	Account            string           `yaml:"accountNumber"`
	OpeningBalance     string           `yaml:"openingBalance"`
	ClosingBalance     string           `yaml:"closingBalance"`
	OpeningDate        string           `yaml:"openingDate"`
	ClosingDate        string           `yaml:"closingDate"`
	DepositsBalance    string           `yaml:"depositsBalance"`
	Deposits           []TransactionRaw `yaml:"deposits"`
	ChecksBalance      string           `yaml:"checksBalance"`
	Checks             []TransactionRaw `yaml:"checks"`
	WithdrawalsBalance string           `yaml:"withdrawalsBalance"`
	Withdrawals        []TransactionRaw `yaml:"withdrawals"`
}

type CheckingConfig struct {
	Account            MatcherConfig `yaml:"accountNumber"`
	OpeningBalance     MatcherConfig `yaml:"openingBalance"`
	ClosingBalance     MatcherConfig `yaml:"closingBalance"`
	OpeningDate        MatcherConfig `yaml:"openingDate"`
	ClosingDate        MatcherConfig `yaml:"closingDate"`
	Transaction        MatcherConfig `yaml:"transaction"`
	DepositsBalance    MatcherConfig `yaml:"depositsBalance"`
	DepositsStart      MatcherConfig `yaml:"depositsStart"`
	DepositsEnd        MatcherConfig `yaml:"depositsEnd"`
	ChecksBalance      MatcherConfig `yaml:"checksBalance"`
	ChecksStart        MatcherConfig `yaml:"checksStart"`
	ChecksEnd          MatcherConfig `yaml:"checksEnd"`
	WithdrawalsBalance MatcherConfig `yaml:"withdrawalsBalance"`
	WithdrawalsStart   MatcherConfig `yaml:"withdrawalsStart"`
	WithdrawalsEnd     MatcherConfig `yaml:"withdrawalsEnd"`
}

type TransactionRaw struct {
	Date   string
	Memo   string
	Amount string
}
