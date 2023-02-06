package model

import (
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/boundedinfinity/go-commoner/stringer"
	"github.com/boundedinfinity/rfc3339date"
)

type gnuCashDate struct {
	rfc3339date.Rfc3339Date
}

func (t gnuCashDate) String() string {
	t.Rfc3339Date.Format("02/01/2006")
	return ""
}

type GnuCashTransation struct {
	Date              gnuCashDate `csv:"Date"`
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
	ReconcileDate     gnuCashDate `csv:"Reconcile Date"`
	RatePrice         float32     `csv:"Rate/Price"`
}

func NewGnuCashTrasaction() GnuCashTransation {
	return GnuCashTransation{
		CommodityCurrency: "CURRENCY::USD",
		RatePrice:         1,
	}
}

func Trasaction2GnuCash(tx Transaction, fullAccountName string) GnuCashTransation {
	gnu := NewGnuCashTrasaction()
	gnu.FullAccountName = fullAccountName
	gnu.AccountName = slicer.Last(stringer.Split(fullAccountName, "::"))
	gnu.Description = tx.Memo
	gnu.Date = gnuCashDate{tx.Date}
	gnu.AmountNum = tx.Amount
	gnu.Number = tx.Number

	return gnu
}
