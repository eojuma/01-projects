package pkg


import "strings"

func FixArticle(text string) string {
	words := strings.Fields(text)

	for i := 0; i < len(words)-1; i++ {
		w := words[i]

		if w == "a" || w == "A" {
			next := strings.ToLower(words[i+1])

			if len(next) > 0 {
				first := next[0]
				if first == 'a' || first == 'e' || first == 'i' ||
					first == 'o' || first == 'u' || first == 'h' {

					if w == "A" {
						words[i] = "An"
					} else {
						words[i] = "an"
					}
				}
			}
		}
	}

	return strings.Join(words, " ")
}
