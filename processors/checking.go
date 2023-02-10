package processors

import (
	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) getChaseChecking() *model.StatementDescriptor {
	return &model.StatementDescriptor{
		List: []*model.LineDescriptor{
			model.NewLineWithField("Account", `Account\sNumber:\s*(?P<Account>[\d\s]+)`),
			model.NewLineWithFieldAndKey("OpeningBalance", "Amount", `^Beginning Balance\s+`+usdPattern),
			model.NewLineWithFieldAndKey("OpeningDate", "Date", `(?P<Date>\w+\s+\d+,\s+\d+)\s+through`),
			model.NewLineWithFieldAndKey("ClosingDate", "Date", `through\s+(?P<Date>\w+\s+\d+,\s+\d+)`),
			model.NewLineWithFieldAndKey("DepositsBalance", "Amount", `Deposits and Additions\s+`+usdPattern),
			model.NewLineWithField("DepositsStart", `(?P<DepositsStart>DEPOSITS AND ADDITIONS)`),
			model.NewLineWithField("DepositsEnd", `(?P<DepositsEnd>Total Deposits and Additions)`),
			model.NewLineWithFieldAndKey("WithdrawalsBalance", "Amount", `^Electronic Withdrawals\s+`+usdPattern),
			model.NewLineWithField("AtmDebitWithdrawalsStart", `(?P<AtmDebitWithdrawalsStart>ATM & DEBIT CARD WITHDRAWALS)`),
			model.NewLineWithField("AtmDebitWithdrawalsEnd", `^(?P<AtmDebitWithdrawalsEnd>Total ATM & Debit Card Withdrawals)`),
			model.NewLineWithField("WithdrawalsStart", `(?P<WithdrawalsStart>ELECTRONIC WITHDRAWALS)`),
			model.NewLineWithField("WithdrawalsEnd", `^(?P<WithdrawalsEnd>Total Electronic Withdrawals)`),
			model.NewLineWithFieldAndKey("ChecksBalance", "Amount", `^Checks Paid\s+`+usdPattern),
			model.NewLineWithField("ChecksStart", `(?P<ChecksStart>CHECKS PAID)`),
			model.NewLineWithField("ChecksEnd", `(?P<ChecksEnd>Total Checks Paid)`),
			model.NewLine(
				"Page",
				`Page (?P<Current>\d+) of (?P<Total>\d+)`,
				model.NewField("Current"), model.NewField("Total"),
			),
			model.NewLine(
				"Transaction",
				`^(?P<Date>\d{2}/\d{2})[_\.]*\s+(?P<Memo>.*?)\s{10,}`+usdPattern,
				model.NewField("Date"), model.NewField("Memo"), model.NewField("Amount"),
			),
			model.NewLine(
				"Transaction",
				`^(?P<Number>\d+)\s+(?P<Memo>.*?)\s+(?P<Date>\d{2}/\d{2})\s+`+usdPattern,
				model.NewField("Date"), model.NewField("Memo"), model.NewField("Amount"),
			),
		},
	}
}
