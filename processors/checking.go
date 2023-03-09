package processors

import (
	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) getChaseChecking() *model.StatementDescriptor {
	return &model.StatementDescriptor{
		List: []*model.LineDescriptor{
			model.NewLineWithField("Account", `Account\sNumber:\s*(?P<Account>[\d\s]+?)\s{5,}`),
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

func (t *ProcessManager) transformChecking(statement *model.CheckingStatement) error {
	var section []model.Transaction

	for _, ext := range t.pc.Extracted {
		switch ext.Name {
		case "DepositsStart":
			section = make([]model.Transaction, 0)
		case "DepositsEnd":
			statement.Deposits = section
		case "WithdrawalsStart":
			section = make([]model.Transaction, 0)
		case "WithdrawalsEnd":
			statement.Withdrawals = section
		case "ChecksStart":
			section = make([]model.Transaction, 0)
		case "ChecksEnd":
			statement.Checks = section
		case "AtmDebitWithdrawalsStart":
			section = make([]model.Transaction, 0)
		case "AtmDebitWithdrawalsEnd":
			statement.AtmDebit = section
		case "Transaction":
			var transaction model.Transaction

			if err := convertTransaction(ext.Values, &transaction, statement.OpeningDate, statement.ClosingDate); err != nil {
				return err
			}

			section = append(section, transaction)
		case "CheckTransaction":
			var transaction model.Transaction

			if err := convertTransaction(ext.Values, &transaction, statement.OpeningDate, statement.ClosingDate); err != nil {
				return err
			}

			section = append(section, transaction)
		case "Account":
			if err := convertString(ext.Values, "Account", &statement.Account, accountCleanup...); err != nil {
				return err
			}
		case "OpeningBalance":
			if err := convertFloat(ext.Values, "Amount", &statement.OpeningBalance, usdCleanup...); err != nil {
				return err
			}
		case "ClosingBalance":
			if err := convertFloat(ext.Values, "Amount", &statement.ClosingBalance, usdCleanup...); err != nil {
				return err
			}
		case "DepositsBalance":
			if err := convertFloat(ext.Values, "Amount", &statement.DepositsBalance, usdCleanup...); err != nil {
				return err
			}
		case "ChecksBalance":
			if err := convertFloat(ext.Values, "Amount", &statement.ChecksBalance, usdCleanup...); err != nil {
				return err
			}
		case "WithdrawalsBalance":
			if err := convertFloat(ext.Values, "Amount", &statement.WithdrawalsBalance, usdCleanup...); err != nil {
				return err
			}
		case "OpeningDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat1, &statement.OpeningDate, dateCleanup...); err != nil {
				return err
			}
		case "ClosingDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat1, &statement.ClosingDate, dateCleanup...); err != nil {
				return err
			}
		}
	}

	return nil
}
