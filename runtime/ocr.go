package runtime

import (
	"fmt"
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/environmenter"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/boundedinfinity/go-commoner/slicer"
	"github.com/boundedinfinity/go-commoner/stringer"
)

func (t Runtime) getImageFiles(workDir, ext string) ([]string, error) {
	imageFiles, err := pather.GetFiles(workDir)

	if err != nil {
		return imageFiles, err
	}

	imageFiles = slicer.Filter(imageFiles, func(p string) bool {
		return stringer.EndsWith(p, ext)
	})

	return imageFiles, nil
}

func (t *Runtime) OcrSingle(ocr *model.OcrContext) error {
	ocr.Source = environmenter.Sub(ocr.Source)
	pdfFullname := pather.Base(ocr.Source)
	noExt := extentioner.Strip(pdfFullname)
	pdfExt := extentioner.Ext(pdfFullname)
	txtFullname := extentioner.Swap(pdfFullname, pdfExt, ".txt")
	workDir := pather.Join(t.userConfig.WorkPath, noExt)
	ocr.Pdf = pather.Join(workDir, pdfFullname)
	ocr.Text = pather.Join(workDir, txtFullname)

	if t.userConfig.Reprocess {
		if err := util.EnsureDelete(workDir); err != nil {
			return err
		}
	}

	if err := pather.DirEnsure(workDir); err != nil {
		return err
	}

	if t.userConfig.Reprocess || !pather.PathExists(ocr.Pdf) {
		if err := util.CopyFile(ocr.Pdf, ocr.Source); err != nil {
			return err
		}
	}

	imageFiles, err := t.getImageFiles(workDir, t.imageExt)

	if err != nil {
		return err
	}

	if len(imageFiles) == 0 || t.userConfig.Reprocess {
		env := map[string]string{
			"WORK_DIR": workDir,
		}

		// -density 300 -quiet $pdf_fullname $pdf_nameonly-%04d.$IMAGE_EXT
		imageMagickArgs := []string{
			"imagemagick", "-quiet", "-density", "300",
			pdfFullname, fmt.Sprintf("%v-%%04d%v", noExt, t.imageExt),
		}

		if stdOut, err := t.runDocker(env, imageMagickArgs); err != nil {
			return err
		} else {
			fmt.Println(stdOut)
		}

		imageFiles, err = t.getImageFiles(workDir, t.imageExt)

		if err != nil {
			return err
		}
	}

	ocr.Images = imageFiles
	util.PrintLabeleds("Images", ocr.Images)

	if pather.PathExists(ocr.Text) {
		if err := os.Remove(ocr.Text); err != nil {
			return err
		}
	}

	textFiles, err := t.getImageFiles(workDir, t.textExt)

	if err != nil {
		return err
	}

	if len(textFiles) == 0 || t.userConfig.Reprocess {
		env := map[string]string{
			"WORK_DIR": workDir,
		}

		imageBases := slicer.Map(imageFiles, func(f string) string {
			return pather.Base(f)
		})

		for _, imageBase := range imageBases {
			//  --oem 1 -l eng --psm 6 -c preserve_interword_spaces=1 $image_fullname $image_onlyname
			tesseractArgs := []string{
				"tesseract",
				"--oem", "1", "-l", "eng", "--psm", "6", "-c", "preserve_interword_spaces=1",
				imageBase, extentioner.Strip(imageBase),
			}

			if stdOut, err := t.runDocker(env, tesseractArgs); err != nil {
				return err
			} else {
				fmt.Println(stdOut)
			}
		}

		textFiles, err = t.getImageFiles(workDir, t.textExt)

		if err != nil {
			return err
		}
	}

	ocr.Texts = textFiles
	util.PrintLabeleds("Texts", ocr.Texts)

	if err := util.AppendFile(ocr.Text, textFiles...); err != nil {
		return err
	}

	util.PrintLabeled("Text", ocr.Text)

	return nil
}
