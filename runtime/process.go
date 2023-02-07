package runtime

import (
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/gocarina/gocsv"
)

func (t *Runtime) Process(ocr *model.OcrContext) error {
	manager := processors.NewManager(t.logger, t.userConfig, ocr)
	classifier, err := manager.GetClassifier(ocr)

	if err != nil {
		return err
	}

	if err := manager.Extract(ocr, classifier); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	processor, err := manager.Lookup(ocr)

	if err != nil {
		return err
	}

	if err := manager.Extract(ocr, processor); err != nil {
		return err
	}

	if err := processor.Convert(); err != nil {
		return err
	}

	processor.Print()

	return nil
}

func (t *Runtime) gnuCash(ocr *model.OcrContext) []model.GnuCashTransaction {
	var gtxs []model.GnuCashTransaction

	for _, tx := range ocr.Statement.Deposits {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate{Rfc3339Date: tx.Date}
		gtx.Description = tx.Memo
		gtx.FullAccountName = model.CHECKING_INCOMING_ACCOUNT
		gtx.AccountName = util.FullAccountName2AccountName(gtx.FullAccountName)
		gtx.AmountNum = tx.Amount
		gtxs = append(gtxs, gtx)
	}

	for _, tx := range ocr.Statement.Checks {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate{Rfc3339Date: tx.Date}
		gtx.Description = tx.Memo
		gtx.FullAccountName = model.CHECKING_OUTGOING_ACCOUNT
		gtx.AccountName = util.FullAccountName2AccountName(gtx.FullAccountName)
		gtx.AmountNum = -tx.Amount
		gtxs = append(gtxs, gtx)
	}

	for _, tx := range ocr.Statement.Withdrawals {
		gtx := model.NewGnuCashTrasaction()
		gtx.Date = model.GnuCashDate{Rfc3339Date: tx.Date}
		gtx.Description = tx.Memo
		gtx.FullAccountName = model.CHECKING_OUTGOING_ACCOUNT
		gtx.AccountName = util.FullAccountName2AccountName(gtx.FullAccountName)
		gtx.AmountNum = -tx.Amount
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
