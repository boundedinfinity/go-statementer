package model

type OcrContext struct {
	Source     string      `yaml:"source"`
	WorkDir    string      `yaml:"workDir"`
	WorkPdf    string      `yaml:"workPdf"`
	WorkImages []string    `yaml:"workImages"`
	WorkTexts  []string    `yaml:"workTexts"`
	WorkText   string      `yaml:"workText"`
	WorkCsv    string      `yaml:"workCsv"`
	WorkYaml   string      `yaml:"workYaml"`
	DestCsv    string      `yaml:"destCsv"`
	DestPdf    string      `yaml:"destPdf"`
	DestYaml   string      `yaml:"destYaml"`
	Data       []Extracted `yaml:"extracted"`
	UserConfig UserConfigStatement
	Statement  *CheckingStatement `yaml:"statement"`
}

func NewOcrContext() *OcrContext {
	return &OcrContext{
		WorkImages: make([]string, 0),
		WorkTexts:  make([]string, 0),
		Data:       make([]Extracted, 0),
		Statement:  NewCheckingStatement(),
	}
}

type Extracted struct {
	Name   string            `yaml:"name"`
	Values map[string]string `yaml:"values"`
}
