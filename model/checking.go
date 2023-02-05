package model

import (
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/boundedinfinity/rfc3339date"
	"github.com/oriser/regroup"
)

type MatchHandler func(map[string]string)

type FieldDescriptor struct {
	Name string
	Key  string
}

func NewFieldAndKey(name, key string) *FieldDescriptor {
	return &FieldDescriptor{
		Name: name,
		Key:  key,
	}
}

func NewField(name string) *FieldDescriptor {
	return NewFieldAndKey(name, name)
}

type LineDescriptor struct {
	Name    string
	Pattern string
	Fields  []*FieldDescriptor
	Regex   *regroup.ReGroup
}

func NewLine(name, pattern string, fields ...*FieldDescriptor) *LineDescriptor {
	return &LineDescriptor{
		Name:    name,
		Pattern: pattern,
		Fields:  fields,
	}
}

func NewLineWithField(name, pattern string) *LineDescriptor {
	return NewLine(name, pattern, NewField(name))
}

func NewLineWithFieldAndKey(name, key, pattern string) *LineDescriptor {
	return NewLine(name, pattern, NewFieldAndKey(name, key))
}

// type SectionDescriptor struct {
// 	Name  string
// 	Start *LineDescriptor
// 	End   *LineDescriptor
// 	Line  *LineDescriptor
// }

// func (t SectionDescriptor) Lines() []*LineDescriptor {
// 	return []*LineDescriptor{t.Start, t.End, t.Line}
// }

type CheckingDescriptor struct {
	Account                *LineDescriptor
	OpeningBalance         *LineDescriptor
	ClosingBalance         *LineDescriptor
	OpeningDate            *LineDescriptor
	ClosingDate            *LineDescriptor
	DepositsBalance        *LineDescriptor
	DepositsStart          *LineDescriptor
	DepositsEnd            *LineDescriptor
	DepositsTransaction    *LineDescriptor
	ChecksBalance          *LineDescriptor
	ChecksStart            *LineDescriptor
	ChecksEnd              *LineDescriptor
	WithdrawalsBalance     *LineDescriptor
	WithdrawalsStart       *LineDescriptor
	WithdrawalsEnd         *LineDescriptor
	Transaction            *LineDescriptor
	Page                   *LineDescriptor
	AnnualPercentageEarned *LineDescriptor
	InterestEarned         *LineDescriptor
	InterestPaid           *LineDescriptor
}

func (t CheckingDescriptor) Lines() []*LineDescriptor {
	lines := make([]*LineDescriptor, 0)
	lines = append(lines, t.Account, t.Transaction, t.Page,
		t.OpeningBalance, t.OpeningDate,
		t.ClosingBalance, t.ClosingDate,
		t.WithdrawalsBalance, t.WithdrawalsStart,
		t.ChecksBalance, t.ChecksStart, t.ChecksEnd,
		t.DepositsBalance, t.DepositsStart, t.DepositsEnd,
		t.AnnualPercentageEarned, t.InterestEarned, t.InterestPaid,
	)

	return slicer.Filter(lines, func(line *LineDescriptor) bool { return line != nil })
}

type CheckingStatement struct {
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

type CheckingMatchers struct {
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

func (t CheckingMatchers) All() []MatcherConfig {
	return []MatcherConfig{
		t.Account,
		t.OpeningBalance, t.OpeningDate,
		t.ClosingBalance, t.ClosingDate,
		t.ChecksBalance, t.ChecksStart, t.ChecksEnd,
		t.WithdrawalsBalance, t.WithdrawalsStart, t.WithdrawalsEnd,
		t.DepositsBalance, t.DepositsStart, t.DepositsEnd,
	}
}
