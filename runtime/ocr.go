package runtime

import (
	"fmt"
	"os"

	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/docsorter/util"
	"github.com/boundedinfinity/go-commoner/extentioner"
	"github.com/boundedinfinity/go-commoner/pather"
	"github.com/boundedinfinity/go-commoner/slicer"
)

func (t *Runtime) OcrSingle(ocr *model.OcrContext) error {
	if err := t.prepareDirectory(ocr); err != nil {
		return err
	}

	if err := t.pdf2Images(ocr); err != nil {
		return err
	}

	if err := t.images2Text(ocr); err != nil {
		return err
	}

	if err := util.AppendFile(ocr.WorkText, ocr.WorkTexts...); err != nil {
		return err
	}

	util.PrintLabeled("Text", ocr.WorkText)

	return nil
}

func (t *Runtime) prepareDirectory(ocr *model.OcrContext) error {
	workDir := extentioner.Strip(pather.Base(ocr.Source))
	ocr.WorkDir = pather.Join(t.userConfig.WorkPath, workDir)
	ocr.WorkPdf = pather.Join(ocr.WorkDir, pather.Base(ocr.Source))

	if t.userConfig.Reprocess {
		if err := util.EnsureDelete(ocr.WorkDir); err != nil {
			return err
		}
	}

	if err := pather.DirEnsure(ocr.WorkDir); err != nil {
		return err
	}

	if t.userConfig.Reprocess || !pather.PathExists(ocr.WorkPdf) {
		if err := util.CopyFile(ocr.WorkPdf, ocr.Source); err != nil {
			return err
		}
	}

	return nil
}

func (t *Runtime) pdf2Images(ocr *model.OcrContext) error {
	imageFiles, err := util.GetFilteredFiles(ocr.WorkDir, t.extImage)

	if err != nil {
		return err
	}

	if len(imageFiles) == 0 || t.userConfig.Reprocess {
		env := map[string]string{
			"WORK_DIR": ocr.WorkDir,
		}

		// -density 300 -quiet $pdf_fullname $pdf_nameonly-%04d.$IMAGE_EXT
		imageMagickArgs := []string{
			"imagemagick", "-quiet", "-density", "300",
			pather.Base(ocr.Source),
			fmt.Sprintf("%v-%%04d%v", pather.Base(ocr.WorkDir), t.extImage),
		}

		if stdOut, err := t.runDocker(env, imageMagickArgs); err != nil {
			return err
		} else {
			fmt.Println(stdOut)
		}

		imageFiles, err = util.GetFilteredFiles(ocr.WorkDir, t.extImage)

		if err != nil {
			return err
		}
	}

	ocr.WorkImages = imageFiles
	util.PrintLabeleds("Images", ocr.WorkImages)

	return nil
}

func (t *Runtime) images2Text(ocr *model.OcrContext) error {
	ocr.WorkText = extentioner.Swap(ocr.WorkPdf, t.extPdf, t.extText)

	if pather.PathExists(ocr.WorkText) {
		if err := os.Remove(ocr.WorkText); err != nil {
			return err
		}
	}

	textFiles, err := util.GetFilteredFiles(ocr.WorkDir, t.extText)

	if err != nil {
		return err
	}

	if len(textFiles) == 0 || t.userConfig.Reprocess {
		env := map[string]string{
			"WORK_DIR": ocr.WorkDir,
		}

		imageBases := slicer.Map(ocr.WorkImages, func(f string) string {
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

		textFiles, err = util.GetFilteredFiles(ocr.WorkDir, t.extText)

		if err != nil {
			return err
		}
	}

	ocr.WorkTexts = textFiles
	util.PrintLabeleds("Texts", ocr.WorkTexts)

	return nil
}
