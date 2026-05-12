package libs

import "fmt"

func NewBoard(size int) [][]rune {
	board := make([][]rune, size)

	for i := range board {
		board[i] = make([]rune, size)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}
	return board
}

func MinBoardSize(n int) int {
	size := 2

	for size*size < n*4 {
		size++
	}
	return size
}

func PrintBoard(board [][]rune) {
	for _, row := range board {
		fmt.Println(string(row))
	}
}
