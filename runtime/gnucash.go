package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/go-commoner/optioner"
)

func (t *Runtime) gnuCash(ocr *model.ProcessContext) []model.GnuCashTransaction {
	switch ocr.UserConfig.Processor {
	case "chase-checking":
		return t.gnuCashChecking(ocr)
	case "chase-credit-card":
		var gtxs []model.GnuCashTransaction
		return gtxs
	default:
		var gtxs []model.GnuCashTransaction
		return gtxs
	}
}

func (t *Runtime) gnuCashCredit(ocr *model.ProcessContext) []model.GnuCashTransaction {
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

	for _, tx := range ocr.CreditCard.Payments {
		incoming(tx, "Payments")
	}

	for _, tx := range ocr.CreditCard.Purchases {
		outgoing(tx, "Purchases")
	}

	for _, tx := range ocr.CreditCard.Redemptions {
		incoming(tx, "Redemptions")
		outgoing(tx, "Redemptions")
	}

	return gtxs
}

func (t *Runtime) gnuCashChecking(ocr *model.ProcessContext) []model.GnuCashTransaction {
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

	for _, tx := range ocr.Checking.Deposits {
		incoming(tx, "DEPOSITS AND ADDITIONS")
	}

	for _, tx := range ocr.Checking.Checks {
		outgoing(tx, "CHECKS PAID")
	}

	for _, tx := range ocr.Checking.Withdrawals {
		outgoing(tx, "ELECTRONIC WITHDRAWALS")
	}

	for _, tx := range ocr.Checking.AtmDebit {
		outgoing(tx, "ATM & DEBIT CARD WITHDRAWALS")
	}

	return gtxs
}
