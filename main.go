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
	Input       string `short:"i" long:"input" description:"An input markdown file path." required:"true"`
	Output      string `short:"o" long:"output" description:"An output html directory that exported by vscode-reveal." required:"true"`
	Removelines int    `short:"r" long:"removelines" description:"Remove top lines of the markdown file." default:"7"`
	ImageDir    string `short:"m" long:"imagedir" description:"An image file directory." default:""`
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
	fmt.Println("outut: ", opts.Output)

	minCss, _ := box.Find("css/font-awesome.min.css")
	webfontSvg, _ := box.Find("fonts/fontawesome-webfont.svg")
	otf, _ := box.Find("fonts/FontAwesome.otf")
	webfontWoff2, _ := box.Find("fonts/fontawesome-webfont.woff2")
	webfontWoff, _ := box.Find("fonts/fontawesome-webfont.woff")
	webfontEot, _ := box.Find("fonts/fontawesome-webfont.eot")

	if !strings.HasSuffix(opts.Input, ".md") {
		log.Fatal("<< ERROR >> " + opts.Input + " is not markdown file.")
		os.Exit(1)
	}

	if !Exists(opts.Input) {
		log.Fatal("<< ERROR >> " + opts.Input + " doesn't exist.")
		os.Exit(1)
	}

	if !Exists(opts.Output) {
		log.Fatal("<< ERROR >> " + opts.Output + " doesn't exist.")
		os.Exit(1)
	}

	if !Exists(filepath.Join(opts.Output, "index.html")) {
		log.Fatal("<< ERROR >> index.html doesn't exist in " + opts.Output + ".")
		os.Exit(1)
	}

	os.Mkdir(filepath.Join(opts.Output, "css"), 0777)
	os.Mkdir(filepath.Join(opts.Output, "fonts"), 0777)
	ioutil.WriteFile(filepath.Join(opts.Output, "css", "font-awesome.min.css"), minCss, 0777)
	ioutil.WriteFile(filepath.Join(opts.Output, "fonts", "fontawesome-webfont.svg"), webfontSvg, 0777)
	ioutil.WriteFile(filepath.Join(opts.Output, "fonts", "FontAwesome.otf"), otf, 0777)
	ioutil.WriteFile(filepath.Join(opts.Output, "fonts", "fontawesome-webfont.woff2"), webfontWoff2, 0777)
	ioutil.WriteFile(filepath.Join(opts.Output, "fonts", "fontawesome-webfont.woff"), webfontWoff, 0777)
	ioutil.WriteFile(filepath.Join(opts.Output, "fonts", "fontawesome-webfont.eot"), webfontEot, 0777)

	htmlBytes, _ := ioutil.ReadFile(filepath.Join(opts.Output, "index.html"))
	html := string(htmlBytes)
	re := regexp.MustCompile(`(<section data-markdown="/markdown.md" data-separator="\^\[[\s\S]*</section>)`)
	html = re.ReplaceAllString(html, PLACEHOLDER)
	html = strings.Replace(html, `<link rel="stylesheet" href="libs/reveal.js/font-awesome-4.7.0/css/font-awesome.min.css">`, `<link rel="stylesheet" href="css/font-awesome.min.css">`, -1)

	file, err := os.Open(opts.Input)
	if err != nil {
		log.Panic(err)
	}

	scanner := bufio.NewScanner(file)
	cnt := 0
	markdownString := ""
	isTopEmptyLines := true
	for scanner.Scan() {
		if cnt < opts.Removelines {
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
		panic(err)
	}

	markdownHtml := `
            <!-- https://github.com/hakimel/reveal.js/issues/929#issuecomment-80734215 -->
            <!-- <section data-markdown data-separator="^\n---$" data-separator-vertical="^\n--\n$"> -->
            <section data-markdown data-separator="^\n---$" data-separator-vertical="^\n--\n$">
                <textarea data-template>
` + markdownString + `
                </textarea>
            </section>
`
	fixedHtml := strings.Replace(html, PLACEHOLDER, markdownHtml, -1)
	ioutil.WriteFile(filepath.Join(opts.Output, "index.html"), []byte(fixedHtml), 0777)

	if Exists(opts.ImageDir) && opts.ImageDir != "" {
		CopyDir(opts.ImageDir, opts.Output)
	}
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func CopyDir(src string, dest string) {
	// First, remove target directory
	topDirName := filepath.Base(src)
	destFilePath := filepath.Join(dest, topDirName)
	os.RemoveAll(destFilePath)
	os.Mkdir(destFilePath, 0777)
	doCopyDir(src, destFilePath)
}

func doCopyDir(src string, dest string) {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Fatal("<< ERROR >> " + src + " doesn't exist.")
		os.Exit(1)
	}

	for _, file := range files {
		srcFilePath := filepath.Join(src, file.Name())
		destFilePath := filepath.Join(dest, file.Name())
		if file.IsDir() {
			os.Mkdir(destFilePath, 0777)
			doCopyDir(srcFilePath, destFilePath)
		} else {
			// copy file
			bytes, err := ioutil.ReadFile(srcFilePath)
			if err != nil {
				log.Fatal("<< ERROR >> " + srcFilePath + " cannot read.")
				os.Exit(1)
			}
			ioutil.WriteFile(destFilePath, bytes, 0777)
		}
	}
}
