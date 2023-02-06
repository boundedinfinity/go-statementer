package processors

import (
	"errors"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/oriser/regroup"
	"github.com/sirupsen/logrus"
)

type ChaseCheckingStatementProcessor struct {
	name                   string
	userConfig             model.UserConfig
	Statement              model.CheckingStatement
	logger                 *logrus.Logger
	ocr                    *model.OcrContext
	Account                *model.LineDescriptor
	OpeningBalance         *model.LineDescriptor
	ClosingBalance         *model.LineDescriptor
	OpeningDate            *model.LineDescriptor
	ClosingDate            *model.LineDescriptor
	DepositsBalance        *model.LineDescriptor
	DepositsStart          *model.LineDescriptor
	DepositsEnd            *model.LineDescriptor
	ChecksBalance          *model.LineDescriptor
	ChecksStart            *model.LineDescriptor
	ChecksEnd              *model.LineDescriptor
	WithdrawalsBalance     *model.LineDescriptor
	WithdrawalsStart       *model.LineDescriptor
	WithdrawalsEnd         *model.LineDescriptor
	Transaction            *model.LineDescriptor
	Page                   *model.LineDescriptor
	AnnualPercentageEarned *model.LineDescriptor
	InterestEarned         *model.LineDescriptor
	InterestPaid           *model.LineDescriptor
}

var _ model.Processor = &ChaseCheckingStatementProcessor{}

