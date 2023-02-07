package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
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

	for _, path := range allPaths {
		ocr := model.NewOcrContext()
		ocr.Source = path
		ocrs = append(ocrs, *ocr)
	}

	return ocrs, nil
}
