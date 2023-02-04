package processors

import (
	"bufio"
	"os"

	"github.com/boundedinfinity/docsorter/model"
)

func Descriminator(path string) (model.StatementDiscriminator, error) {
	statement := model.StatementDiscriminator{
		Account: "",
	}

	// processor, err := NewProcessor(configs)

	// if err != nil {
	// 	return statement, err
	// }

	// if err := Process(path, processor); err != nil {
	// 	return statement, err
	// }

	return statement, nil
}

func ProcessStatement(path string, descriminator model.StatementDiscriminator) (model.CheckingStatementRaw, error) {
	zero := model.CheckingStatementRaw{}
	processor, err := lookup(path, descriminator)

	if err != nil {
		return zero, err
	}

	if err := Process(path, processor); err != nil {
		return zero, err
	}

	return *processor.raw, nil
}

func Process(path string, processor *StatementProcessor) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if err := processor.Process(scanner.Text()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
