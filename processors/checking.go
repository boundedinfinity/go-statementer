package processors

import (
	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) getChaseChecking() (model.Processor, error) {
	processor, err := newChaseChecking(t.logger, t.userConfig, t.ocr)

	return processor, err
}
