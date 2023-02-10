package processors

import "github.com/boundedinfinity/docsorter/model"

func (t *ProcessManager) Transform(descriptor *model.StatementDescriptor) error {
	var section []model.Transaction

	for _, ext := range t.ocr.Data {
		switch ext.Name {
		case "DepositsStart":
			section = make([]model.Transaction, 0)
		case "DepositsEnd":
			t.ocr.Statement.Deposits = section
		case "WithdrawalsStart":
			section = make([]model.Transaction, 0)
		case "WithdrawalsEnd":
			t.ocr.Statement.Withdrawals = section
		case "ChecksStart":
			section = make([]model.Transaction, 0)
		case "ChecksEnd":
			t.ocr.Statement.Checks = section
		case "AtmDebitWithdrawalsStart":
			section = make([]model.Transaction, 0)
		case "AtmDebitWithdrawalsEnd":
			t.ocr.Statement.AtmDebit = section
		case "Transaction":
			var transaction model.Transaction
			if err := converTransaction(ext.Values, &transaction, t.ocr.Statement.OpeningDate, t.ocr.Statement.ClosingDate); err != nil {
				return err
			}
			section = append(section, transaction)
		case "CheckTransaction":
			var transaction model.Transaction
			if err := converTransaction(ext.Values, &transaction, t.ocr.Statement.OpeningDate, t.ocr.Statement.ClosingDate); err != nil {
				return err
			}
			section = append(section, transaction)
		case "Account":
			if err := convertString(ext.Values, "Account", &t.ocr.Statement.Account, accountCleanup...); err != nil {
				return err
			}
		case "OpeningBalance":
			if err := convertFloat(ext.Values, "Amount", &t.ocr.Statement.OpeningBalance, usdCleanup...); err != nil {
				return err
			}
		case "ClosingBalance":
			if err := convertFloat(ext.Values, "Amount", &t.ocr.Statement.ClosingBalance, usdCleanup...); err != nil {
				return err
			}
		case "DepositsBalance":
			if err := convertFloat(ext.Values, "Amount", &t.ocr.Statement.DepositsBalance, usdCleanup...); err != nil {
				return err
			}
		case "ChecksBalance":
			if err := convertFloat(ext.Values, "Amount", &t.ocr.Statement.ChecksBalance, usdCleanup...); err != nil {
				return err
			}
		case "WithdrawalsBalance":
			if err := convertFloat(ext.Values, "Amount", &t.ocr.Statement.WithdrawalsBalance, usdCleanup...); err != nil {
				return err
			}
		case "OpeningDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat1, &t.ocr.Statement.OpeningDate, dateCleanup...); err != nil {
				return err
			}
		case "ClosingDate":
			if err := convertDate(ext.Values, "Date", chaseDateFormat1, &t.ocr.Statement.ClosingDate, dateCleanup...); err != nil {
				return err
			}
		}
	}

	return nil
}
