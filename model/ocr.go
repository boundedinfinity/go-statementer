package model

type ProcessStage struct {
	Source string   `yaml:"source"`
	Dir    string   `yaml:"dir"`
	Pdf    string   `yaml:"pdf"`
	Image  string   `yaml:"image"`
	Images []string `yaml:"images"`
	Text   string   `yaml:"text"`
	Texts  []string `yaml:"texts"`
	Csv    string   `yaml:"csv"`
	Yaml   string   `yaml:"yaml"`
}

func NewProcessStage() ProcessStage {
	return ProcessStage{
		Images: make([]string, 0),
		Texts:  make([]string, 0),
	}
}

type OcrContext struct {
	Stage1     ProcessStage `yaml:"stage1"`
	Stage2     ProcessStage `yaml:"stage2"`
	Dest       ProcessStage `yaml:"dest"`
	Data       []Extracted  `yaml:"extracted"`
	UserConfig UserConfigStatement
	Statement  *CheckingStatement `yaml:"statement"`
}

func NewOcrContext() *OcrContext {
	return &OcrContext{
		Stage1:    NewProcessStage(),
		Stage2:    NewProcessStage(),
		Dest:      NewProcessStage(),
		Data:      make([]Extracted, 0),
		Statement: NewCheckingStatement(),
	}
}

type Extracted struct {
	Name   string            `yaml:"name"`
	Values map[string]string `yaml:"values"`
}
