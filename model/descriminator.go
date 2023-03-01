package model

import (
	"time"

	"github.com/boundedinfinity/rfc3339date"
)

type DescriminatorProcessor struct {
	account string `yaml:"accountNumber"`
}

func (t *DescriminatorProcessor) Account() string {
	return t.account
}

func (t *DescriminatorProcessor) ClosingDate() rfc3339date.Rfc3339Date {
	return rfc3339date.NewDate(time.Now())
}

var _ StatementProcessor = &DescriminatorProcessor{}
