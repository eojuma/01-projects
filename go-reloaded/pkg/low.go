package pkg

import (
	"strings"
)

func ToLower(str string) string {
	var result string
	for _, char := range str {
		if char >= 'A' && char <= 'Z' {
			char += 32
		}
		result += string(char)
	}
	return result
}


func ExtractNumb(word string) int {
	n := 0
	for i := 0; i < len(word); i++ {
		if word[i] >= '0' && word[i] <= '9' {
			n = n*10 + int(word[i]-'0')
		}
	}
	return n
}

func ApplyLow(text string) string {
	words := strings.Fields(text)

	for i := 0; i < len(words); i++ {

		
		if words[i] == "(low)" {
			if i-1 >= 0 {
				words[i-1] = ToLower(words[i-1])
			}
			words[i] = ""

		} else if words[i] == "(low," && i+1 < len(words) {
			n := ExtractNumb(words[i+1])
			for j := i - n; j < i; j++ {
				if j >= 0 {
					words[j] = ToLower(words[j])
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

