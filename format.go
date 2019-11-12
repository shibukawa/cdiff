package cdiff

import (
	"strconv"
	"strings"
)

func format(lines []Line, builder *strings.Builder, theme map[Tag]string) {
	for _, l := range lines {
		switch l.Ope {
		case Insert:
			builder.WriteString(theme[OpenInsertedLine])
			builder.WriteString(theme[OpenInsertedNotModified])
			builder.WriteString("+")
			builder.WriteString(theme[CloseInsertedNotModified])
			for _, f := range l.Fragments {
				if f.Changed {
					builder.WriteString(theme[OpenInsertedModified])
					builder.WriteString(f.Text)
					builder.WriteString(theme[CloseInsertedModified])
				} else {
					builder.WriteString(theme[OpenInsertedNotModified])
					builder.WriteString(f.Text)
					builder.WriteString(theme[CloseInsertedNotModified])
				}
			}
			builder.WriteString(theme[CloseInsertedLine])
		case Delete:
			builder.WriteString(theme[OpenDeletedLine])
			builder.WriteString(theme[OpenDeletedNotModified])
			builder.WriteString("-")
			builder.WriteString(theme[CloseDeletedNotModified])
			for _, f := range l.Fragments {
				if f.Changed {
					builder.WriteString(theme[OpenDeletedModified])
					builder.WriteString(f.Text)
					builder.WriteString(theme[CloseDeletedModified])
				} else {
					builder.WriteString(theme[OpenDeletedNotModified])
					builder.WriteString(f.Text)
					builder.WriteString(theme[CloseDeletedNotModified])
				}
			}
			builder.WriteString(theme[CloseDeletedLine])
		case Keep:
			builder.WriteString(theme[OpenKeepLine] + " ")
			for _, f := range l.Fragments {
				builder.WriteString(f.Text)
			}
			builder.WriteString(theme[CloseKeepLine])
		}
	}
}

// Format returns formatted text
func (r Result) Format(theme map[Tag]string) string {
	var builder strings.Builder
	format(r.Lines, &builder, theme)
	return builder.String()
}

type block struct {
	start int
	end   int
}

func (b block) section(lines []Line) string {
	minRemoved := 0
	minInserted := 0
	getValue := func(orig, newValue int) int {
		if orig != 0 {
			return orig
		}
		return newValue
	}
	for i := b.start; i <= b.end; i++ {
		switch lines[i].Ope {
		case Insert:
			minInserted = getValue(minInserted, lines[i].NewLineNumber)
			if minRemoved != 0 {
				break
			}
		case Delete:
			minRemoved = getValue(minRemoved, lines[i].OldLineNumber)
			if minInserted != 0 {
				break
			}
		case Keep:
			minRemoved = getValue(minRemoved, lines[i].OldLineNumber)
			minInserted = getValue(minInserted, lines[i].NewLineNumber)
			break
		}
	}

	maxRemoved := 0
	maxInserted := 0
	for i := b.end; i >= b.start; i-- {
		switch lines[i].Ope {
		case Insert:
			maxInserted = getValue(maxInserted, lines[i].NewLineNumber)
			if maxRemoved != 0 {
				break
			}
		case Delete:
			maxRemoved = getValue(maxRemoved, lines[i].OldLineNumber)
			if minInserted != 0 {
				break
			}
		case Keep:
			maxRemoved = getValue(maxRemoved, lines[i].OldLineNumber)
			maxInserted = getValue(maxInserted, lines[i].NewLineNumber)
			break
		}
	}
	render := func(min, max int) string {
		start := strconv.FormatInt(int64(min), 10)
		if max == min {
			return start
		}
		return start + "," + strconv.FormatInt(int64(max-min+1), 10)
	}
	return "@@ -" + render(minRemoved, maxRemoved) + " +" + render(minInserted, maxInserted) + " @@"
}

func grouping(lines []Line, extraLine int) []block {
	blockStart := -1
	var blocks []block
	for i, l := range lines {
		if blockStart > -1 {
			if l.Ope == Keep {
				blocks = append(blocks, block{start: blockStart - extraLine, end: i - 1 + extraLine})
				blockStart = -1
			}
		} else {
			if l.Ope != Keep {
				blockStart = i
			}
		}
	}
	if blockStart > -1 {
		blocks = append(blocks, block{start: blockStart - extraLine, end: len(lines) - 1 + extraLine})
	}
	if len(blocks) > 0 {
		if blocks[0].start < 0 {
			blocks[0].start = 0
		}
		if blocks[len(blocks)-1].end >= len(lines) {
			blocks[len(blocks)-1].end = len(lines) - 1
		}
	}
	result := make([]block, 0, len(blocks))
	for i, block := range blocks {
		if i == 0 {
			result = append(result, block)
		} else {
			if result[len(result)-1].end >= block.start {
				result[len(result)-1].end = block.end
			} else {
				result = append(result, block)
			}
		}
	}
	return result
}

// Unified returns unified format diff text
func (r Result) Unified(oldTitle, newTitle string, l int, theme map[Tag]string) string {
	var builder strings.Builder
	builder.WriteString(theme[OpenHeader])
	builder.WriteString("--- " + oldTitle)
	builder.WriteString(theme[CloseHeader])
	builder.WriteString(theme[OpenHeader])
	builder.WriteString("+++ " + newTitle)
	builder.WriteString(theme[CloseHeader])
	blocks := grouping(r.Lines, l)
	for _, block := range blocks {
		builder.WriteString(theme[OpenSection])
		builder.WriteString(block.section(r.Lines))
		builder.WriteString(theme[CloseSection])
		lines := r.Lines[block.start : block.end+1]
		format(lines, &builder, theme)
	}
	return builder.String()
}
