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

	if err := processors.ExtractStatement(t.userConfig, ocr); err != nil {
		return err
	}

	if t.userConfig.Debug {
		for _, e := range ocr.Extracted {
			util.PrintLabeled(e.Name, fmt.Sprintf("%v", e.Values))
		}
	}

	return nil
}
