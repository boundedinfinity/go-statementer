package model

import (
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
)

type FileSet struct {
	Source string   `yaml:"source"`
	Stem   string   `yaml:"stem"`
	Dir    string   `yaml:"dir"`
	Pdf    string   `yaml:"pdf"`
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

func FileSetFromSource(base, src string) FileSet {
	fs := NewFileSet()

	fs.Source = src
	fs.Pdf = pather.Base(fs.Source)
	fs.Stem = extentioner.Strip(fs.Pdf)
	fs.Dir = pather.Join(base, fs.Pdf)
	fs.Text = extentioner.Join(fs.Stem, ".txt")
	fs.Yaml = extentioner.Join(fs.Stem, ".yaml")
	fs.Csv = extentioner.Join(fs.Stem, ".csv")

	return fs
}

// func FileSetFromFileset(base, stem string, src FileSet) FileSet {
// 	fs := NewFileSet()

// 	fs.Source = src.Source
// 	fs.Pdf = strings.ReplaceAll(src.Pdf, src.Dir)
// 	fs.Stem = extentioner.Strip(fs.Pdf)
// 	fs.Dir = pather.Join(base, fs.Pdf)
// 	fs.Text = extentioner.Join(fs.Stem, ".txt")
// 	fs.Yaml = extentioner.Join(fs.Stem, ".yaml")
// 	fs.Csv = extentioner.Join(fs.Stem, ".csv")

// 	return fs
// }
