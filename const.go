package cdiff

import "fmt"

// DiffType is an option of Diff()
type DiffType int

const (
	// LineByLine returns simple result
	LineByLine DiffType = iota + 1
	// WordByWord returns complex result that includes word by word diffs in line by line diffs
	WordByWord
)

// Ope defines the operation of a diff item.
type Ope int8

func (o Ope) String() string {
	switch o {
	case Delete:
		return "@Delete"
	case Insert:
		return "@Insert"
	case Keep:
		return "@Keep"
	}
	panic(fmt.Sprintf("unknown value of Ope: %d", int(o)))
}

const (
	// Delete item represents a delete diff.
	Delete Ope = -1
	// Insert item represents an insert diff.
	Insert Ope = 1
	// Keep item represents an equal diff.
	Keep Ope = 0
)

// Tag is for formatting text
type Tag int

const (
	OpenDeletedLine Tag = iota + 1
	CloseDeletedLine
	OpenDeletedModified
	CloseDeletedModified
	OpenDeletedNotModified
	CloseDeletedNotModified
	OpenInsertedLine
	CloseInsertedLine
	OpenInsertedModified
	CloseInsertedModified
	OpenInsertedNotModified
	CloseInsertedNotModified
	OpenKeepLine
	CloseKeepLine
	OpenSection
	CloseSection
	OpenHeader
	CloseHeader
)

// GooKitColorTheme is a theme for Result.Format() method for coloring console
var GooKitColorTheme = map[Tag]string{
	OpenDeletedLine:          "",
	CloseDeletedLine:         "\n",
	OpenDeletedModified:      "<fg=black;bg=red;>",
	CloseDeletedModified:     "</>",
	OpenDeletedNotModified:   "<red>",
	CloseDeletedNotModified:  "</>",
	OpenInsertedLine:         "",
	CloseInsertedLine:        "\n",
	OpenInsertedModified:     "<fg=black;bg=green;>",
	CloseInsertedModified:    "</>",
	OpenInsertedNotModified:  "<green>",
	CloseInsertedNotModified: "</>",
	OpenKeepLine:             "",
	CloseKeepLine:            "\n",
	OpenSection:              "<cyan>",
	CloseSection:             "</>\n",
	OpenHeader:               "",
	CloseHeader:              "\n",
}

// HTMLTheme is a theme for Result.Format() method for generating HTML
var HTMLTheme = map[Tag]string{
	OpenDeletedLine:          `<div style="background-color: #ffecec;">`,
	CloseDeletedLine:         `</div>`,
	OpenDeletedModified:      `<span style="background-color: #f8cbcb;">`,
	CloseDeletedModified:     `</span>`,
	OpenDeletedNotModified:   "",
	CloseDeletedNotModified:  "",
	OpenInsertedLine:         `<div style="background-color: #eaffea;">`,
	CloseInsertedLine:        `</div>`,
	OpenInsertedModified:     `<span style="background-color: #a6f3a6;">`,
	CloseInsertedModified:    `</span>`,
	OpenInsertedNotModified:  "",
	CloseInsertedNotModified: "",
	OpenKeepLine:             `<div style="background-color: #ffffff;">`,
	CloseKeepLine:            "</div>",
}
