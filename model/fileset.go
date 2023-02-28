package model

type FileSet struct {
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

func NewFileSet() FileSet {
	return FileSet{
		Images: make([]string, 0),
		Texts:  make([]string, 0),
	}
}
