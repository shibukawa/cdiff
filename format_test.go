package cdiff

import (
	"testing"

	"github.com/bmizerany/assert"
)

func Test_grouping(t *testing.T) {
	type args struct {
		lines []Line
		l     int
	}
	tests := []struct {
		name string
		args args
		want []block
	}{
		{
			name: "no extra lines (1)",
			args: args{
				lines: []Line{
					{Ope: Insert, NewLineNumber: 1, OldLineNumber: -1},
					{Ope: Keep, NewLineNumber: 2, OldLineNumber: 1},
					{Ope: Delete, NewLineNumber: -1, OldLineNumber: 2},
				},
				l: 0,
			},
			want: []block{
				{start: 0, end: 0},
				{start: 2, end: 2},
			},
		},
		{
			name: "no extra lines (2)",
			args: args{
				lines: []Line{
					{Ope: Keep, NewLineNumber: 1, OldLineNumber: 1},
					{Ope: Insert, NewLineNumber: 2, OldLineNumber: -1},
					{Ope: Keep, NewLineNumber: 3, OldLineNumber: 2},
					{Ope: Delete, NewLineNumber: -1, OldLineNumber: 3},
					{Ope: Keep, NewLineNumber: 4, OldLineNumber: 4},
				},
				l: 0,
			},
			want: []block{
				{start: 1, end: 1},
				{start: 3, end: 3},
			},
		},
		{
			name: "extra lines no merge (1)",
			args: args{
				lines: []Line{
					{Ope: Keep, NewLineNumber: 1, OldLineNumber: 1},
					{Ope: Insert, NewLineNumber: 2, OldLineNumber: -1},
					{Ope: Keep, NewLineNumber: 3, OldLineNumber: 2},
					{Ope: Keep, NewLineNumber: 4, OldLineNumber: 3},
					{Ope: Keep, NewLineNumber: 5, OldLineNumber: 4},
					{Ope: Delete, NewLineNumber: -1, OldLineNumber: 5},
					{Ope: Keep, NewLineNumber: 6, OldLineNumber: 6},
				},
				l: 1,
			},
			want: []block{
				{start: 0, end: 2},
				{start: 4, end: 6},
			},
		},
		{
			name: "extra lines no merge (2)",
			args: args{
				lines: []Line{
					{Ope: Keep, NewLineNumber: 1, OldLineNumber: 1},
					{Ope: Insert, NewLineNumber: 2, OldLineNumber: -1},
					{Ope: Keep, NewLineNumber: 3, OldLineNumber: 2},
				},
				l: 2,
			},
			want: []block{
				{start: 0, end: 2},
			},
		},
		{
			name: "extra lines merge (1)",
			args: args{
				lines: []Line{
					{Ope: Keep, NewLineNumber: 1, OldLineNumber: 1},
					{Ope: Insert, NewLineNumber: 2, OldLineNumber: -1},
					{Ope: Keep, NewLineNumber: 3, OldLineNumber: 2},
					{Ope: Keep, NewLineNumber: 4, OldLineNumber: 3},
					{Ope: Keep, NewLineNumber: 5, OldLineNumber: 4},
					{Ope: Delete, NewLineNumber: -1, OldLineNumber: 5},
					{Ope: Keep, NewLineNumber: 6, OldLineNumber: 6},
				},
				l: 2,
			},
			want: []block{
				{start: 0, end: 6},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := grouping(tt.args.lines, tt.args.l)
			assert.Equal(t, got, tt.want)
		})
	}
}
