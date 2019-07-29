package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gobuffalo/packr"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const PLACEHOLDER = "XXXXXXXXXX"

var helpFlag = flag.Bool("help", false, "help")
var inputMarkdownFlag = flag.String("input", "", "[required] Input file (.md) path.")
var outputDirectoryFlag = flag.String("output", "", "[required] Output directory ( exported html directory by vscode-reveal ) path.")
var removeHeaderLinesFlag = flag.Int("remove-header-lines", 7, "[optional] Remove header lines at markdown file..")
var box = packr.NewBox("./font-awesome-4.7.0")

func init() {
	flag.BoolVar(helpFlag, "h", false, "= -help")
	flag.StringVar(inputMarkdownFlag, "i", "", "= -input")
	flag.StringVar(outputDirectoryFlag, "o", "", "= -output")
	flag.IntVar(removeHeaderLinesFlag, "r", 7, "= -remove-header-lines")
}

func main() {

	flag.Parse()
	if *helpFlag || *inputMarkdownFlag == "" || *outputDirectoryFlag == "" {
		flag.Usage()
		os.Exit(0)
	}
	fmt.Println("input: ", *inputMarkdownFlag)
	fmt.Println("outut: ", *outputDirectoryFlag)

	minCss, _ := box.Find("css/font-awesome.min.css")
	webfontSvg, _ := box.Find("fonts/fontawesome-webfont.svg")
	otf, _ := box.Find("fonts/FontAwesome.otf")
	webfontWoff2, _ := box.Find("fonts/fontawesome-webfont.woff2")
	webfontWoff, _ := box.Find("fonts/fontawesome-webfont.woff")
	webfontEot, _ := box.Find("fonts/fontawesome-webfont.eot")

	if !strings.HasSuffix(*inputMarkdownFlag, ".md") {
		log.Fatal("<< ERROR >> " + *inputMarkdownFlag + " is not markdown file.")
		os.Exit(1)
	}

	if !Exists(*inputMarkdownFlag) {
		log.Fatal("<< ERROR >> " + *inputMarkdownFlag + " is not exist.")
		os.Exit(1)
	}

	if !Exists(*outputDirectoryFlag) {
		log.Fatal("<< ERROR >> " + *outputDirectoryFlag + " is not exist.")
		os.Exit(1)
	}

	if !Exists(filepath.Join(*outputDirectoryFlag, "index.html")) {
		log.Fatal("<< ERROR >> index.html is not exist in " + *outputDirectoryFlag + ".")
		os.Exit(1)
	}

	os.Mkdir(filepath.Join(*outputDirectoryFlag, "css"), 0777)
	os.Mkdir(filepath.Join(*outputDirectoryFlag, "fonts"), 0777)
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "css", "font-awesome.min.css"), minCss, 0777)
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "fonts", "fontawesome-webfont.svg"), webfontSvg, 0777)
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "fonts", "FontAwesome.otf"), otf, 0777)
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "fonts", "fontawesome-webfont.woff2"), webfontWoff2, 0777)
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "fonts", "fontawesome-webfont.woff"), webfontWoff, 0777)
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "fonts", "fontawesome-webfont.eot"), webfontEot, 0777)

	htmlBytes, _ := ioutil.ReadFile(filepath.Join(*outputDirectoryFlag, "index.html"))
	html := string(htmlBytes)
	re := regexp.MustCompile(`(<section data-markdown="/markdown.md" data-separator="\^\[[\s\S]*</section>)`)
	html = re.ReplaceAllString(html, PLACEHOLDER)
	html = strings.Replace(html, `<link rel="stylesheet" href="libs/reveal.js/font-awesome-4.7.0/css/font-awesome.min.css">`, `<link rel="stylesheet" href="css/font-awesome.min.css">`, -1)

	file, err := os.Open(*inputMarkdownFlag)
	if err != nil {
		log.Panic(err)
	}

	scanner := bufio.NewScanner(file)
	cnt := 0
	markdownString := ""
	isTopEmptyLines := true
	for scanner.Scan() {
		if cnt < *removeHeaderLinesFlag {
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
	ioutil.WriteFile(filepath.Join(*outputDirectoryFlag, "index.html"), []byte(fixedHtml), 0777)

}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
