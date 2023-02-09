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

func (t *Runtime) Output(dst, src *model.ProcessStage) error {

	if err := pather.DirEnsure(dst.Dir); err != nil {
		return err
	}

	if err := util.CopyFile(dst.Csv, src.Csv); err != nil {
		return err
	}

	if err := util.CopyFile(dst.Pdf, src.Pdf); err != nil {
		return err
	}

	return nil
}

func (t *Runtime) DumpCvs(ocr *model.OcrContext) error {
	ocr.Stage2.Csv = extentioner.Swap(ocr.Stage2.Pdf, t.extPdf, t.extCvs)
	file, err := os.OpenFile(ocr.Stage2.Csv, os.O_RDWR|os.O_CREATE, os.ModePerm)

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
	ocr.Stage2.Yaml = extentioner.Swap(ocr.Stage2.Pdf, t.extPdf, t.extYaml)
	Stage2Stage(t.userConfig.OutputPath, pather.Base(ocr.Stage2.Dir), &ocr.Dest, ocr.Stage2)

	bs, err := yaml.Marshal(ocr)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(ocr.Stage2.Yaml, bs, 0755); err != nil {
		return err
	}

	return nil
}
