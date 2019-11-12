# cdiff

[![GoDoc](https://godoc.org/github.com/shibukawa/cdiff?status.svg)](https://godoc.org/github.com/shibukawa/cdiff)

`cdiff` generates line-by-line diff or word-by-word diff (like github) and formats with color.

![screenshot](https://raw.githubusercontent.com/shibukawa/cdiff/master/images/screenshot.png)

## Usage

```go
import (
    "github.com/gookit/color"
)

func main() {
	diff := cdiff.Diff(string(oldDoc), string(newDoc), cdiff.WordByWord)
	color.Print(diff.Unified(oldDocPath, newDocPath, 3, cdiff.GooKitColorTheme))
}
```

## Reference

### func Diff(oldText, newText string, diffType DiffType) Result

It returns diff information from oldText/newText. `diffType` should be `WordByWord` or `LineByLine`.

### Result.Unified(oldDocPath, newDocPath string, keepLines int, theme map[Tag]string) string

It returns string representation of unified format.

`keepLines` is like `diff -U n`. No changed lines count around diffs.

`theme` is a style text to create string representation. There is `GooKitColorTheme` and `HTMLTheme`.
`GooKitColorTheme` has tags of github.com/gookit/color.

### Result.Format(theme map[Tag]string) string

It doesn't omit unchanged line blocks instead of `Unified()`.

### Result.String() sring

It is similar to `Format()` but it desn't have extra text like theme.

## License

Apache 2