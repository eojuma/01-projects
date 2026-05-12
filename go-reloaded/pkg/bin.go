package pkg

import (
	"strconv"
	"strings"
)

func BinToDecimal(str string) int {
	result := 0
	for _, char := range str {
	
		if char =='0' || char == '1' {
			
			result = result*2 
			 if char == '1' {
            result += 1
        } else if char != '0' {
            return -1 
        }
		}
	}
	return result
}

func Final(text string) string {
	words := strings.Fields(text)

	for i := 1; i < len(words); i++ {
		if words[i] == "(bin)" {
			dec := BinToDecimal(words[i-1])
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
