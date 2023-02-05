package processors

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/sirupsen/logrus"
)

func lookup(logger *logrus.Logger, userConfig model.UserConfig, ocr *model.OcrContext) (*StatementProcessor, error) {
	// txDescriptor := model.NewLine(
	// 	"Transaction",
	// 	`(?P<Date>\d{2}/\d{2})\s+(?P<Memo>.*?)\s+`+usdPattern,
	// 	model.NewField("Date"), model.NewField("Memo"), model.NewField("Amount"),
	// )

	descriptor := model.CheckingDescriptor{
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

	processor, err := NewProcessor(logger, userConfig, ocr, descriptor)

	if err != nil {
		return processor, err
	}

	return processor, nil
}
