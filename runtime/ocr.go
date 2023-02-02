package runtime

import (
	"fmt"

	"github.com/boundedinfinity/go-commoner/environmenter"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/boundedinfinity/go-commoner/stringer"
)

func (t *Runtime) OcrSingle(path string) (string, error) {
	var workTxt string
	imgExt := ".png"

	pdfPath := environmenter.Sub(path)
	pdfFullname := pather.Base(pdfPath)
	noExt := extentioner.Strip(pdfFullname)
	pdfExt := extentioner.Ext(pdfFullname)
	extFullname := extentioner.Swap(pdfFullname, pdfExt, ".txt")
	workDir := pather.Join(t.config.WorkPath, noExt)
	workPdf := pather.Join(workDir, pdfFullname)
	workTxt = pather.Join(workDir, extFullname)

	if err := prepOutDir(workDir); err != nil {
		return workTxt, err
	}

	if err := copyFile(workPdf, pdfPath); err != nil {
		return workTxt, err
	}

	env := map[string]string{
		"WORK_DIR": workDir,
	}

	// -density 300 -quiet $pdf_fullname $pdf_nameonly-%04d.$IMAGE_EXT
	imageMagickArgs := []string{
		"imagemagick", "-quiet", "-density", "300",
		pdfFullname, fmt.Sprintf("%v-%%04d%v", noExt, imgExt),
	}

	if stdOut, err := t.runDocker(env, imageMagickArgs); err != nil {
		return workTxt, err
	} else {
		fmt.Println(stdOut)
	}

	allFiles, err := pather.GetFiles(workDir)

	if err != nil {
		return workTxt, err
	}

	filteredFiles := slicer.Filter(allFiles, func(p string) bool {
		return stringer.EndsWith(p, imgExt)
	})

	filteredFiles = slicer.Map(filteredFiles, func(f string) string {
		return pather.Base(f)
	})

	for _, file := range filteredFiles {
		//  --oem 1 -l eng --psm 6 -c preserve_interword_spaces=1 $image_fullname $image_onlyname
		tesseractArgs := []string{
			"tesseract",
			"--oem", "1", "-l", "eng", "--psm", "6", "-c", "preserve_interword_spaces=1",
			file, extentioner.Strip(file),
		}

		if stdOut, err := t.runDocker(env, tesseractArgs); err != nil {
			return workTxt, err
		} else {
			fmt.Println(stdOut)
		}
	}

	allFiles, err = pather.GetFiles(workDir)

	if err != nil {
		return workTxt, err
	}

	outputFiles := slicer.Filter(allFiles, func(p string) bool {
		return stringer.EndsWith(p, ".txt")
	})

	if err := appendFile(workTxt, outputFiles...); err != nil {
		return workTxt, err
	}

	return workTxt, nil
}
