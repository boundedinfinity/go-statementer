package processors

import (
	"bufio"
	"errors"
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/oriser/regroup"
)

func (t *ProcessManager) Extract(statement *model.StatementDescriptor) error {
	file, err := os.Open(t.ocr.Stage1.Text)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text := scanner.Text()
		for _, line := range statement.List {
			t.logger.Tracef("[[[[%v]]]][[[[%v]]]]\n", line.Pattern, text)

			groups, err := line.Regex.Groups(text)

			if err != nil {
				if errors.Is(err, &regroup.NoMatchFoundError{}) {
					continue
				} else {
					return err
				}
			}

			t.ocr.Extracted = append(t.ocr.Extracted, model.Extracted{
				Name:   line.Name,
				Values: groups,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
