package runtime

import (
	"io/ioutil"
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

func (t *Runtime) Output(ocr *model.OcrContext) error {

	if err := pather.DirEnsure(pather.Dir(ocr.DestCsv)); err != nil {
		return err
	}

	if err := util.CopyFile(ocr.DestCsv, ocr.WorkCsv); err != nil {
		return err
	}

	if err := util.CopyFile(ocr.DestPdf, ocr.WorkPdf); err != nil {
		return err
	}

	if err := util.CopyFile(ocr.DestYaml, ocr.WorkYaml); err != nil {
		return err
	}

	return nil
}

func (t *Runtime) DumpCvs(ocr *model.OcrContext) error {
	ocr.WorkCsv = extentioner.Swap(ocr.WorkPdf, t.extPdf, t.extCvs)
	file, err := os.OpenFile(ocr.WorkCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)

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

func (t *Runtime) DumpYaml(ocr *model.OcrContext) error {
	ocr.WorkYaml = extentioner.Swap(ocr.WorkPdf, t.extPdf, t.extYaml)
	ocr.DestCsv = pather.Join(t.userConfig.OutputPath, pather.Base(ocr.WorkCsv))
	ocr.DestPdf = pather.Join(t.userConfig.OutputPath, pather.Base(ocr.WorkPdf))
	ocr.DestYaml = pather.Join(t.userConfig.OutputPath, pather.Base(ocr.WorkYaml))

	bs, err := yaml.Marshal(ocr)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(ocr.WorkYaml, bs, 0755); err != nil {
		return err
	}

	return nil
}
