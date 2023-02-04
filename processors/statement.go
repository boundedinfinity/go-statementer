package processors

import (
	"errors"
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/oriser/regroup"
)

type StatementProcessor struct {
	Name      string
	section   string
	configs   []model.MatchConfig
	matchers  []*model.MatchContext
	raw       *model.CheckingStatementRaw
	statement *model.CheckingStatement
}

func NewProcessor(path string, config model.CheckingConfig) (*StatementProcessor, error) {
	processor := &StatementProcessor{
		raw: &model.CheckingStatementRaw{
			Source:   path,
			Deposits: make([]model.TransactionRaw, 0),
		},
	}

	addMatcher := func(pattern string, fn func(map[string]string)) error {
		matcher, err := regroup.Compile(pattern)

		if err != nil {
			return err
		}

		processor.matchers = append(processor.matchers, &model.MatchContext{
			Pattern: pattern,
			Extract: fn,
			Matcher: matcher,
		})

		return nil
	}

	extractTransaction := func(groups map[string]string) {
		transaction := model.TransactionRaw{}
		extractStringFn("date", &transaction.Date)(groups)
		extractStringFn("memo", &transaction.Memo)(groups)
		extractStringFn("usd", &transaction.Amount)(groups)

		switch processor.section {
		case "deposits":
			processor.raw.Deposits = append(processor.raw.Deposits, transaction)
		case "withdrawals":
			processor.raw.Withdrawals = append(processor.raw.Withdrawals, transaction)
		case "checks":
			processor.raw.Checks = append(processor.raw.Checks, transaction)
		}
	}

	sectionStart := func(name string) func(map[string]string) {
		return func(_ map[string]string) {
			processor.section = name
		}
	}

	addMatcher(config.Account.Pattern, extractStringFn("account", &processor.raw.Account))
	addMatcher(config.OpeningBalance.Pattern, extractStringFn("usd", &processor.raw.OpeningBalance))
	addMatcher(config.OpeningDate.Pattern, extractStringFn("date", &processor.raw.OpeningDate))
	addMatcher(config.ClosingBalance.Pattern, extractStringFn("usd", &processor.raw.ClosingBalance))
	addMatcher(config.ClosingDate.Pattern, extractStringFn("date", &processor.raw.ClosingDate))
	addMatcher(config.DepositsBalance.Pattern, extractStringFn("usd", &processor.raw.DepositsBalance))
	addMatcher(config.ChecksBalance.Pattern, extractStringFn("usd", &processor.raw.ChecksBalance))
	addMatcher(config.WithdrawalsBalance.Pattern, extractStringFn("usd", &processor.raw.WithdrawalsBalance))
	addMatcher(config.DepositsStart.Pattern, sectionStart("deposits"))
	addMatcher(config.ChecksStart.Pattern, sectionStart("checks"))
	addMatcher(config.WithdrawalsStart.Pattern, sectionStart("withdrawals"))
	addMatcher(config.Transaction.Pattern, extractTransaction)

	return processor, nil
}

func (p *StatementProcessor) Process(line string) error {
	for _, matcher := range p.matchers {
		fmt.Printf("%v: on %v\n", matcher.Pattern, line)
		groups, err := matcher.Matcher.Groups(line)

		if err != nil {
			if errors.Is(err, &regroup.NoMatchFoundError{}) {
				continue
			} else {
				return err
			}
		}

		if matcher.Extract != nil {
			matcher.Extract(groups)
		}
	}

	return nil
}
