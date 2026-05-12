package ascii

import "fmt"

func BuildASCII(input []string, charmap map[rune][]string) []string {
	var res []string

	for _, word := range input {
		if word == "" {
			res = append(res, "")
			continue
		}
		for i := 0; i < 8; i++ {
			var line string
			for _, letter := range word {
				line += charmap[letter][i]
			}
			res = append(res, line)
		}
	}
	return res
}

func OutputCharacters(input []string, charmap map[rune][]string) {
	for _, line := range BuildASCII(input, charmap) {
		fmt.Println(line)
	}
}