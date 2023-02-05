package model

type OcrContext struct {
	Source        string                 `yaml:"source"`
	Pdf           string                 `yaml:"pdf"`
	Images        []string               `yaml:"images"`
	Texts         []string               `yaml:"texts"`
	Text          string                 `yaml:"texts"`
	Discriminator StatementDiscriminator `yaml:"discriminator"`
	Extracted     []Extracted            `yaml:"extracted"`
	Statement     CheckingStatement      `yaml:"statement"`
}

type Extracted struct {
	Name   string            `yaml:"name"`
	Values map[string]string `yaml:"values"`
}

type StatementDiscriminator struct {
	Account string `yaml:"account"`
}

type StatementDiscriminatorConfig struct {
	Account string `yaml:"account"`
}
