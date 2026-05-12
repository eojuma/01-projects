
---

# Tetris Optimizer

## Project Overview

The **Tetris Optimizer** is a Go program that reads a list of tetrominoes from a text file and assembles them to create the **smallest possible square** on a grid. Each tetromino is identified by an uppercase Latin letter (A, B, C, ...), and the program ensures that:

* Each tetromino is placed without overlapping others.
* All tetrominoes are connected in their valid form.
* The output board uses the minimal square size necessary to fit all pieces.
* Errors in input format or invalid tetrominoes are detected and reported.

This project is designed to teach:

* Parsing input files in Go.
* Validating structured shapes (tetrominoes).
* Using recursive backtracking algorithms.
* Efficient board management and printing.
* Algorithmic thinking for optimization problems.

---

## Features

* Reads a `.txt` file containing tetrominoes.
* Validates each tetromino for:

  * Proper 4x4 format.
  * Exactly 4 `#` blocks.
  * Connectivity (tetromino blocks must be adjacent horizontally or vertically).
* Assembles tetrominoes in order to form the **smallest possible square**.
* Prints the final board to the terminal using uppercase letters.
* Gracefully handles invalid input by printing `ERROR`.

---

## File Structure

```
tetris-optimizer/
├── board.go       # Functions to create, print, and manage the board
├── main.go        # Entry point of the program
├── parser.go      # File reading and tetromino parsing
├── solver.go      # Backtracking algorithm to place tetrominoes
├── type.go        # Data structures for tetrominoes and points
├── utils.go       # Utility functions (abs, normalize coordinates)
├── validator.go   # Tetromino format and connectivity validation
├── sample.txt     # Example input file with tetrominoes
├── go.mod         # Go module file
```

---

## Tetromino File Format

* Each tetromino occupies **4 lines**, each with 4 characters (`#` or `.`).
* Tetrominoes are separated by **one empty line**.
* Example (`sample.txt`):

```
...#
...#
...#
...#

....
....
..##
..##

.###
...#
....
....
```

* `#` represents a block, `.` represents an empty space.
* Each tetromino must have **exactly 4 `#` blocks** and all blocks must be connected.

---

## How to Run

1. Make sure you have **Go installed** (version >= 1.24).
2. Clone or navigate to the project directory.
3. Use the command:

```bash
go run . sample.txt
```

4. The output will be the smallest square board with tetrominoes placed. Example output:

```
ABB.C
ABB.C
AEEEC
ADDEC
.DD..

```

5. If the file has invalid tetrominoes or formatting issues, the program prints:

```
ERROR
```

---

## Key Functions Overview

### `board.go`

* `NewBoard(size int) [][]rune` → Creates an empty square board of given size.
* `MinBoardSize(n int) int` → Calculates the minimal board size to fit `n` tetrominoes.
* `PrintBoard(board [][]rune)` → Prints the board to the terminal.

### `parser.go`

* `ParseFile(path string) ([]Tetromino, error)` → Reads tetrominoes from a file.
* `ParseBlock(lines []string, letter rune)` → Converts a 4x4 block into a `Tetromino` struct.

### `validator.go`

* `IsConnected(blocks []Point) bool` → Checks if a tetromino’s blocks are adjacent.
* `IsValidFormat(lines []string) bool` → Validates tetromino format (4x4, exactly 4 `#`).

### `solver.go`

* `Solve(tetros []Tetromino) [][]rune` → Places tetrominoes using recursive backtracking.
* `Backtrack(board [][]rune, tetros []Tetromino, index int) bool` → Core recursive solver.
* `CanPlace`, `Place`, `Remove` → Helper functions to manipulate tetrominoes on the board.

### `type.go`

* `Point` → Struct representing x, y coordinates.
* `Tetromino` → Struct representing a tetromino and its identifying letter.

### `utils.go`

* `abs(int) int` → Absolute value helper.
* `Normalize([]Point) []Point` → Normalizes coordinates to start at (0,0).

---

## Error Handling

The program will output `ERROR` if:

* The file is missing or cannot be read.
* A tetromino is not 4x4.
* A tetromino contains invalid characters.
* A tetromino has less or more than 4 `#` blocks.
* Tetromino blocks are disconnected.

---

## Example Input & Output

**Input (`sample.txt`):**

```
...#
...#
...#
...#

....
....
..##
..##
```

**Output:**

```
ABB.C
ABB.C
AEEEC
ADDEC
.DD..
```

---

## Notes

* The program uses **recursive backtracking**, which ensures it will find the **smallest possible square**.
* Tetrominoes are placed **in order**, and each is represented by a **unique uppercase letter**.
* The project follows **clean Go practices**, separating parsing, validation, board management, and solving logic into separate files.

---

## Author
* Written and implemented by **Evans Juma**
* Designed to meet the Tetris Optimizer project requirements in Go.

---

