package runtime

import (
	"os"
	"path"
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

	for _, path := range t.userConfig.InputPaths {
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

func Stage2Stage(dir, name string, dst *model.ProcessStage, src model.ProcessStage) {
	dst.Source = src.Source
	dst.Dir = path.Join(dir, name)
	srcRep := path.Join(src.Dir, pather.Base(src.Dir))
	dstRep := path.Join(dst.Dir, name)

	replace := func(s string) string {
		return strings.ReplaceAll(s, srcRep, dstRep)
	}

	dst.Pdf = replace(src.Pdf)
	dst.Image = replace(src.Image)
	dst.Images = slicer.Map(src.Images, replace)
	dst.Text = replace(src.Text)
	dst.Texts = slicer.Map(src.Texts, replace)
	dst.Csv = replace(src.Csv)
	dst.Yaml = replace(src.Yaml)
}

func (t *Runtime) Rename(ocr2 model.OcrContext, dst *model.ProcessStage, src model.ProcessStage) error {
	name := ocr2.Statement.Account
	name = name[len(name)-4:]
	name += "-" + ocr2.Statement.ClosingDate.String()

	Stage2Stage(t.userConfig.WorkPath, name, dst, src)

	if err := pather.DirEnsure(dst.Dir); err != nil {
		return err
	}

	files, err := pather.GetFiles(src.Dir)

	if err != nil {
		return err
	}

	for _, old := range files {
		new := strings.ReplaceAll(old, pather.Base(src.Dir), name)

		if err := os.Rename(old, new); err != nil {
			return err
		}
	}

	if err := util.EnsureDelete(src.Dir); err != nil {
		return nil
	}

	return nil
}

func (t *Runtime) prepareDirectory(stage *model.ProcessStage) error {
	sourceName := extentioner.Strip(pather.Base(stage.Source))
	stage.Dir = pather.Join(t.userConfig.WorkPath, sourceName)
	stage.Pdf = pather.Join(stage.Dir, pather.Base(stage.Source))

	if t.userConfig.Reprocess {
		if err := util.EnsureDelete(stage.Dir); err != nil {
			return err
		}
	}

	if err := pather.DirEnsure(stage.Dir); err != nil {
		return err
	}

	if t.userConfig.Reprocess || !pather.PathExists(stage.Pdf) {
		if err := util.CopyFile(stage.Pdf, stage.Source); err != nil {
			return err
		}
	}

	return nil
}
