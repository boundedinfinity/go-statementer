package processors

import (
	"errors"
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/oriser/regroup"
)

type StatementProcessor struct {
	Name       string
	userConfig model.UserConfig
	ocr        *model.OcrContext
	desc       model.CheckingDescriptor
}

func NewProcessor(userConfig model.UserConfig, ocr *model.OcrContext, desc model.CheckingDescriptor) (*StatementProcessor, error) {
	processor := &StatementProcessor{
		ocr:        ocr,
		userConfig: userConfig,
	}

	for _, line := range desc.Lines() {
		if err := validateLineRegex(*line); err != nil {
			return processor, err
		}
	}

	for _, line := range desc.Lines() {
		matcher, err := regroup.Compile(line.Pattern)

		if err != nil {
			return processor, err
		}

		line.Regex = matcher
	}

	// for _, section := range desc.Sections() {
	// 	if err := createSectionRegex(section); err != nil {
	// 		return processor, err
	// 	}
	// }

	processor.desc = desc

	return processor, nil
}

func (p *StatementProcessor) Extract(line string) error {
	for _, lineDesc := range p.desc.Lines() {
		if p.userConfig.Debug {
			fmt.Printf("%v: on %v\n", lineDesc.Pattern, line)
		}

		groups, err := lineDesc.Regex.Groups(line)

		if err != nil {
			if errors.Is(err, &regroup.NoMatchFoundError{}) {
				continue
			} else {
				return err
			}
		}

		p.ocr.Extracted = append(p.ocr.Extracted, model.Extracted{
			Name:   lineDesc.Name,
			Values: groups,
		})
	}

	return nil
}
