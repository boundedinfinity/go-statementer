package runtime

// func classify(runtime *Runtime) error {
// 	for path, desc := range runtime.files {
// 		fmt.Printf("classify: %v\n", path)

// 		if len(desc.PropList) <= 0 {
// 			fmt.Printf("    0 classifications\n")
// 			continue
// 		} else if len(desc.PropList) > 1 {
// 			fmt.Printf("    more than 1 classification\n")
// 			continue
// 		}

// 		match, mok := su.FindFn(runtime.config.Match, func(m ConfigMatchDescriptor) bool {
// 			return m.Name == desc.PropList[0].Name
// 		})

// 		if !mok {
// 			continue
// 		}

// 		if len(desc.PropList[0].Props) != len(match.Patterns) {
// 			continue
// 		}

// 		desc.Account = desc.PropList[0].Props["AccountNumber"]
// 		desc.Name = desc.PropList[0].Name

// 		if mok {
// 			prop, pok := desc.PropList[0].Props[match.DateKey]
// 			dateStr := strings.Title(prop)

// 			if pok {
// 				date, err := time.Parse(match.DateFormat, dateStr)

// 				if err != nil {
// 					return err
// 				}

// 				desc.Date = date
// 			}
// 		} else {
// 			continue
// 		}

// 		desc.Classified = true
// 	}

// 	return nil
// }

// func findRegexMatches(runtime *Runtime) error {
// 	processFn := func(path string) error {
// 		fmt.Printf("regex: %v\n", path)
// 		finfo := runtime.files[path]

// 		err := hasRegexMatch(runtime, path, &finfo.PropList)

// 		if err != nil {
// 			return err
// 		}

// 		for _, props := range finfo.PropList {
// 			fmt.Printf("    %v\n", props.Name)
// 			for k, v := range props.Props {
// 				fmt.Printf("        %v: %v\n", k, v)
// 			}
// 		}

// 		return nil
// 	}

// 	return fu.WalkFilePaths(runtime.config.InputPath, inputExtFilter(runtime), processFn)
// }

////////////////////////////////////////////////////////
// Utilities
////////////////////////////////////////////////////////

// func inputExtFilter(runtime *Runtime) func(string) bool {
// 	return func(path string) bool {
// 		return runtime.hasInputExt(path) && !runtime.isIgnorePath(path)
// 	}
// }

// func lcFile(runtime *Runtime, path string) error {
// 	textPath := runtime.orig2txt(path)
// 	bs, err := os.ReadFile(textPath)

// 	if err != nil {
// 		return err
// 	}

// 	bs = bytes.ToLower(bs)
// 	err = os.WriteFile(textPath, bs, 0755)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func getSumMatches(sums map[string]string, path, sum string) []string {
// 	matches := []string{}

// 	for previousPath, previousSum := range sums {
// 		if path != previousPath && sum == previousSum {
// 			matches = append(matches, previousPath)
// 		}
// 	}

// 	return matches
// }

// func hasRegexMatch(runtime *Runtime, path string, propList *[]FileProps) error {
// 	textPath := runtime.orig2txt(path)
// 	bs, err := os.ReadFile(textPath)
// 	content := string(bs)

// 	if err != nil {
// 		return err
// 	}

// 	for _, matchConfig := range runtime.config.Match {
// 		props := FileProps{
// 			Name:  matchConfig.Name,
// 			Props: map[string]string{},
// 		}

// 		for _, pattern := range matchConfig.Patterns {
// 			re, err := regexp.Compile(pattern.Pattern)

// 			if err != nil {
// 				return err
// 			}

// 			groups := re.FindStringSubmatch(content)

// 			if len(groups) == 0 {
// 				// No match
// 			} else if len(groups) == 1 {
// 				props.Props[pattern.Label] = strings.TrimSpace(groups[0])
// 			} else if len(groups) == 2 {
// 				props.Props[pattern.Label] = strings.TrimSpace(groups[1])
// 			} else {
// 				fmt.Printf("more than 2 groups: %v", groups)
// 			}
// 		}

// 		if len(props.Props) == len(matchConfig.Patterns) {
// 			*propList = append(*propList, props)
// 		}
// 	}

// 	return nil
// }

// func createSumFile(runtime Runtime, path string, sum string) error {
// 	sumPath := runtime.orig2sum(path)

// 	if err := os.WriteFile(sumPath, []byte(sum), 0755); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func removeSumFile(runtime Runtime, path string) error {
// 	sumPath := runtime.orig2sum(path)

// 	if fu.PathExists(sumPath) {
// 		if err := os.Remove(sumPath); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func getSumfileContents(path string) (string, error) {
// 	bs, err := os.ReadFile(path)

// 	if err != nil {
// 		return "", nil
// 	}

// 	return string(bs), nil
// }

// func calculateSum(path string) (string, error) {
// 	hasher := sha256.New()
// 	input, err := os.Open(path)

// 	if err != nil {
// 		return "", err
// 	}

// 	defer input.Close()

// 	if _, err := io.Copy(hasher, input); err != nil {
// 		return "", err
// 	}

// 	sum := hasher.Sum(nil)
// 	encoded := hex.EncodeToString(sum)
// 	return encoded, nil
// }
