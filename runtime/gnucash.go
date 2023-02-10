package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/go-commoner/optioner"
)

func (t *Runtime) gnuCash(ocr *model.OcrContext) []model.GnuCashTransaction {
	var gtxs []model.GnuCashTransaction

	incoming := func(conanical model.Transaction, notes string) {
		gnucash := model.NewGnuCashTrasaction()
		gnucash.Date = model.GnuCashDate(conanical.Date)
		gnucash.Description = conanical.Memo
		gnucash.AccountName = optioner.OfZ(ocr.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gnucash.Incoming = model.GnuCashFloat(conanical.Amount)
		gnucash.Notes = notes
		gtxs = append(gtxs, gnucash)
	}

	outgoing := func(conanical model.Transaction, notes string) {
		gnucash := model.NewGnuCashTrasaction()
		gnucash.CheckNumber = conanical.Number
		gnucash.Date = model.GnuCashDate(conanical.Date)
		gnucash.Description = conanical.Memo
		gnucash.AccountName = optioner.OfZ(ocr.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gnucash.Outgoing = model.GnuCashFloat(conanical.Amount)
		gnucash.Notes = notes
		gtxs = append(gtxs, gnucash)
	}

	for _, tx := range ocr.Statement.Deposits {
		incoming(tx, "DEPOSITS AND ADDITIONS")
	}

	for _, tx := range ocr.Statement.Checks {
		outgoing(tx, "CHECKS PAID")
	}

	for _, tx := range ocr.Statement.Withdrawals {
		outgoing(tx, "ELECTRONIC WITHDRAWALS")
	}

	for _, tx := range ocr.Statement.AtmDebit {
		outgoing(tx, "ATM & DEBIT CARD WITHDRAWALS")
	}

	return gtxs
}
