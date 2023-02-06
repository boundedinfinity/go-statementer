package processors

import (
	"errors"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/oriser/regroup"
	"github.com/sirupsen/logrus"
)

type ClassifierProcessor struct {
	name       string
	userConfig model.UserConfig
	classifier model.AccountClassifier
	logger     *logrus.Logger
	ocr        *model.OcrContext
	accounts   []*model.LineDescriptor
}

var _ model.Processor = &ClassifierProcessor{}

func newClassifier(logger *logrus.Logger, userConfig model.UserConfig, ocr *model.OcrContext) (*ClassifierProcessor, error) {
	processor := &ClassifierProcessor{
		ocr:        ocr,
		userConfig: userConfig,
		logger:     logger,
		accounts: []*model.LineDescriptor{
			model.NewLineWithField("Account", `Account\sNumber:\s*(?P<Account>[\d\s]+)`),
		},
	}

	for _, line := range processor.Lines() {
		if err := util.ValidateLineRegex(*line); err != nil {
			return processor, err
		}
	}

	for _, line := range processor.Lines() {
		matcher, err := regroup.Compile(line.Pattern)

		if err != nil {
			return processor, err
		}

		line.Regex = matcher
	}

	return processor, nil
}

func (t ClassifierProcessor) Convert() error {
	return nil
}

func (t ClassifierProcessor) Name() string {
	return t.name
}

func (t ClassifierProcessor) Lines() []*model.LineDescriptor {
	return t.accounts
}

func (t ClassifierProcessor) Print() {
}

func (p *ClassifierProcessor) Extract(line string) error {
	for _, lineDesc := range p.Lines() {
		p.logger.Tracef("%v: on %v\n", lineDesc.Pattern, line)

		groups, err := lineDesc.Regex.Groups(line)

		if err != nil {
			if errors.Is(err, &regroup.NoMatchFoundError{}) {
				continue
			} else {
				return err
			}
		}

		p.ocr.Data = append(p.ocr.Data, model.Extracted{
			Name:   lineDesc.Name,
			Values: groups,
		})
	}

	return nil
}

// # -   name: chase-credit-card
//     #     dateFormat: 01/02/06
//     #     dateKey: StartDate
//     #     patterns:
//     #         -   label: AccountNumber
//     #             pattern: account number\:\s(\d+\s\d+\s\d+\s\d+\s)
//     #         -   label: StartDate
//     #             pattern: (\d+/\d+/\d+)\s-\s\d+/\d+/\d+
//     #         -   label: EndDate
//     #             pattern: \d+/\d+/\d+\s-\s(\d+/\d+/\d+)
//     #         -   label: Other1
//     #             pattern: wilmington, de 19850-5298

//     # -   name: donna-farrow-ownwer-statement
//     #     dateFormat: Jan 2006
//     #     dateKey: StartDate
//     #     patterns:
//     #         -   label: AccountNumber
//     #             pattern: donna farrow & company
//     #         -   label: StartDate
//     #             pattern: 'months: (\w+\s+\d+)\s+-\s+\w+\s+\d+'
//     #         -   label: EndDate
//     #             pattern: \w+\s+\d+\s+-\s+(\w+\s+\d+)
//     #         -   label: Other1
//     #             pattern: owner\s+statement

//     # -   name: cornerstone-owner-statement
//     #     dateFormat: Jan 02, 2006
//     #     dateKey: StartDate
//     #     patterns:
//     #         -   label: AccountNumber
//     #             pattern: cornerstone
//     #         -   label: StartDate
//     #             pattern: (\w+\s+\d+,\s+\d+)\s+-\s+\w+\s+\d+,\s+\d+
//     #         -   label: EndDate
//     #             pattern: \w+\s+\d+,\s+\d+\s+-\s+(\w+\s+\d+,\s+\d+)
//     #         -   label: Other1
//     #             pattern: owner\s+statement
//     #         -   label: Other2
//     #             pattern: www.rentfd.com

//     # -   name: ally
//     #     dateFormat: 01/02/2006
//     #     dateKey: StatementDate
//     #     patterns:
//     #         -   label: AccountNumber
//     #             pattern: ally bank
//     #         -   label: StatementDate
//     #             pattern: statement date\n(\d+/\d+/\d+)
//     #         -   label: Other1
//     #             pattern: ally bank

//     # -   name: wells-fargo-credit-card
//     #     dateFormat: 01/02/2006
//     #     dateKey: EndDate
//     #     patterns:
//     #         -   label: AccountNumber
//     #             pattern: ending in\s+(\d+)
//     #         -   label: EndDate
//     #             pattern: to\s+(\d+/\d+/\d+)
//     #         -   label: Other1
//     #             pattern: wells fargo cash wise visa signature card

//     # -   name: novo-checking
//     #     dateFormat: Jan 02, 2006
//     #     dateKey: StartDate
//     #     patterns:
//     #         -   label: EndDate
//     #             pattern: thru\s+(\w+\s+\d+,\s+\d+)
//     #         -   label: StartDate
//     #             pattern: (\w+\s+\d+,\s+\d+)\s+thru
//     #         -   label: Type
//     #             pattern: middlesex federal savings
//     #         -   label: AccountNumber
//     #             pattern: account number\nxxxx\s+(\d+)

//     # -   name: chase-checking
//     #     dateFormat: January 2, 2006
//     #     dateKey: StartDate
//     #     patterns:
//     #         -   label: EndDate
//     #             pattern: through\s+(\w+\s+\d+,\s+\d+)
//     #         -   label: StartDate
//     #             pattern: (\w+\s+\d+,\s+\d+)\s+through
//     #         -   label: AccountNumber
//     #             pattern: account number:\s+(\d+)
//     #         -   label: Type
//     #             pattern: checking summary

//     # -   name: chase-savings
//     #     dateFormat: January 2, 2006
//     #     dateKey: StartDate
//     #     patterns:
//     #         -   label: EndDate
//     #             pattern: through\s+(\w+\s+\d+,\s+\d+)
//     #         -   label: StartDate
//     #             pattern: (\w+\s+\d+,\s+\d+)\s+through
//     #         -   label: AccountNumber
//     #             pattern: account number:\s+(\d+)
//     #         -   label: Type
//     #             pattern: savings summary

//     # -   name: wells-fargo-mortgage
//     #     dateFormat: 01/02/06
//     #     dateKey: StatementDate
//     #     patterns:
//     #         -   label: StatementDate
//     #             pattern: statement date(?:\n.*){4}(\d{2}/\d{2}/\d{2})
//     #         -   label: AccountNumber
//     #             pattern: loan number(?:\n.*){3}\n(\d+)
//     #         -   label: Type
//     #             pattern: wells fargo home mortgage
