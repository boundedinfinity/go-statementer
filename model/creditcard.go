package model

type CreditCardStatement struct {
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
