# go-vscode-reveal-html-postprocessor

The post-processor for the exported html by vscode-reveal.

go-vscode-reveal-html-postprocessor post-processes the index.html to be able to start as stand alone. As a result, your presentation will be distributable to other.

# User requirement

 - vscode-reveal 3.3.2

# Usage

1. Export html and css, js files by vscode-reveal
2. post process by [go-vscode-reveal-html-postprocessor](https://github.com/xshoji/go-vscode-reveal-html-postprocessor/releases)

```
$ ./go-vscode-reveal-html-postprocessor -i /path/to/markdown.md -o /path/to/export
```

3. Open index.html then you will be able to start your presentation on browser!

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