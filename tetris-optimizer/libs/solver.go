package libs

func Solve(tetros []Tetromino) [][]rune {
	size := MinBoardSize(len(tetros))

	for {
		board := NewBoard(size)
		if Backtrack(board, tetros, 0) {
			return board
		}
		size++
	}
}

func Backtrack(board [][]rune, tetros []Tetromino, index int) bool {
	if index == len(tetros) {
		return true
	}

	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board); x++ {
			if CanPlace(board, tetros[index], x, y) {
				Place(board, tetros[index], x, y)
				if Backtrack(board, tetros, index+1) {
					return true
				}
				Remove(board, tetros[index], x, y)
			}
		}
	}
	return false
}

func CanPlace(board [][]rune, t Tetromino, x, y int) bool {
	for _, block := range t.Blocks {
		nx := x + block.X
		ny := y + block.Y

		if ny >= len(board) || nx >= len(board) || board[ny][nx] != '.' {
			return false
		}
	}
	return true
}

func Place(board [][]rune, t Tetromino, x, y int) {
	for _, block := range t.Blocks {
		board[y+block.Y][x+block.X] = t.Letter
	}
}

func Remove(board [][]rune, t Tetromino, x, y int) {
	for _, block := range t.Blocks {
		board[y+block.Y][x+block.X] = '.'
	}
}
