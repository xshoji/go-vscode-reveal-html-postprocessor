# go-vscode-reveal-html-postprocessor

The post-processor for the exported html by vscode-reveal.

go-vscode-reveal-html-postprocessor post-processes the index.html generated by vscode-reveal to be able to start in stand-alone. 
As a result, your presentation will be distributable to others.

# User requirement

 - vscode-reveal 3.3.2
   - `1.` [download vscode-reveal 3.3.2](https://evilz.gallery.vsassets.io/_apis/public/gallery/publisher/evilz/extension/vscode-reveal/3.3.2/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage)
   - `2.` `mv Microsoft.VisualStudio.Services.VSIXPackage Microsoft.VisualStudio.Services.VSIX`
   - `3.` Start Visual Studio Code, then F1 -> `Extensions: Install from VSIX...`

# Usage

1. Export html and css, js files by vscode-reveal.
2. Post-process by [go-vscode-reveal-html-postprocessor](https://github.com/xshoji/go-vscode-reveal-html-postprocessor/releases).

```
$ ./go-vscode-reveal-html-postprocessor -i /path/to/markdown.md -o /path/to/export
```

3. Open index.html then you will be able to start your presentation on browser!

### Tool usage

```
Usage:
  go-vscode-reveal-html-postprocessor [OPTIONS]

Application Options:
  -i, --input=                 [ required ] An input markdown file path.
  -o, --output=                [ required ] An output html directory that exported by vscode-reveal.
  -r, --removeLines=           Remove top lines of the markdown file. (default: 7)
  -d, --dataSeparator=         The data-separator. (default: ---)
  -v, --dataSeparatorVertical= The data-separator-vertical. (default: --)
  -m, --imageDir=              An image file directory.

Help Options:
  -h, --help                   Show this help message
```


# Build requirement

 - dep
 - packr
 - goreleaser
 - go 1.12

# Release

```
// create tag, push tag, release
git tag v0.0.1 && git push --tags && goreleaser --rm-dist
```

# Delete tag

```
git tag -d v0.0.1 && git push origin :v0.0.1
```