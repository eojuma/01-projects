package libs

type Point struct {
	X int
	Y int
}

type Tetromino struct {
	Blocks []Point
	Letter rune
}
