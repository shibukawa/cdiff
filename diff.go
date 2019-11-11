package cdiff

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type blockDiff struct {
	Ope           Ope
	Text          string
	NewLineNumber int
	OldLineNumber int
}

func calcBlockDiff(oldText, newText string) []blockDiff {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(oldText, newText)
	diffs := dmp.DiffMain(a, b, true)
	diffByLines := dmp.DiffCharsToLines(diffs, c)
	result := make([]blockDiff, len(diffByLines))
	newLineNum := 1
	oldLineNum := 1
	for i, diff := range diffByLines {
		f := blockDiff{
			Ope:           Ope(diff.Type),
			Text:          diff.Text,
			NewLineNumber: -1,
			OldLineNumber: -1,
		}
		inc := strings.Count(diff.Text, "\n")
		switch f.Ope {
		case Insert:
			f.NewLineNumber = newLineNum
			newLineNum += inc
		case Delete:
			f.OldLineNumber = oldLineNum
			oldLineNum += inc
		case Keep:
			f.NewLineNumber = newLineNum
			f.OldLineNumber = oldLineNum
			newLineNum += inc
			oldLineNum += inc
		}
		result[i] = f
	}
	return result
}

// Fragment contains the smallest text fragment of diff
type Fragment struct {
	Text    string
	Changed bool
}

// Line represents a line
type Line struct {
	Ope           Ope
	NewLineNumber int
	OldLineNumber int
	Fragments     []Fragment
}

func (d Line) String() string {
	var builder strings.Builder
	for _, f := range d.Fragments {
		builder.WriteString(f.Text)
	}
	return builder.String()
}

// Result contains diff result
type Result struct {
	Lines []Line
}

func (r Result) String() string {
	var builder strings.Builder
	for _, l := range r.Lines {
		switch l.Ope {
		case Insert:
			builder.WriteString("+ ")
		case Delete:
			builder.WriteString("- ")
		case Keep:
			builder.WriteString("  ")
		}
		builder.WriteString(l.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func lineDiff(oldText, newText string) Result {
	var result Result
	blocks := calcBlockDiff(oldText, newText)
	for _, block := range blocks {
		lines := strings.Split(block.Text, "\n")
		for i, line := range lines[:len(lines)-1] {
			lineObj := Line{
				Ope: block.Ope,
				Fragments: []Fragment{
					{
						Text: line,
					},
				},
				OldLineNumber: -1,
				NewLineNumber: -1,
			}
			if block.NewLineNumber > 0 {
				lineObj.NewLineNumber = block.NewLineNumber + i
			}
			if block.OldLineNumber > 0 {
				lineObj.OldLineNumber = block.OldLineNumber + i
			}
			result.Lines = append(result.Lines, lineObj)
		}
	}
	return result
}

func splitDiffsByNewLine(diffs []diffmatchpatch.Diff) []diffmatchpatch.Diff {
	result := make([]diffmatchpatch.Diff, len(diffs))
	for _, diff := range diffs {
		texts := strings.Split(diff.Text, "\n")
		for i, text := range texts {
			if i != len(texts)-1 {
				text = text + "\n"
			}
			if text != "" {
				result = append(result, diffmatchpatch.Diff{
					Text: text,
					Type: diff.Type,
				})
			}
		}
	}
	return result
}

func wordDiff(oldText, newText string) Result {
	var result Result
	blocks := calcBlockDiff(oldText, newText)
	dmp := diffmatchpatch.New()
	for i := 0; i < len(blocks); i++ {
		if i != len(blocks)-1 && blocks[i].Ope == Delete && blocks[i+1].Ope == Insert {
			oldLineNumber := blocks[i].OldLineNumber
			newLineNumber := blocks[i].NewLineNumber
			diffs := dmp.DiffMain(blocks[i].Text, blocks[i+1].Text, true)
			diffs = splitDiffsByNewLine(diffs)
			var fragments []Fragment
			for _, diff := range diffs {
				hasNewLine := strings.HasSuffix(diff.Text, "\n")
				text := strings.TrimRight(diff.Text, "\n")
				if diff.Type == diffmatchpatch.DiffEqual {
					fragments = append(fragments, Fragment{
						Changed: false,
						Text:    text,
					})
				} else if diff.Type == diffmatchpatch.DiffDelete {
					fragments = append(fragments, Fragment{
						Changed: true,
						Text:    text,
					})
				}
				if hasNewLine {
					result.Lines = append(result.Lines, Line{
						Ope:           Delete,
						NewLineNumber: -1,
						OldLineNumber: oldLineNumber,
						Fragments:     fragments,
					})
					fragments = nil
					oldLineNumber++
				}
			}
			fragments = nil
			for _, diff := range diffs {
				hasNewLine := strings.HasSuffix(diff.Text, "\n")
				text := strings.TrimRight(diff.Text, "\n")
				if diff.Type == diffmatchpatch.DiffEqual {
					fragments = append(fragments, Fragment{
						Changed: false,
						Text:    text,
					})
				} else if diff.Type == diffmatchpatch.DiffInsert {
					fragments = append(fragments, Fragment{
						Changed: true,
						Text:    text,
					})
				}
				if hasNewLine {
					result.Lines = append(result.Lines, Line{
						Ope:           Insert,
						NewLineNumber: newLineNumber,
						OldLineNumber: -1,
						Fragments:     fragments,
					})
					fragments = nil
					newLineNumber++
				}
			}
			i++
		} else {
			block := blocks[i]
			lines := strings.Split(block.Text, "\n")
			for i, line := range lines[:len(lines)-1] {
				lineObj := Line{
					Ope: block.Ope,
					Fragments: []Fragment{
						{
							Text:    line,
							Changed: block.Ope != Keep,
						},
					},
					OldLineNumber: -1,
					NewLineNumber: -1,
				}
				if block.NewLineNumber > 0 {
					lineObj.NewLineNumber = block.NewLineNumber + i
				}
				if block.OldLineNumber > 0 {
					lineObj.OldLineNumber = block.OldLineNumber + i
				}
				result.Lines = append(result.Lines, lineObj)
			}
		}
	}
	return result
}

// Diff calcs diff of text
func Diff(oldText, newText string, diffType DiffType) Result {
	if diffType == LineByLine {
		return lineDiff(oldText, newText)
	}
	return wordDiff(oldText, newText)
}
