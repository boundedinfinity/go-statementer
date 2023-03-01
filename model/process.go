package model

import "github.com/boundedinfinity/rfc3339date"

type StatementProcessor interface {
	Account() string
	ClosingDate() rfc3339date.Rfc3339Date
}

type ProcessContext struct {
	Stage1     FileSet     `yaml:"stage1"`
	Stage2     FileSet     `yaml:"stage2"`
	Dest       FileSet     `yaml:"dest"`
	Extracted  []Extracted `yaml:"extracted"`
	UserConfig UserConfigStatement
	Checking   CheckingStatement   `yaml:"checking"`
	CreditCard CreditCardStatement `yaml:"creditcard"`
}

func NewProcessContext() *ProcessContext {
	return &ProcessContext{
		Stage1:    NewFileSet(),
		Stage2:    NewFileSet(),
		Dest:      NewFileSet(),
		Extracted: make([]Extracted, 0),
	}
}

type Extracted struct {
	Name   string            `yaml:"name"`
	Values map[string]string `yaml:"values"`
}
