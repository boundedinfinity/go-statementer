package model

import "github.com/boundedinfinity/rfc3339date"

type Transaction struct {
	Number string                  `yaml:"number"`
	Date   rfc3339date.Rfc3339Date `yaml:"date"`
	Memo   string                  `yaml:"memo"`
	Amount float32                 `yaml:"amount"`
}
