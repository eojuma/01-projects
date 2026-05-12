package libs

func Abs(X int) int {
	if X < 0 {
		return -X
	}

	return X
}

func Normalize(blocks []Point) []Point {
	minX := blocks[0].X
	minY := blocks[0].Y

	for _, b := range blocks {
		if b.X < minX {
			minX = b.X
		}
		if b.Y < minY {
			minY = b.Y
		}

	}

	for i := range blocks {
		blocks[i].X -= minX
		blocks[i].Y -= minY
	}
	return blocks
}
