package runtime

import (
	"errors"
	"os"
	"strings"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/boundedinfinity/go-commoner/slicer"
)

func (t *Runtime) LoadFiles() ([]model.OcrContext, error) {
	var ocrs []model.OcrContext
	allPaths := make([]string, 0)

	for _, path := range t.UserConfig.InputPaths {
		paths, err := util.GetFilteredFiles(path, t.extPdf)

		if err != nil {
			return ocrs, err
		}

		allPaths = append(allPaths, paths...)
	}

	allPaths = slicer.Dedup(allPaths)

	for _, path := range allPaths {
		ocr := model.NewOcrContext()
		ocr.Stage1.Source = path
		ocrs = append(ocrs, *ocr)
	}

	return ocrs, nil
}

func (t *Runtime) CalcFiles(dir, name string, dst *model.ProcessStage, src model.ProcessStage) error {
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

func (t *Runtime) Rename(ocr model.OcrContext, dst *model.ProcessStage, src model.ProcessStage) error {
	name := ocr.Statement.Account
	name = name[len(name)-4:]
	name += "-" + ocr.Statement.ClosingDate.String()

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

func (t *Runtime) prepareDirectory(stage *model.ProcessStage) error {
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
