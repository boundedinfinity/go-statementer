package processors

import (
	"fmt"

	"github.com/boundedinfinity/docsorter/model"
)

func (t *ProcessManager) Transform(ocr *model.OcrContext) error {
	switch ocr.UserConfig.Processor {
	case "chase-checking":
		t.ocr.Checking = model.NewCheckingStatement()
		return t.TransformChecking(&t.ocr.Checking)
	case "chase-credit-card":
		t.ocr.CreditCard = model.NewCreditCardStatement()
		t.TransformCreditCard(&t.ocr.CreditCard)
	default:
		return fmt.Errorf("error transformer for %v", ocr.UserConfig.Account)
	}

	return nil
}

func (t *ProcessManager) TransformCreditCard(statement *model.CreditCardStatement) error {
	var section []model.Transaction

	for _, ext := range t.ocr.Extracted {
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
			if err := convertString(ext.Values, "Account", &statement.AccountNumber, accountCleanup...); err != nil {
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

func (t *ProcessManager) TransformChecking(statement *model.CheckingStatement) error {
	var section []model.Transaction

	for _, ext := range t.ocr.Extracted {
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
