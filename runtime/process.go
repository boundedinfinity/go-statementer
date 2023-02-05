package runtime

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/processors"
	"github.com/boundedinfinity/docsorter/util"
)

func (t *Runtime) Process(ocr *model.OcrContext) error {
	if err := processors.Descriminator(ocr); err != nil {
		return err
	}

	if err := processors.ExtractStatement(t.logger, t.userConfig, ocr); err != nil {
		return err
	}

	for _, e := range ocr.Extracted {
		t.logger.Debugf(util.PrintLabeled(e.Name, fmt.Sprintf("%v", e.Values)))
	}

	return nil
}
