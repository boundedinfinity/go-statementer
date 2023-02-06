package processors

import (
	"bufio"
	"os"

	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) GetClassifier(ocr *model.OcrContext) (model.Processor, error) {
	proecssor, err := newClassifier(t.logger, t.userConfig, ocr)

	if err != nil {
		return nil, err
	}

	return proecssor, nil
}

func (t *ProcessManager) Extract(ocr *model.OcrContext, processor model.Processor) error {
	file, err := os.Open(ocr.Text)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if err := processor.Extract(scanner.Text()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (t *ProcessManager) Convert(ocr *model.OcrContext, processor model.Processor) error {
	return nil
}
