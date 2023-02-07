package runtime

import (
	"github.com/boundedinfinity/docsorter/model"
	"github.com/boundedinfinity/go-commoner/pather"
)

func (t *Runtime) LoadFiles() ([]model.OcrContext, error) {
	var ocrs []model.OcrContext
	allPaths := make([]string, 0)

	for _, path := range t.userConfig.InputPaths {
		paths, err := pather.GetFiles(path)

		if err != nil {
			return ocrs, err
		}

		allPaths = append(allPaths, paths...)
	}

	for _, path := range allPaths {
		ocr := model.NewOcrContext()
		ocr.Source = path
		ocrs = append(ocrs, *ocr)
	}

	return ocrs, nil
}

// func move(runtime *Runtime) error {
// 	for path, desc := range runtime.files {
// 		if !desc.Classified {
// 			continue
// 		}

// 		bs, err := ioutil.ReadFile(path)

// 		if err != nil {
// 			return err
// 		}

// 		acct := desc.Account[len(desc.Account)-4:]
// 		date := desc.Date.Format("2006.01.02")
// 		filename := fmt.Sprintf("%v-%v-%v.pdf", acct, date, desc.Name)
// 		newPath := filepath.Join(runtime.config.OutputPath, filename)

// 		fmt.Printf("Writing: %v\n", newPath)

// 		if err := ioutil.WriteFile(newPath, bs, 0755); err != nil {
// 			return err
// 		}

// 		if err := clean(runtime, path); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func clean(runtime *Runtime, path string) error {
// 	paths := []string{
// 		path,
// 		runtime.orig2txt(path),
// 		runtime.orig2sum(path),
// 	}

// 	for _, p := range paths {
// 		if fu.PathExists(p) {
// 			if err := os.Remove(p); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
