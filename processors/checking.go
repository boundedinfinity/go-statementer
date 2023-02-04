package processors

import "github.com/boundedinfinity/docsorter/model"

func lookup(path string, descriminator model.StatementDiscriminator) (*StatementProcessor, error) {
	config := model.CheckingConfig{
		Account:            model.MatcherConfig{Pattern: `Account\sNumber:\s*(?P<account>[\d\s]+)`},
		Transaction:        model.MatcherConfig{Pattern: `(?P<date>\d{2}/\d{2})\s+(?P<memo>.*?)\s+` + usdPattern},
		OpeningDate:        model.MatcherConfig{Pattern: `(?P<date>\w+\s+\d+,\s+\d+)\s+through`},
		OpeningBalance:     model.MatcherConfig{Pattern: `Beginning Balance\s+` + usdPattern},
		ClosingDate:        model.MatcherConfig{Pattern: `through\s+(?P<date>\w+\s+\d+,\s+\d+)`},
		ClosingBalance:     model.MatcherConfig{Pattern: `Ending Balance\s+` + usdPattern},
		DepositsBalance:    model.MatcherConfig{Pattern: `Deposits and Additions\s+` + usdPattern},
		DepositsStart:      model.MatcherConfig{Pattern: `(?P<depositsStart>DEPOSITS AND ADDITIONS)`},
		DepositsEnd:        model.MatcherConfig{Pattern: `(?P<depositsEnd>Total Deposits and Additions)`},
		ChecksBalance:      model.MatcherConfig{Pattern: `Checks Paid\s+` + usdPattern},
		ChecksStart:        model.MatcherConfig{Pattern: `(?P<checksStart>CHECKS PAID)`},
		ChecksEnd:          model.MatcherConfig{Pattern: `(?P<checksEnd>Total Checks Paid)`},
		WithdrawalsBalance: model.MatcherConfig{Pattern: `Electronic Withdrawals\s+` + usdPattern},
		WithdrawalsStart:   model.MatcherConfig{Pattern: `(?P<withdrawalsStart>ELECTRONIC WITHDRAWALS)`},
		WithdrawalsEnd:     model.MatcherConfig{Pattern: `(?P<withdrawalsEnd>Total Electronic Withdrawals)`},
	}

	processor, err := NewProcessor(path, config)

	if err != nil {
		return processor, err
	}

	return processor, nil
}
