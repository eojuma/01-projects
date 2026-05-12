package pkg

import (
	"strings"
)

func firstCap(word string) string {
	if len(word) == 0 {
		return ""
	}

	first := word[0]
	if first >= 'a' && first <= 'z' {
		first -= 32
	}

	return string(first) + word[1:]
}


func extractNum(word string) int {
	n := 0
	for i := 0; i < len(word); i++ {
		if word[i] >= '0' && word[i] <= '9' {
			n = n*10 + int(word[i]-'0')
		}
	}
	return n
}

func ApplyCap(text string) string {
	words := strings.Fields(text)

	for i := 0; i < len(words); i++ {

		
		if words[i] == "(cap)" {
			if i-1 >= 0 {
				words[i-1] = firstCap(words[i-1])
			}
			words[i] = ""
		}

	
	if words[i] == "(cap," && i+1 < len(words) {
	n := extractNum(words[i+1])

	for j := i - n; j < i; j++ {
		if j >= 0 {
			words[j] = firstCap(words[j])
		}
	}

	words[i] = ""
	words[i+1] = ""
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


