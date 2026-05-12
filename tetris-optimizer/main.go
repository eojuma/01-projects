package main

import (
	"fmt"
	"os"
	"tetris-optimizer/libs"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR")
		return
	}
	input := os.Args[1]
	tetrominoes, err := libs.ParseFile(input)
	if err != nil {
		fmt.Println("ERROR")
		return
	}
	board := libs.Solve(tetrominoes)
	libs.PrintBoard(board)
}
