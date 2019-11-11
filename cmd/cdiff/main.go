package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gookit/color"
	"github.com/shibukawa/cdiff"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	length     = kingpin.Flag("unified", "output NUM (default 3) lines of unified context").Short('U').Default("3").PlaceHolder("NUM").Int()
	oldDocPath = kingpin.Arg("OLD", "old file to compare").Required().ExistingFile()
	newDocPath = kingpin.Arg("NEW", "new file to compare").Required().ExistingFile()
)

func main() {
	kingpin.Parse()
	oldDoc, err := ioutil.ReadFile(*oldDocPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open old document %q: %v", *oldDocPath, err)
		os.Exit(1)
	}
	newDoc, err := ioutil.ReadFile(*newDocPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open new document %q: %v", *newDocPath, err)
		os.Exit(1)
	}
	diff := cdiff.Diff(string(oldDoc), string(newDoc), cdiff.WordByWord)
	color.Print(diff.Unified(*oldDocPath, *newDocPath, *length, cdiff.GooKitColorTheme))
}
