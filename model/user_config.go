package model

type UserConfig struct {
	InputPaths  []string              `yaml:"inputPaths"`
	OutputPath  string                `yaml:"outputPath"`
	WorkPath    string                `yaml:"workPath"`
	SumExt      string                `yaml:"sumExt"`
	InputExt    string                `yaml:"inputExt"`
	Log         string                `yaml:"log"`
	Debug       bool                  `yaml:"debug"`
	Reprocess   bool                  `yaml:"reprocess"`
	IgnorePaths []string              `yaml:"ignorePaths"`
	Statements  []UserConfigStatement `yaml:"statements"`
}

type UserConfigStatement struct {
	Name      string `yaml:"name"`
	Account   string `yaml:"account"`
	Processor string `yaml:"processor"`
}
