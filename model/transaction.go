package model

import "github.com/boundedinfinity/rfc3339date"

type Transaction struct {
	Date   rfc3339date.Rfc3339Date `yaml:"date"`
	Memo   string                  `yaml:"memo"`
	Amount float32                 `yaml:"amount"`
}
