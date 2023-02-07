package model

type OcrContext struct {
	Source     string      `yaml:"source"`
	WorkDir    string      `yaml:"workDir"`
	Pdf        string      `yaml:"pdf"`
	Images     []string    `yaml:"images"`
	Texts      []string    `yaml:"texts"`
	Text       string      `yaml:"texts"`
	Csv        string      `yaml:"csv"`
	Data       []Extracted `yaml:"extracted"`
	UserConfig UserConfigStatement
	Statement  *CheckingStatement `yaml:"statement"`
}

func NewOcrContext() *OcrContext {
	return &OcrContext{
		Images:    make([]string, 0),
		Texts:     make([]string, 0),
		Data:      make([]Extracted, 0),
		Statement: NewCheckingStatement(),
	}
}

type Extracted struct {
	Name   string            `yaml:"name"`
	Values map[string]string `yaml:"values"`
}
