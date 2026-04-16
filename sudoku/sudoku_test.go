package sudoku

import "testing"

func TestSudoku(t *testing.T) {
	p, s := NewSudoku(Easy)
	for i := range p {
		for j := range p[i] {
			if p[i][j] == 0 {
				continue
			}
			if p[i][j] != s[i][j] {
				t.Error("bad sudoku")
			}
		}
	}
}
