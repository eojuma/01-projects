package pkg

import (
	"strconv"
	"strings"
)

func HexToDecimal(str string) int {
	result := 0
	for _, char := range str {
		var value int
		if char >= '0' && char <= '9' {
			value = int(char - '0')
		} else if char >= 'A' && char <= 'F' {
			value = int(char-'A') + 10
		} else if char >= 'a' && char <= 'f' {
			value = int(char-'a') + 10
		} else {
			return -1 
		}
		result = result*16 + value
	}
	return result
}

func Finally(text string) string {
	words := strings.Fields(text)

	for i := 1; i < len(words); i++ {
		if words[i] == "(hex)" {
			dec := HexToDecimal(words[i-1])
			if dec != -1 {
				words[i-1] = strconv.Itoa(dec)
			}
			words[i] = "" 
		}
	}


	var res []string
	for _, w := range words {
		if w != "" {
			res = append(res, w)
		}
	}

	return strings.Join(res, " ")
}
