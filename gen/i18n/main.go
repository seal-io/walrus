package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
	"golang.org/x/text/message/pipeline"

	"github.com/seal-io/walrus/utils/log"
)

func main() {
	err := generate()
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}
}

// generate takes from golang.org/x/text/cmd/gotext/update.go,
// and updates the i18n catalog.go file,
// which avoids compiling gotext according to go runtime.
//
// See https://github.com/golang/go/issues/58751#issuecomment-1706999475.
func generate() error {
	// Prepare.
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}

	cfg := pipeline.Config{
		Packages:            []string{"github.com/seal-io/walrus/pkg/i18n"}, // The packages to scan.
		Dir:                 filepath.Join(pwd, "pkg/i18n/locales"),         // Find text and translation file.
		TranslationsPattern: `messages\.(.*)\.json$`,                        // The regex pattern to match translation files.
		GenPackage:          "github.com/seal-io/walrus/pkg/i18n",           // The package name of the generated file.
		GenFile:             filepath.Join(pwd, "pkg/i18n/catalog.go"),      // The file name of the generated file.
		SourceLanguage:      language.English,
		Supported:           []language.Tag{language.English, language.Chinese},
	}

	state, err := pipeline.Extract(&cfg)
	if err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	if err = state.Import(); err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	if err = state.Merge(); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	if err = state.Export(); err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	if err = state.Generate(); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	return nil
}
