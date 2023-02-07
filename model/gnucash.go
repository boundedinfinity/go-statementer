package model

import (
	"fmt"

	"github.com/boundedinfinity/rfc3339date"
)

const (
	CHECKING_INCOMING_ACCOUNT = "Equity::Imported"
	CHECKING_OUTGOING_ACCOUNT = "Expenses::Imported"
)

type GnuCashDate rfc3339date.Rfc3339Date

func (t GnuCashDate) String() string {
	return rfc3339date.Rfc3339Date(t).String()
}

func (t GnuCashDate) MarshalCSV() (string, error) {
	if t.IsZero() {
		return "", nil
	} else {
		return t.Format("2006-01-02"), nil
	}
}

type GnuCashFloat float32

func (t *GnuCashFloat) MarshalCSV() (string, error) {
	return fmt.Sprintf("%.2f", *t), nil
}

type GnuCashRate float32

func (t *GnuCashRate) MarshalCSV() (string, error) {
	return fmt.Sprintf("%.4f", *t), nil
}

type GnuCashTransaction struct {
	Date          GnuCashDate  `csv:"Date"`
	TransactionID string       `csv:"Transaction ID"`
	CheckNumber   string       `csv:"CheckNumber"`
	Memo          string       `csv:"Memo"`
	Description   string       `csv:"Description"`
	Notes         string       `csv:"Notes"`
	AccountName   string       `csv:"AccountName"`
	Incoming      GnuCashFloat `csv:"Incoming"`
	Outgoing      GnuCashFloat `csv:"Outgoing"`
}

func NewGnuCashTrasaction() GnuCashTransaction {
	return GnuCashTransaction{}
}
