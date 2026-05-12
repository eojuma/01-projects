package ascii

import (
	"fmt"
	"strings"
)

func Generate(input string, banner string) ([]string, error) {

	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("input is empty")
	}
	input = strings.ReplaceAll(input, `\n`, "\n")

	bannerFile := banner + ".txt"

	lines, err := ReadLinesFromFile(bannerFile)
	if err != nil {
		return nil, fmt.Errorf("could not read banner file '%s': %v", bannerFile, err)
	}

	charmap := MapToASCII(lines)

	inputLines := strings.Split(input, "\n")

	result := BuildASCII(inputLines, charmap)

	return result, nil
}
