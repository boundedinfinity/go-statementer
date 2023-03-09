package processors

import "github.com/boundedinfinity/docsorter/model"

func (t *ProcessManager) getCreditCard() *model.StatementDescriptor {
	return &model.StatementDescriptor{
		List: []*model.LineDescriptor{
			model.NewLineWithField("Account", `Account\sNumber:\s*(?P<Account>[\d\s]+?)\s{5,}`),
			model.NewLineWithField("OpeningBalance", `^(?P<OpeningBalance>Previous Balance)\s+`+usdPattern),
			model.NewLineWithField("ClosingBalance", `^(?P<ClosingBalance>New Balance)\s+`+usdPattern),
			model.NewLineWithFieldAndKey("OpeningDate", "Date", `(?P<Date>\d+/\d+/\d+) - \d+/\d+/\d+`),
			model.NewLineWithFieldAndKey("ClosingDate", "Date", `\d+/\d+/\d+ - (?P<Date>\d+/\d+/\d+)`),
			model.NewLineWithFieldAndKey("Payments", "Amount", `Payment, Credits\s+`+usdPattern),
			model.NewLineWithFieldAndKey("Purchases", "Amount", `Purchases\s+`+usdPattern),
			model.NewLineWithFieldAndKey("CashAdvances", "Amount", `Cash Advances\s+`+usdPattern),
			model.NewLineWithFieldAndKey("BalanceTransfers", "Amount", `Balance Transfers\s+`+usdPattern),
			model.NewLineWithFieldAndKey("FeesCharged", "Amount", `Fees Charged\s+`+usdPattern),
			model.NewLineWithFieldAndKey("InterestCharged", "Amount", `Interest Charged\s+`+usdPattern),
			model.NewLineWithField("PaymentsStart", `^(?P<PaymentsStart>PAYMENTS AND OTHER CREDITS)\s*$`),
			model.NewLineWithField("PurchasesStart", `^(?P<PurchasesStart>^PURCHASES?)\s*$`),
			model.NewLineWithField("PurchasesAndRedemptionsStart", `^(?P<PurchasesAndRedemptionsStart>PURCHASES AND REDEMPTIONS)\s*$`),
			model.NewLineWithField("ImportantNews", `(?P<ImportantNews>IMPORTANT NEWS)\s*`),
			model.NewLine(
				"Page",
				`Page (?P<Current>\d+) of (?P<Total>\d+)`,
				model.NewField("Current"), model.NewField("Total"),
			),
			model.NewLine(
				"Transaction",
				`(?P<Date>\d{2}/\d{2})[_\.]*\s+(?P<Memo>.*?)\s{10,}`+usdPattern,
				model.NewField("Date"), model.NewField("Memo"), model.NewField("Amount"),
			),
		},
	}
}

func (t *ProcessManager) transformCreditCard(statement *model.CreditCardStatement) error {
	var section []model.Transaction

	for _, ext := range t.pc.Extracted {
		switch ext.Name {
		case "PaymentsStart":
			section = make([]model.Transaction, 0)
		case "PurchasesStart":
			if len(statement.Payments) <= 0 {
				statement.Payments = section
				section = make([]model.Transaction, 0)
			}
		case "PurchasesAndRedemptionsStart":
			statement.Purchases = section
			section = make([]model.Transaction, 0)
		case "ImportantNews":
			statement.Redemptions = section
			section = make([]model.Transaction, 0)
		case "Transaction":
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
		case "OpeningDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat3, &statement.OpeningDate, dateCleanup...); err != nil {
				return err
			}
		case "ClosingDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat3, &statement.ClosingDate, dateCleanup...); err != nil {
				return err
			}
		}
	}

	return nil
}
