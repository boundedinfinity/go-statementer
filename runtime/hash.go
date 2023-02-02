package runtime

// func getCurrSums(runtime *Runtime) error {
// 	processFn := func(path string) error {
// 		fmt.Println(path)
// 		sum, err := calculateSum(path)

// 		if err != nil {
// 			return err
// 		}

// 		fmt.Printf("    sum: %v\n", sum)

// 		fdesc := runtime.files[path]
// 		fdesc.Sum = sum

// 		return nil
// 	}

// 	return fu.WalkFilePaths(runtime.config.InputPath, inputExtFilter(runtime), processFn)
// }

// func writeSumFiles(runtime *Runtime) error {
// 	currPaths := mu.Keys(runtime.files)

// 	createSumFn := func(path string) error {
// 		fdesc := runtime.files[path]
// 		sumPath := runtime.orig2sum(path)
// 		return createSumFile(*runtime, sumPath, fdesc.Sum)
// 	}

// 	return su.EachErr(currPaths, createSumFn)
// }

// func compareSums(runtime *Runtime) error {
// 	for path1, fdesc1 := range runtime.files {
// 		fmt.Printf("Checking %v\n", path1)

// 		for path2, fdesc2 := range runtime.files {
// 			if path1 == path2 {
// 				continue
// 			}

// 			if fdesc1.Sum == fdesc2.Sum {
// 				fmt.Printf("    matched %v\n", path2)
// 			}
// 		}
// 	}

// 	return nil
// }
