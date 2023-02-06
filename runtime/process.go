package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
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