func newChaseChecking(logger *logrus.Logger, userConfig model.UserConfig, ocr *model.OcrContext) (*ChaseCheckingStatementProcessor, error) {
	processor := &ChaseCheckingStatementProcessor{
		ocr:                ocr,
		userConfig:         userConfig,
		logger:             logger,
		Statement:          *model.NewCheckingStatement(),
		Account:            model.NewLineWithField("Account", `Account\sNumber:\s*(?P<Account>[\d\s]+)`),
		OpeningBalance:     model.NewLineWithFieldAndKey("OpeningBalance", "Amount", `^Beginning Balance\s+`+usdPattern),
		OpeningDate:        model.NewLineWithFieldAndKey("OpeningDate", "Date", `(?P<Date>\w+\s+\d+,\s+\d+)\s+through`),
		ClosingBalance:     model.NewLineWithFieldAndKey("ClosingBalance", "Amount", `^Ending Balance\s+`+usdPattern),
		ClosingDate:        model.NewLineWithFieldAndKey("ClosingDate", "Date", `through\s+(?P<Date>\w+\s+\d+,\s+\d+)`),
		DepositsBalance:    model.NewLineWithFieldAndKey("DepositsBalance", "Amount", `Deposits and Additions\s+`+usdPattern),
		DepositsStart:      model.NewLineWithField("DepositsStart", `(?P<DepositsStart>DEPOSITS AND ADDITIONS)`),
		DepositsEnd:        model.NewLineWithField("DepositsEnd", `(?P<DepositsEnd>Total Deposits and Additions)`),
		WithdrawalsBalance: model.NewLineWithFieldAndKey("WithdrawalsBalance", "Amount", `^Electronic Withdrawals\s+`+usdPattern),
		WithdrawalsStart:   model.NewLineWithField("WithdrawalsStart", `(?P<WithdrawalsStart>ELECTRONIC WITHDRAWALS)`),
		WithdrawalsEnd:     model.NewLineWithField("WithdrawalsEnd", `(?P<WithdrawalsEnd>Total Electronic Withdrawals)`),
		ChecksBalance:      model.NewLineWithFieldAndKey("ChecksBalance", "Amount", `^Checks Paid\s+`+usdPattern),
		ChecksStart:        model.NewLineWithField("ChecksStart", `(?P<ChecksStart>CHECKS PAID)`),
		ChecksEnd:          model.NewLineWithField("ChecksEnd", `(?P<ChecksEnd>Total Checks Paid)`),
		Page: model.NewLine(
			"Page",
			`Page (?P<Current>\d+) of (?P<Total>\d+)`,
			model.NewField("Current"), model.NewField("Total"),
		),
		Transaction: model.NewLine(
			"Transaction",
			`(?P<Date>\d{2}/\d{2})\s+(?P<Memo>.*?)\s+`+usdPattern,
			model.NewField("Date"), model.NewField("Memo"), model.NewField("Amount"),
		),
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

func (t ChaseCheckingStatementProcessor) Name() string {
	return t.name
}

func (t *ChaseCheckingStatementProcessor) Convert() error {
	var section []model.Transaction

	for _, ext := range t.ocr.Data {
		switch ext.Name {
		case "DepositsStart":
			section = make([]model.Transaction, 0)
		case "DepositsEnd":
			t.Statement.Deposits = section
		case "WithdrawalsStart":
			section = make([]model.Transaction, 0)
		case "WithdrawalsEnd":
			t.Statement.Withdrawals = section
		case "ChecksStart":
			section = make([]model.Transaction, 0)
		case "ChecksEnd":
			t.Statement.Checks = section
		case "Transaction":
			var transaction model.Transaction
			if err := converTransaction(ext.Values, &transaction, t.Statement.OpeningDate, t.Statement.ClosingDate); err != nil {
				return err
			}
			section = append(section, transaction)
		case "Account":
			if err := convertString(ext.Values, "Account", &t.Statement.Account, accountCleanup...); err != nil {
				return err
			}
		case "OpeningBalance":
			if err := convertFloat(ext.Values, "Amount", &t.Statement.OpeningBalance, usdCleanup...); err != nil {
				return err
			}
		case "ClosingBalance":
			if err := convertFloat(ext.Values, "Amount", &t.Statement.ClosingBalance, usdCleanup...); err != nil {
				return err
			}
		case "DepositsBalance":
			if err := convertFloat(ext.Values, "Amount", &t.Statement.DepositsBalance, usdCleanup...); err != nil {
				return err
			}
		case "ChecksBalance":
			if err := convertFloat(ext.Values, "Amount", &t.Statement.ChecksBalance, usdCleanup...); err != nil {
				return err
			}
		case "WithdrawalsBalance":
			if err := convertFloat(ext.Values, "Amount", &t.Statement.WithdrawalsBalance, usdCleanup...); err != nil {
				return err
			}
		case "OpeningDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat1, &t.Statement.OpeningDate, dateCleanup...); err != nil {
				return err
			}
		case "ClosingDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat1, &t.Statement.ClosingDate, dateCleanup...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t ChaseCheckingStatementProcessor) Print() {
	t.logger.Info(util.PrintLabeled("Name", t.ocr.UserConfig.Name))
	t.logger.Info(util.PrintLabeled("Account", t.Statement.Account))
	t.logger.Info(util.PrintLabeled("OpeningBalance", t.Statement.OpeningBalance))
	t.logger.Info(util.PrintLabeled("ClosingBalance", t.Statement.ClosingBalance))
	t.logger.Info(util.PrintLabeled("OpeningDate", t.Statement.OpeningDate))
	t.logger.Info(util.PrintLabeled("ClosingDate", t.Statement.ClosingDate))
}

func (t ChaseCheckingStatementProcessor) Lines() []*model.LineDescriptor {
	lines := make([]*model.LineDescriptor, 0)
	lines = append(lines, t.Account, t.Transaction, t.Page,
		t.OpeningBalance, t.OpeningDate,
		t.ClosingBalance, t.ClosingDate,
		t.WithdrawalsBalance, t.WithdrawalsStart,
		t.ChecksBalance, t.ChecksStart, t.ChecksEnd,
		t.DepositsBalance, t.DepositsStart, t.DepositsEnd,
		t.AnnualPercentageEarned, t.InterestEarned, t.InterestPaid,
	)

	return slicer.Filter(lines, func(line *model.LineDescriptor) bool { return line != nil })
}

func (p *ChaseCheckingStatementProcessor) Extract(line string) error {
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
