package runtime

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/boundedinfinity/rfc3339date"
)

func (t *Runtime) LoadFiles() ([]model.ProcessContext, error) {
	var pcs []model.ProcessContext
	allPaths := make([]string, 0)

	for _, path := range t.UserConfig.InputPaths {
		paths, err := util.GetFilteredFiles(path, t.extPdf)

		if err != nil {
			return pcs, err
		}

		allPaths = append(allPaths, paths...)
	}

	allPaths = slicer.Dedup(allPaths)

	for _, path := range allPaths {
		pc := model.NewProcessContext()
		pc.Stage1.Source = path
		pcs = append(pcs, *pc)
	}

	return pcs, nil
}

func (t *Runtime) CalcFiles(dir, name string, dst *model.FileSet, src model.FileSet) error {
	if dir == "" {
		return errors.New("missing source")
	}

	if src.Source == "" {
		return errors.New("missing source")
	}

	dst.Source = src.Source

	if name == "" {
		name = pather.Base(dst.Source)
		name = extentioner.Strip(name)
	}

	dst.Dir = pather.Join(dir, name)
	dst.Pdf = pather.Join(dst.Dir, extentioner.Join(name, t.extPdf))

	replace1 := func(s *string, ext string) {
		*s = extentioner.Swap(dst.Pdf, t.extPdf, ext)
	}

	replace1(&dst.Image, t.extImage)
	replace1(&dst.Text, t.extText)
	replace1(&dst.Csv, t.extCsv)
	replace1(&dst.Yaml, t.extYaml)

	if src.Pdf != "" {
		srcName := extentioner.Strip(src.Pdf)
		dstName := extentioner.Strip(dst.Pdf)

		replace2 := func(s string) string {
			return strings.ReplaceAll(s, srcName, dstName)
		}

		dst.Images = slicer.Map(src.Images, replace2)
		dst.Texts = slicer.Map(src.Texts, replace2)
	}

	return nil
}

func (t *Runtime) Rename(ocr model.ProcessContext, dst *model.FileSet, src model.FileSet) error {
	var account string
	var closingDate rfc3339date.Rfc3339Date

	switch ocr.UserConfig.Processor {
	case "chase-checking":
		account = ocr.Checking.Account
		closingDate = ocr.Checking.ClosingDate
	case "chase-credit-card":
		account = ocr.CreditCard.Account
		closingDate = ocr.CreditCard.ClosingDate
	default:
		return fmt.Errorf("error transformer for %v", ocr.UserConfig.Account)
	}

	name := account
	name = name[len(name)-4:]
	name += "-" + closingDate.String()

	if err := t.CalcFiles(t.UserConfig.WorkPath, name, dst, src); err != nil {
		return err
	}

	if err := pather.DirEnsure(dst.Dir); err != nil {
		return err
	}

	srcFiles, err := pather.GetFiles(src.Dir)

	if err != nil {
		return err
	}

	for _, srcFile := range srcFiles {
		srcName := extentioner.Strip(src.Pdf)
		dstName := extentioner.Strip(dst.Pdf)
		dstFile := strings.ReplaceAll(srcFile, srcName, dstName)

		if err := os.Rename(srcFile, dstFile); err != nil {
			return err
		}
	}

	if err := util.EnsureDelete(src.Dir); err != nil {
		return nil
	}

	return nil
}

func (t *Runtime) prepareDirectory(stage *model.FileSet) error {
	if t.UserConfig.Reprocess {
		if err := util.EnsureDelete(stage.Dir); err != nil {
			return err
		}
	}

	if err := pather.DirEnsure(stage.Dir); err != nil {
		return err
	}

	if t.UserConfig.Reprocess || !pather.PathExists(stage.Pdf) {
		if err := util.CopyFile(stage.Pdf, stage.Source); err != nil {
			return err
		}
	}

	return nil
}
