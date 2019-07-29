# go-vscode-reveal-html-postprocessor

The post processor for an exported html by vscode reveal.

# User requirement

 - vscode-reveal 3.3.2

# Usage

1. export html file by vscode-reveal
2. post process by this script

```
$ ./go-vscode-reveal-html-postprocessor -i /path/to/markdown.md -o /path/to/export
```

# Build requirement

 - dep
 - packr
 - goreleaser
 - go 1.12
