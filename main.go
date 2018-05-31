package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Language string `short:"l" long:"lang" description:"language (default: en)"`
}

type errorContent struct {
	line    int
	column  int
	message string
}

type errorFile struct {
	contents map[string]*errorContent
}

func newErrorFile() *errorFile {
	return &errorFile{contents: map[string]*errorContent{}}
}

type errorsContainer struct {
	files map[string]*errorFile
}

func newErrorsContainer() *errorsContainer {
	return &errorsContainer{map[string]*errorFile{}}
}

func (errs *errorsContainer) printErrors() {
	// TODO: sort
	filenames := make([]string, 0, len(errs.files))
	for k := range errs.files {
		filenames = append(filenames, k)
	}
	sort.Strings(filenames)

	for _, filename := range filenames {
		errors := errs.files[filename]

		keys := make([]string, 0, len(errors.contents))
		for k := range errors.contents {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			contents := errors.contents[k]

			message := strings.Replace(contents.message, "\n", "", -1)
			fmt.Printf("%s:%d:%d: %s\n", filename, contents.line, contents.column, message)
		}
	}
}

func (errs *errorsContainer) addError(filename string, line int, column int, message string) {
	file, ok := errs.files[filename]
	if !ok {
		file = newErrorFile()
		errs.files[filename] = file
	}

	// to prevent duplication
	key := fmt.Sprintf("%08d:%08d:%m", line, column, message)
	_, ok = file.contents[key]
	if !ok {
		file.contents[key] = &errorContent{line, column, message}
	}
}

type xbugsErrorTranslator struct {
	errors *errorsContainer

	bugDescriptions map[string]string
	srcDirs         []string
}

func (xtr *xbugsErrorTranslator) toAbsPath(source string) string {
	minseps := math.MaxInt32
	var res string

	for _, dir := range xtr.srcDirs {
		if strings.HasSuffix(dir, string(os.PathSeparator)+source) {
			seps := strings.Count(dir[0:len(dir)-len(source)-1], string(os.PathSeparator))
			if seps < minseps {
				minseps = seps
				res = dir
			}
		}
	}

	if res == "" {
		return source
	}

	return res
}

func (xtr *xbugsErrorTranslator) addError(bug BugInstance, source SourceLine) int {
	if source.Start == 0 {
		return 0
	}

	sourcePath := xtr.toAbsPath(source.SourcePath)

	message := bug.Type + ": " + xtr.bugDescriptions[bug.Type]

	xtr.errors.addError(sourcePath, source.Start, 1, message)

	return 1
}

func (xtr *xbugsErrorTranslator) parseXbugsErrors(BugCollection BugCollection) {
	for _, bug := range BugCollection.BugInstances {
		counts := 0

		for _, source := range bug.SourceLines {
			counts += xtr.addError(bug, source)
		}
		if counts > 0 {
			continue
		}

		for _, source := range bug.FieldSourceLines {
			counts += xtr.addError(bug, source)
		}
		if counts > 0 {
			continue
		}

		for _, source := range bug.MethodSourceLines {
			counts += xtr.addError(bug, source)
		}
		if counts > 0 {
			continue
		}

		for _, source := range bug.ClassSourceLines {
			counts += xtr.addError(bug, source)
		}
	}
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
		return
	}

	if opts.Language == "" {
		opts.Language = "en"
	}

	var bugDescriptions map[string]string
	switch opts.Language {
	case "en":
		bugDescriptions = BugDescriptionEn
	case "ja":
		bugDescriptions = BugDescriptionJa
	case "fr":
		bugDescriptions = BugDescriptionFr
	default:
		panic("Unsupported language: " + opts.Language)
	}

	errors := newErrorsContainer()

	decoder := xml.NewDecoder(os.Stdin)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}

		switch token := t.(type) {
		case xml.StartElement:
			if token.Name.Local == "BugCollection" {
				document := BugCollection{}
				decoder.DecodeElement(&document, &token)

				translator := xbugsErrorTranslator{errors, bugDescriptions, document.Project.SrcDirs}
				translator.parseXbugsErrors(document)
			}
		}
	}

	errors.printErrors()
}
