package runtime

import (
	"io/ioutil"
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

func (t *Runtime) Output(dst, src *model.FileSet) error {
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

func (t *Runtime) DumpCvs(pc *model.ProcessContext) error {
	file, err := os.OpenFile(pc.Stage2.Csv, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	txs := t.gnuCash(pc)

	if err := gocsv.Marshal(txs, file); err != nil {
		return err
	}

	return nil
}

func (t *Runtime) DumpYaml(pc *model.ProcessContext) error {
	name := pather.Base(pc.Stage2.Dir)

	if err := t.CalcFiles(t.UserConfig.OutputPath, name, &pc.Dest, pc.Stage2); err != nil {
		return err
	}

	bs, err := yaml.Marshal(pc)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(pc.Stage2.Yaml, bs, 0755); err != nil {
		return err
	}

	return nil
}
