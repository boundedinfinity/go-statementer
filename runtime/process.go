package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
)

func (t *Runtime) Process(pc *model.ProcessContext) error {
	manager := processors.NewManager(t.logger, t.UserConfig, pc)
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

	if err := manager.Transform(pc); err != nil {
		return err
	}

	return nil
}
