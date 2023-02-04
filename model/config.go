package model

import (
	"github.com/boundedinfinity/rfc3339date"
)

type Config struct {
	InputPath   string            `yaml:"inputPath"`
	OutputPath  string            `yaml:"outputPath"`
	WorkPath    string            `yaml:"workPath"`
	SumExt      string            `yaml:"sumExt"`
	InputExt    string            `yaml:"inputExt"`
	Debug       bool              `yaml:"debug"`
	IgnorePaths []string          `yaml:"ignorePaths"`
	Statements  []StatementMapper `yaml:"statements"`
}

type StatementMapper struct {
	Name      string `yaml:"name"`
	Account   string `yaml:"account"`
	Processor string `yaml:"processor"`
}

type StatementDiscriminator struct {
	Account string `yaml:"account"`
}

type StatementDiscriminatorConfig struct {
	Account string `yaml:"account"`
}

type Transaction struct {
	Date   rfc3339date.Rfc3339Date `yaml:"date"`
	Memo   string                  `yaml:"memo"`
	Amount float32                 `yaml:"amount"`
}
