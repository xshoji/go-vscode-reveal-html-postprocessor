package main

import (
	"bufio"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const PLACEHOLDER = "XXXXXXXXXX"

type options struct {
	Input                 string `short:"i" long:"input" description:"[ required ] An input markdown file path." required:"true"`
	Output                string `short:"o" long:"output" description:"[ required ] An output html directory that exported by vscode-reveal." required:"true"`
	RemoveLines           int    `short:"r" long:"removeLines" description:"Remove top lines of the markdown file." default:"7"`
	DataSeparator         string `short:"d" long:"dataSeparator" description:"The data-separator." default:"---"`
	DataSeparatorVertical string `short:"v" long:"dataSeparatorVertical" description:"The data-separator-vertical." default:"--"`
	ImageDir              string `short:"m" long:"imageDir" description:"An image file directory." default:""`
}

var box = packr.NewBox("./font-awesome-4.7.0")

func main() {

	opts := *new(options)
	parser := flags.NewParser(&opts, flags.Default)
	// set name
	parser.Name = "go-vscode-reveal-html-postprocessor"
	if _, err := parser.Parse(); err != nil {
		flagsError, _ := err.(*flags.Error)
		// help時は何もしない
		if flagsError.Type == flags.ErrHelp {
			return
		}
		fmt.Println()
		parser.WriteHelp(os.Stdout)
		fmt.Println()
		return
	}
	fmt.Println("input: ", opts.Input)
	fmt.Println("output: ", opts.Output)

	if !strings.HasSuffix(opts.Input, ".md") {
		log.Fatal("<< ERROR >> " + opts.Input + " is not markdown file.")
	}

	if !exists(opts.Input) {
		log.Fatal("<< ERROR >> " + opts.Input + " doesn't exist.")
	}

	if !exists(opts.Output) {
		log.Fatal("<< ERROR >> " + opts.Output + " doesn't exist.")
	}

	if !exists(filepath.Join(opts.Output, "index.html")) {
		log.Fatal("<< ERROR >> index.html doesn't exist in " + opts.Output + ".")
	}

	makeDirectory(filepath.Join(opts.Output, "css"))
	makeDirectory(filepath.Join(opts.Output, "fonts"))

	copyFileFromBox("css/font-awesome.min.css", filepath.Join(opts.Output, "css", "font-awesome.min.css"))
	copyFileFromBox("fonts/fontawesome-webfont.svg", filepath.Join(opts.Output, "fonts", "fontawesome-webfont.svg"))
	copyFileFromBox("fonts/FontAwesome.otf", filepath.Join(opts.Output, "fonts", "FontAwesome.otf"))
	copyFileFromBox("fonts/fontawesome-webfont.woff2", filepath.Join(opts.Output, "fonts", "fontawesome-webfont.woff2"))
	copyFileFromBox("fonts/fontawesome-webfont.woff", filepath.Join(opts.Output, "fonts", "fontawesome-webfont.woff"))
	copyFileFromBox("fonts/fontawesome-webfont.eot", filepath.Join(opts.Output, "fonts", "fontawesome-webfont.eot"))

	htmlBytes, _ := ioutil.ReadFile(filepath.Join(opts.Output, "index.html"))
	html := string(htmlBytes)
	re := regexp.MustCompile(`(<section data-markdown="/markdown.md" data-separator="\^\[[\s\S]*</section>)`)
	html = re.ReplaceAllString(html, PLACEHOLDER)
	html = strings.Replace(html, `<link rel="stylesheet" href="libs/reveal.js/font-awesome-4.7.0/css/font-awesome.min.css">`, `<link rel="stylesheet" href="css/font-awesome.min.css">`, -1)

	file, err := os.Open(opts.Input)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	cnt := 0
	markdownString := ""
	isTopEmptyLines := true
	for scanner.Scan() {
		if cnt < opts.RemoveLines {
			cnt++
			continue
		}
		line := scanner.Text()
		if isTopEmptyLines && line == "" {
			continue
		}
		isTopEmptyLines = false
		markdownString = markdownString + "                    " + line + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	markdownHtml := `
            <!-- https://github.com/hakimel/reveal.js/issues/929#issuecomment-80734215 -->
            <section data-markdown data-separator="^\n` + opts.DataSeparator + `$" data-separator-vertical="^\n` + opts.DataSeparatorVertical + `$">
                <textarea data-template>
` + markdownString + `
                </textarea>
            </section>
`
	fixedHtml := strings.Replace(html, PLACEHOLDER, markdownHtml, -1)
	writeFile(filepath.Join(opts.Output, "index.html"), []byte(fixedHtml))

	if exists(opts.ImageDir) && opts.ImageDir != "" {
		copyDir(opts.ImageDir, opts.Output)
	}
}

func makeDirectory(path string) {
	err := os.Mkdir(path, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func writeFile(path string, contentBytes []byte) {
	err := ioutil.WriteFile(path, contentBytes, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func copyFileFromBox(boxPath string, outputPath string) {
	bytes, _ := box.Find(boxPath)
	writeFile(outputPath, bytes)
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func copyDir(src string, dest string) {
	// First, remove target directory
	topDirName := filepath.Base(src)
	destFilePath := filepath.Join(dest, topDirName)
	err := os.RemoveAll(destFilePath)
	if err != nil {
		log.Fatal(err)

	}
	makeDirectory(destFilePath)
	doCopyDir(src, destFilePath)
}

func doCopyDir(src string, dest string) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal("<< ERROR >> " + src + " doesn't exist.")
	}

	for _, file := range files {
		srcFilePath := filepath.Join(src, file.Name())
		destFilePath := filepath.Join(dest, file.Name())
		if file.IsDir() {
			makeDirectory(destFilePath)
			doCopyDir(srcFilePath, destFilePath)
		} else {
			// copy file
			bytes, err := ioutil.ReadFile(srcFilePath)
			if err != nil {
				log.Fatal("<< ERROR >> " + srcFilePath + " cannot read.")
			}
			writeFile(destFilePath, bytes)
		}
	}
}
