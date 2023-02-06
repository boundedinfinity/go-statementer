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

	if err := util.AppendFile(ocr.Text, ocr.Texts...); err != nil {
		return err
	}

	util.PrintLabeled("Text", ocr.Text)

	return nil
}

func (t *Runtime) prepareDirectory(ocr *model.OcrContext) error {
	workDir := extentioner.Strip(pather.Base(ocr.Source))
	ocr.WorkDir = pather.Join(t.userConfig.WorkPath, workDir)
	ocr.Pdf = pather.Join(ocr.WorkDir, pather.Base(ocr.Source))

	if t.userConfig.Reprocess {
		if err := util.EnsureDelete(ocr.WorkDir); err != nil {
			return err
		}
	}

	if err := pather.DirEnsure(ocr.WorkDir); err != nil {
		return err
	}

	if t.userConfig.Reprocess || !pather.PathExists(ocr.Pdf) {
		if err := util.CopyFile(ocr.Pdf, ocr.Source); err != nil {
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

	ocr.Images = imageFiles
	util.PrintLabeleds("Images", ocr.Images)

	return nil
}

func (t *Runtime) images2Text(ocr *model.OcrContext) error {
	ocr.Text = extentioner.Swap(ocr.Pdf, t.extPdf, t.extText)

	if pather.PathExists(ocr.Text) {
		if err := os.Remove(ocr.Text); err != nil {
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

		imageBases := slicer.Map(ocr.Images, func(f string) string {
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

	ocr.Texts = textFiles
	util.PrintLabeleds("Texts", ocr.Texts)

	return nil
}
