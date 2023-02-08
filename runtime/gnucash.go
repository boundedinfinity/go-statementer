package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/go-commoner/optioner"
)

func (t *Runtime) gnuCash(ocr *model.OcrContext) []model.GnuCashTransaction {
	var gtxs []model.GnuCashTransaction

	for _, tx := range ocr.Statement.Deposits {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate(tx.Date)
		gtx.Description = tx.Memo
		gtx.AccountName = optioner.OfZ(ocr.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gtx.Incoming = model.GnuCashFloat(tx.Amount)
		gtxs = append(gtxs, gtx)
	}

	for _, tx := range ocr.Statement.Checks {
		gtx := model.NewGnuCashTrasaction()
		gtx.CheckNumber = tx.Number
		gtx.Date = model.GnuCashDate(tx.Date)
		gtx.Description = tx.Memo
		gtx.AccountName = optioner.OfZ(ocr.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gtx.Outgoing = model.GnuCashFloat(tx.Amount)
		gtxs = append(gtxs, gtx)
	}

	for _, tx := range ocr.Statement.Withdrawals {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate(tx.Date)
		gtx.Description = tx.Memo
		gtx.AccountName = optioner.OfZ(ocr.UserConfig.Name).OrElse(model.IMPORTED_UNKOWN)
		gtx.Outgoing = model.GnuCashFloat(tx.Amount)
		gtxs = append(gtxs, gtx)
	}

	return gtxs
}
