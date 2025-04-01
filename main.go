package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/bregydoc/gtranslate"
)

var (
	foldersName  *string
	languageFrom *string
	languageTo   *string
	wg           sync.WaitGroup
	mu           sync.Mutex
	sem          chan struct{} // 信号量，用于限制并发数
)

type Googlelanguage struct {
	from string
	to   string
}

func charTranslation(glanguage *Googlelanguage, filePath string) {
	defer wg.Done()

	// Read the file contents
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[-] %s Translation failed: %v \n", filePath, err)
		return
	}

	// Regular expression matches Chinese characters
	re := regexp.MustCompile(`[\p{Han}]+`)
	chineseTexts := re.FindAllString(string(content), -1)

	// Translate Chinese to English
	translatedTexts := make([]string, len(chineseTexts))
	for i, text := range chineseTexts {
		translated, err := gtranslate.TranslateWithParams(
			text,
			gtranslate.TranslationParams{
				From: glanguage.from,
				To:   glanguage.to,
			},
		)
		if err != nil {
			fmt.Printf("[-] %s Translation failed: %v \n", filePath, err)
			translatedTexts[i] = text // If the translation fails, keep the original Chinese
		} else {
			translatedTexts[i] = translated
		}
	}

	// Replace the Chinese in the file with the translated English
	translatedContent := string(content)
	for i, text := range chineseTexts {
		translatedContent = strings.Replace(translatedContent, text, translatedTexts[i], 1)
	}

	// Write the translated content back to the file
	err = ioutil.WriteFile(filePath, []byte(translatedContent), 0644)
	if err != nil {
		fmt.Printf("[-] %s Translation failed: %v \n", filePath, err)
		return
	}

	fmt.Printf("[>] translation of %s has been completed and the language in the document has been replaced.\n", filePath)
}

// GetDirAllFilePaths gets all the file paths in the specified directory recursively.
func GetDirAllFilePaths(dirname string) ([]string, error) {
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(infos))
	for _, info := range infos {
		path := dirname + string(os.PathSeparator) + info.Name()
		if info.IsDir() {
			tmp, err := GetDirAllFilePaths(path)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmp...)
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func LogFlag() {
	fmt.Println(`

	 _                       _       _       
	| |                     | |     | |      
	| |_ _ __ __ _ _ __  ___| | __ _| |_ ___ 
	| __| '__/ _' | '_ \/ __| |/ _' | __/ _ \
	| |_| | | \_| | | | \__ \ | \_| | ||  __/
	 \__|_|  \__,_|_| |_|___/_|\__,_|\__\___|      @whale-baby
	
	
	`)
}

func init() {
	foldersName = flag.String("folders", "", "Folder Name")
	languageFrom = flag.String("lfrom", "auto", "Source file language (\"auto\" Automatic identification)")
	languageTo = flag.String("lto", "en", "Language to translate")
}

func main() {
	flag.Parse()

	LogFlag()

	if *foldersName != "" {
		var glanguage Googlelanguage
		glanguage.from = *languageFrom
		glanguage.to = *languageTo

		// Set the maximum number of concurrent goroutines (e.g., 10)
		maxConcurrency := 10
		sem = make(chan struct{}, maxConcurrency)

		// Recurse file directory
		fileAll, err := GetDirAllFilePaths(*foldersName)
		if err != nil {
			fmt.Printf("[-] Failed to recurse file directory %v \n", err)
			return
		}

		// Process files concurrently with semaphore
		for _, path := range fileAll {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire a semaphore
				defer func() { <-sem }() // Release the semaphore when done
				charTranslation(&glanguage, path)
			}(path)
		}

		// Wait for all goroutines to complete
		wg.Wait()

		fmt.Println("[+] Mission successfully completed")
	} else {
		fmt.Println("[-] Enter the name of the folder to be translated")
	}
}
