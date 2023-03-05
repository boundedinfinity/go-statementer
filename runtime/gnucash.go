package runtime

import (
	"math"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/go-commoner/optioner"
)

func (t *Runtime) gnuCash(pc *model.ProcessContext) []model.GnuCashTransaction {
	switch pc.UserConfig.Processor {
	case "chase-checking":
		return t.gnuCashChecking(pc)
	case "chase-credit-card":
		return t.gnuCashCredit(pc)
	default:
		var gtxs []model.GnuCashTransaction
		return gtxs
	}
}

func abs(f float32) float32 {
	return float32(math.Abs(float64(f)))
}

func (t *Runtime) gnuCashCredit(pc *model.ProcessContext) []model.GnuCashTransaction {
	var gtxs []model.GnuCashTransaction

	incoming := func(conanical model.Transaction, notes string) {
		gnucash := model.NewGnuCashTrasaction()
		gnucash.Date = model.GnuCashDate(conanical.Date)
		gnucash.Description = conanical.Memo
		gnucash.AccountName = optioner.OfZ(pc.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gnucash.Incoming = model.GnuCashFloat(abs(conanical.Amount))
		gnucash.Notes = notes
		gtxs = append(gtxs, gnucash)
	}

	outgoing := func(conanical model.Transaction, notes string) {
		gnucash := model.NewGnuCashTrasaction()
		gnucash.CheckNumber = conanical.Number
		gnucash.Date = model.GnuCashDate(conanical.Date)
		gnucash.Description = conanical.Memo
		gnucash.AccountName = optioner.OfZ(pc.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gnucash.Outgoing = model.GnuCashFloat(abs(conanical.Amount))
		gnucash.Notes = notes
		gtxs = append(gtxs, gnucash)
	}

	for _, tx := range pc.CreditCard.Payments {
		incoming(tx, "Payments")
	}

	for _, tx := range pc.CreditCard.Purchases {
		outgoing(tx, "Purchases")
	}

	for _, tx := range pc.CreditCard.Redemptions {
		incoming(tx, "Redemptions")
		outgoing(tx, "Redemptions")
	}

	return gtxs
}

func (t *Runtime) gnuCashChecking(pc *model.ProcessContext) []model.GnuCashTransaction {
	var gtxs []model.GnuCashTransaction

	incoming := func(conanical model.Transaction, notes string) {
		gnucash := model.NewGnuCashTrasaction()
		gnucash.Date = model.GnuCashDate(conanical.Date)
		gnucash.Description = conanical.Memo
		gnucash.AccountName = optioner.OfZ(pc.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gnucash.Incoming = model.GnuCashFloat(conanical.Amount)
		gnucash.Notes = notes
		gtxs = append(gtxs, gnucash)
	}

	outgoing := func(conanical model.Transaction, notes string) {
		gnucash := model.NewGnuCashTrasaction()
		gnucash.CheckNumber = conanical.Number
		gnucash.Date = model.GnuCashDate(conanical.Date)
		gnucash.Description = conanical.Memo
		gnucash.AccountName = optioner.OfZ(pc.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gnucash.Outgoing = model.GnuCashFloat(conanical.Amount)
		gnucash.Notes = notes
		gtxs = append(gtxs, gnucash)
	}

	for _, tx := range pc.Checking.Deposits {
		incoming(tx, "DEPOSITS AND ADDITIONS")
	}

	for _, tx := range pc.Checking.Checks {
		outgoing(tx, "CHECKS PAID")
	}

	for _, tx := range pc.Checking.Withdrawals {
		outgoing(tx, "ELECTRONIC WITHDRAWALS")
	}

	for _, tx := range pc.Checking.AtmDebit {
		outgoing(tx, "ATM & DEBIT CARD WITHDRAWALS")
	}

	return gtxs
}
