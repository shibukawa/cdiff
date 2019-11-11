package cdiff

import (
	"reflect"
	"strings"
	"testing"
)

func dumpForTest(r Result) string {
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
		for _, f := range l.Fragments {
			if f.Changed {
				builder.WriteString("[")
				builder.WriteString(f.Text)
				builder.WriteString("]")
			} else {
				builder.WriteString(f.Text)
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()

}

func TestBlockDiff(t *testing.T) {
	type args struct {
		src1 string
		src2 string
	}
	tests := []struct {
		name string
		args args
		want []blockDiff
	}{
		{
			name: "same",
			args: args{
				src1: "abc\ndef\n",
				src2: "abc\ndef\n",
			},
			want: []blockDiff{
				{
					Ope:           Keep,
					Text:          "abc\ndef\n",
					NewLineNumber: 1,
					OldLineNumber: 1,
				},
			},
		},
		{
			name: "diff line",
			args: args{
				src1: "abc\ndef\n",
				src2: "abc\nghi\n",
			},
			want: []blockDiff{
				{
					Ope:           Keep,
					Text:          "abc\n",
					NewLineNumber: 1,
					OldLineNumber: 1,
				},
				{
					Ope:           Delete,
					Text:          "def\n",
					NewLineNumber: -1,
					OldLineNumber: 2,
				},
				{
					Ope:           Insert,
					Text:          "ghi\n",
					NewLineNumber: 2,
					OldLineNumber: -1,
				},
			},
		},
		{
			name: "diff line",
			args: args{
				src1: "abc\ndef\nghi\n",
				src2: "abc\ndef\nghj\n",
			},
			want: []blockDiff{
				{
					Ope:           Keep,
					Text:          "abc\ndef\n",
					NewLineNumber: 1,
					OldLineNumber: 1,
				},
				{
					Ope:           Delete,
					Text:          "ghi\n",
					NewLineNumber: -1,
					OldLineNumber: 3,
				},
				{
					Ope:           Insert,
					Text:          "ghj\n",
					NewLineNumber: 3,
					OldLineNumber: -1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcBlockDiff(tt.args.src1, tt.args.src2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LineDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	type args struct {
		oldText  string
		newText  string
		diffType DiffType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "same",
			args: args{
				oldText:  "abc\ndef\n",
				newText:  "abc\ndef\n",
				diffType: LineByLine,
			},
			want: "  abc\n  def\n",
		},
		{
			name: "remove",
			args: args{
				oldText:  "abc\ndef\n",
				newText:  "abc\n",
				diffType: LineByLine,
			},
			want: "  abc\n- def\n",
		},
		{
			name: "insert",
			args: args{
				oldText:  "abc\n",
				newText:  "abc\ndef\n",
				diffType: LineByLine,
			},
			want: "  abc\n+ def\n",
		},
		{
			name: "diff",
			args: args{
				oldText:  "abc\ndef\n",
				newText:  "abc\ndeg\n",
				diffType: LineByLine,
			},
			want: "  abc\n- def\n+ deg\n",
		},
		{
			name: "github style diff",
			args: args{
				oldText:  "abc\ndef\n",
				newText:  "abc\ndeg\n",
				diffType: WordByWord,
			},
			want: "  abc\n- de[f]\n+ de[g]\n",
		},
		{
			name: "github style diff: multi line blocks",
			args: args{
				oldText:  "abc\ndef\nghi\n",
				newText:  "abc\ndeg\nghj\n",
				diffType: WordByWord,
			},
			want: "  abc\n- de[f]\n- gh[i]\n+ de[g]\n+ gh[j]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.oldText, tt.args.newText, tt.args.diffType); dumpForTest(got) != tt.want {
				for _, line := range calcBlockDiff(tt.want, got.String()) {
					t.Log(line)
				}
				t.Errorf("Diff() = %v, want %v", dumpForTest(got), tt.want)
			}
		})
	}
}
