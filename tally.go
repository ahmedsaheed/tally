package main

import (
	"bufio"
	"os"
)

func Scan(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil { return 0, err }
	defer file.Close()

	reader := bufio.NewReader(file)
	lineCount := 0
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			if err == os.ErrClosed {
				break
			}
			if err != nil {
				return lineCount, nil
			}
		}
		lineCount++
	}

	return lineCount, nil
}

func TallyDirectory(root string) ([]Language, error) {
	recognisedLanguages := GetAllLanguagesInDir(root)
	for _, language := range recognisedLanguages {
		lang, err := CountUp(language, root)
		if err != nil {
			return recognisedLanguages, err
		}
		if lang.TotalCount != 0 { 
			SetEncounteredLangs(lang)
		}
	}
	return GetEncounteredLangs(), nil
}

