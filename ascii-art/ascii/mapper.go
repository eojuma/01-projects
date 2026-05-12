package ascii

func MapToASCII(lines []string) map[rune][]string {
	charmap := make(map[rune][]string)
	start := 1

	for char := ' '; char <= '~'; char++ {
		charmap[char] = lines[start : start+8]
		start += 9
	}

	return charmap
}