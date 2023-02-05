package processors

import (
	"bufio"
	"os"

	"github.com/boundedinfinity/docsorter/model"
)

func Descriminator(ocr *model.OcrContext) error {
	ocr.Discriminator = model.StatementDiscriminator{
		Account: "",
	}

	return nil
}

func ExtractStatement(userConfig model.UserConfig, ocr *model.OcrContext) error {
	processor, err := lookup(userConfig, ocr)

	if err != nil {
		return err
	}

	if err := Extract(ocr, processor); err != nil {
		return err
	}

	return nil
}

func Extract(ocr *model.OcrContext, processor *StatementProcessor) error {
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
