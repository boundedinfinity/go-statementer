package model

import (
	"github.com/boundedinfinity/rfc3339date"
)

const (
	CHECKING_INCOMING_ACCOUNT = "Equity::Imported"
	CHECKING_OUTGOING_ACCOUNT = "Expenses::Imported"
)

type GnuCashDate struct {
	rfc3339date.Rfc3339Date
}

func (t GnuCashDate) String() string {
	return t.Rfc3339Date.String()
}

func (t *GnuCashDate) MarshalCSV() (string, error) {
	if t.Rfc3339Date.IsZero() {
		return "", nil
	} else {
		return t.Rfc3339Date.Format("02/01/2006"), nil
	}
}

type GnuCashTransaction struct {
	Date              GnuCashDate `csv:"Date"`
	TransactionID     string      `csv:"Transaction ID"`
	Number            string      `csv:"Number"`
	Description       string      `csv:"Description"`
	Notes             string      `csv:"Notes"`
	CommodityCurrency string      `csv:"Commodity/Currency"`
	VoidReason        string      `csv:"Void Reason"`
	Action            string      `csv:"Action"`
	Memo              string      `csv:"Memo"`
	FullAccountName   string      `csv:"Full Account Name"`
	AccountName       string      `csv:"Account Name"`
	AmountWithSym     float32     `csv:"Amount With Sym"`
	AmountNum         float32     `csv:"Amount Num."`
	Reconcile         string      `csv:"Reconcile"`
	ReconcileDate     GnuCashDate `csv:"Reconcile Date"`
	RatePrice         float32     `csv:"Rate/Price"`
}

func NewGnuCashTrasaction() GnuCashTransaction {
	return GnuCashTransaction{
		CommodityCurrency: "CURRENCY::USD",
		RatePrice:         1,
	}
}
