package libs

import "strings"

func IsConnected(blocks []Point) bool {
	connections := 0
	for i := range blocks {
		for j := range blocks {
			if i != j {
				dx := Abs(blocks[i].X - blocks[j].X)
				dy := Abs(blocks[i].Y - blocks[j].Y)

				if dx+dy == 1 {
					connections++
				}
			}
		}
	}
	return connections == 6 || connections == 8
}

func IsValidFormat(lines []string) bool {
	if len(lines) != 4 {
		return false
	}

	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != 4 {
			return false
		}

		for _, char := range line {
			if char != '.' && char != '#' {
				return false
			}
			if char == '#' {
				count++
			}
		}
	}
	return count == 4
}
