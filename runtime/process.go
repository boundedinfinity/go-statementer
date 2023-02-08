package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
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
