package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	htmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	extensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_HARD_LINE_BREAK
)

func main() {
	// Parse command line flags and handle default values (if necessary)
	pathPtr := flag.String("input", "", "Filepath for the markdown file you would like to convert")
	outputPtr := flag.String("output", "", "Path for your html output")
	flag.Parse()

	// Fatal error if no input has been specified
	if *pathPtr == "" {
		log.Fatal(errors.New("Error: You must provide a path to your markdown file"))
	}
	// Check the input path and sanitize the filepath
	path := cleanPath(pathPtr)

	// If no output path was provided, override the default pointer and create a new file at the same
	// location as the input file
	outputPath := outputFilePath(outputPtr, path)
	// Read the markdown file into a []Byte
	rawMarkdown := readInput(path)
	// Create a new HmtlRender struct with the htmlFlags const specified above
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")
	// Create an Option struct with the extensions const specified above
	outputOpts := blackfriday.Options{
		Extensions: extensions,
	}
	// Render the markdown file into HTML and return a new []Byte
	htmlBody := blackfriday.MarkdownOptions(rawMarkdown,
		renderer,
		outputOpts,
	)
	// Print in-progress message to stdOut
	fmt.Printf("Converting %s. Output is located at %s\n",
		filepath.Base(path),
		outputPath)
	// Write the htmlBody []Byte out to disk with -rw-r--r-- permissions.
	// 0644 is the standard permission config on for html files on Apache. Seemed like a safe bet.
	ioutil.WriteFile(outputPath, htmlBody, 0644)
}

func cleanPath(pathPtr *string) (cleaned string) {
	// Make sure the input file is a Markdown file
	if filepath.Ext(*pathPtr) != ".md" {
		log.Fatal(errors.New("Error: You must provide a markdown (.md) file."))
	}
	// Tilde expansion for unix
	if string((*pathPtr)[0]) == "~" {
		expandTilde(pathPtr)
	}
	cleaned = filepath.Clean(*pathPtr)
	return
}

func expandTilde(pathPtr *string) {
	usr, _ := user.Current()
	*pathPtr = strings.Replace(*pathPtr, "~", usr.HomeDir, 1)
}

func readInput(filename string) (rawMarkdown []byte) {
	rawMarkdown, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func outputFilePath(outputPtr *string, inputPath string) (output string) {
	if *outputPtr == "" {
		inputDir, inputFile := filepath.Split(inputPath)
		outputFile := strings.TrimSuffix(inputFile, ".md") + ".html"
		output = filepath.Join(inputDir, outputFile)
	} else {
		output = *outputPtr
	}
	output, _ = filepath.Abs(output)
	return
}
