package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/bregydoc/gtranslate"
)

var foldersName *string
var languageFrom *string
var languageTo *string

type Googlelanguage struct {
	from string
	to   string
}

func charTranslation(glanguage *Googlelanguage, filePath string) {

	// Read the file contents
	// filePath := "test.go" // Replace with your file path
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[-] %s Translation failed: %v \n", filePath, err)
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
				// From: "auto",
				// To:   "en",
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
	}

	fmt.Printf("[>] translation of %s has been completed and the language in the document has been replaced.\n", filePath)

}

// GetDirAllFilePaths gets all the file paths in the specified directory recursively.
func GetDirAllFilePaths(dirname string) ([]string, error) {
	// Remove the trailing path separator if dirname has.
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
	
                      ███
                    ███
                   ██
                  ██
                 ███
        █████████████████████
     ███████████████████████████
   ███████████████ ███████████████
 ██████████ ██████ ██████ █████████
 █████████   █████ █████   █████████
█████████     ████ ████     ████████
████████       ███ ███       ████████
██████████████████ ██████████████████
█████████████████   █████████████████
████████   ████████████████  ████████
 ███████                     ███████
  ████████                 ████████
   ██████████████████████████████
    ███████████████████████████                  @whale-baby

	
	
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

		//Set language
		glanguage.from = *languageFrom
		glanguage.to = *languageTo

		//recurse file directory
		fileAll, err := GetDirAllFilePaths(*foldersName)
		if err != nil {
			fmt.Printf("[-] Failed to recurse file directory %v \n", err)
			return
		}

		for _, path := range fileAll {

			charTranslation(&glanguage, path)

		}

		fmt.Println("[+] Mission successfully completed")

	} else {

		fmt.Println("[-] Enter the name of the folder to be translated")

	}

}
