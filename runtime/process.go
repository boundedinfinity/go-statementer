package runtime

import (
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/optioner"
	"github.com/gocarina/gocsv"
)

func (t *Runtime) Process(ocr *model.OcrContext) error {
	manager := processors.NewManager(t.logger, t.userConfig, ocr)
	classifier, err := manager.GetClassifier()

	if err != nil {
		return err
	}

	if err := manager.Init(classifier); err != nil {
		return err
	}

	if err := manager.Extract(classifier); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	statement, err := manager.Lookup()

	if err != nil {
		return err
	}

	if err := manager.Init(statement); err != nil {
		return err
	}

	if err := manager.Extract(statement); err != nil {
		return err
	}

	if err := manager.Transform(statement); err != nil {
		return err
	}

	// processor.Print()

	return nil
}

func (t *Runtime) gnuCash(ocr *model.OcrContext) []model.GnuCashTransaction {
	var gtxs []model.GnuCashTransaction

	for _, tx := range ocr.Statement.Deposits {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate(tx.Date)
		gtx.Description = tx.Memo
		gtx.AccountName = optioner.OfZ(ocr.UserConfig.Deposits).OrElse(model.CHECKING_INCOMING_ACCOUNT)
		gtx.Incoming = model.GnuCashFloat(tx.Amount)
		gtxs = append(gtxs, gtx)
	}

	for _, tx := range ocr.Statement.Checks {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate(tx.Date)
		gtx.Description = tx.Memo
		gtx.AccountName = optioner.OfZ(ocr.UserConfig.Checks).OrElse(model.CHECKING_OUTGOING_ACCOUNT)
		gtx.Outgoing = model.GnuCashFloat(tx.Amount)
		gtxs = append(gtxs, gtx)
	}

	for _, tx := range ocr.Statement.Withdrawals {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate(tx.Date)
		gtx.Description = tx.Memo
		gtx.AccountName = optioner.OfZ(ocr.UserConfig.Withdrawals).OrElse(model.CHECKING_OUTGOING_ACCOUNT)
		gtx.Outgoing = model.GnuCashFloat(tx.Amount)
		gtxs = append(gtxs, gtx)
	}

	return gtxs
}

func (t *Runtime) DumpCvs(ocr *model.OcrContext) error {
	ocr.Csv = extentioner.Swap(ocr.Pdf, t.extPdf, t.extCvs)
	file, err := os.OpenFile(ocr.Csv, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	txs := t.gnuCash(ocr)

	if err := gocsv.Marshal(txs, file); err != nil {
		return err
	}

	return nil
}
