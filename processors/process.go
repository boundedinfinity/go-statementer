package processors

import (
	"bufio"
	"os"

	"github.com/boundedinfinity/docsorter/model"
)

func GetAccount(path string) (string, error) {
	var account string

	processor := BuildLineProcessor("Account Number").
		Match("accunt", accountPattern).
		Contains("account").
		Clean("account", accountCleanup...).
		ExtractString("account", &account).
		Done()

	if err := Process(path, processor); err != nil {
		return account, err
	}

	return account, nil
}

func ProcessStatement(path string, account string) (model.CheckingStatement, error) {
	var statement model.CheckingStatement

	accountP := BuildLineProcessor("Account Number").
		Match("accunt", accountPattern).
		Contains("account").
		Clean("account", trimLeading0, removeSpaces).
		ExtractString("account", &statement.AccountNumber).
		Done()

	openingBalanceP := BuildLineProcessor("Opening Balance").
		Match("usd", openingBalancePatterns...).
		Contains("usd").
		Clean("usd", usdCleanup...).
		ExtractFloat("usd", &statement.OpeningBalance).
		Done()

	closingBalanceP := BuildLineProcessor("Closing Balance").
		Match("usd", closingBalancePatterns...).
		Contains("usd").
		Clean("usd", usdCleanup...).
		ExtractFloat("usd", &statement.ClosingBalance).
		Done()

	processor := &StatementProcessor{
		Name: "Checking Account",
		processors: []Processor{
			accountP,
			openingBalanceP,
			closingBalanceP,
		},
	}

	if err := Process(path, processor); err != nil {
		return statement, err
	}

	return statement, nil
}

func Process(path string, processor Processor) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() && !processor.Completed() {
		if err := processor.Process(scanner.Text()); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
