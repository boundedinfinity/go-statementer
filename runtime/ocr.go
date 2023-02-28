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

func (t *Runtime) OcrSingle(stage *model.FileSet) error {
	if err := t.prepareDirectory(stage); err != nil {
		return err
	}

	if err := t.pdf2Images(stage); err != nil {
		return err
	}

	if err := t.images2Text(stage); err != nil {
		return err
	}

	if err := util.AppendFile(stage.Text, stage.Texts...); err != nil {
		return err
	}

	util.PrintLabeled("Text", stage.Text)

	return nil
}

func (t *Runtime) pdf2Images(stage *model.FileSet) error {
	imageFiles, err := util.GetFilteredFiles(stage.Dir, t.extImage)

	if err != nil {
		return err
	}

	if len(imageFiles) == 0 || t.UserConfig.Reprocess {
		env := map[string]string{
			"WORK_DIR": stage.Dir,
		}

		// -density 300 -quiet $pdf_fullname $pdf_nameonly-%04d.$IMAGE_EXT
		imageMagickArgs := []string{
			"imagemagick", "-quiet", "-density", "300",
			pather.Base(stage.Source),
			fmt.Sprintf("%v-%%04d%v", pather.Base(stage.Dir), t.extImage),
		}

		if stdOut, err := t.runDocker(env, imageMagickArgs); err != nil {
			return err
		} else {
			fmt.Println(stdOut)
		}

		imageFiles, err = util.GetFilteredFiles(stage.Dir, t.extImage)

		if err != nil {
			return err
		}
	}

	stage.Images = imageFiles
	util.PrintLabeleds("Images", stage.Images)

	return nil
}

func (t *Runtime) images2Text(stage *model.FileSet) error {
	if pather.PathExists(stage.Text) {
		if err := os.Remove(stage.Text); err != nil {
			return err
		}
	}

	textFiles, err := util.GetFilteredFiles(stage.Dir, t.extText)

	if err != nil {
		return err
	}

	if len(textFiles) == 0 || t.UserConfig.Reprocess {
		env := map[string]string{
			"WORK_DIR": stage.Dir,
		}

		imageBases := slicer.Map(stage.Images, func(f string) string {
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

		textFiles, err = util.GetFilteredFiles(stage.Dir, t.extText)

		if err != nil {
			return err
		}
	}

	stage.Texts = textFiles
	util.PrintLabeleds("Texts", stage.Texts)

	return nil
}
