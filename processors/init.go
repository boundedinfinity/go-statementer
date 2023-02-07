package processors

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/oriser/regroup"
)

func (t *ProcessManager) Init(descriptor *model.StatementDescriptor) error {
	for _, line := range descriptor.List {
		if err := util.ValidateLineRegex(*line); err != nil {
			return err
		}
	}

	for _, line := range descriptor.List {
		matcher, err := regroup.Compile(line.Pattern)

		if err != nil {
			return err
		}

		line.Regex = matcher
	}

	return nil
}
