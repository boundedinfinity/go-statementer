package model

type OcrContext struct {
	Source     string      `yaml:"source"`
	WorkDir    string      `yaml:"workDir"`
	Pdf        string      `yaml:"pdf"`
	Images     []string    `yaml:"images"`
	Texts      []string    `yaml:"texts"`
	Text       string      `yaml:"texts"`
	Data       []Extracted `yaml:"extracted"`
	UserConfig UserConfigStatement
	Statement  CheckingStatement `yaml:"statement"`
}

type Extracted struct {
	Name   string            `yaml:"name"`
	Values map[string]string `yaml:"values"`
}
