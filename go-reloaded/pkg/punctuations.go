package pkg

func ParsePunc(s string) string {
	ss := []rune(s)

	i := 0
	for i < len(ss) {
		if ss[i] == '\'' {
			if i+1 < len(ss) && ss[i+1] == ' ' {
				ss = append(ss[:i+1], ss[i+2:]...)
			}
			j := i + 1
			for j < len(ss) && ss[j] != '\'' {
				j++
			}
			if j < len(ss) {
				if j > 0 && ss[j-1] == ' ' {
					ss = append(ss[:j-1], ss[j:]...)
					j--
				}
				i = j
			}
		}
		i++
	}

	i = len(ss) - 1
	for i >= 0 {
		if IsPunc(ss[i]) && i > 0 && ss[i-1] == ' ' {
			ss = append(ss[:i-1], ss[i:]...)
		}
		i--
	}

	i = 0
	for i < len(ss) {
		if IsPunc(ss[i]) && i < len(ss)-1 && ss[i+1] != ' ' {
			if ss[i+1] != '\'' {
				ss = append(ss[:i+1], append([]rune{' '}, ss[i+1:]...)...)
				i++
			}
		}
		i++
	}

	return string(ss)
}

func FormatPunctuation(s string) string {
	ss := []rune(s)

	i := 0
	for i < len(ss) {
		if ss[i] == '\'' {
			if i+1 < len(ss) && ss[i+1] == ' ' {
				ss = append(ss[:i+1], ss[i+2:]...)
			}
			j := i + 1
			for j < len(ss) && ss[j] != '\'' {
				j++
			}
			if j < len(ss) {
				if j > 0 && ss[j-1] == ' ' {
					ss = append(ss[:j-1], ss[j:]...)
					j--
				}
				i = j
			}
		}
		i++
	}

	i = len(ss) - 1
	for i >= 0 {
		if IsPunc(ss[i]) && i > 0 && ss[i-1] == ' ' {
			ss = append(ss[:i-1], ss[i:]...)
		}
		i--
	}

	i = 0
	for i < len(ss) {
		if IsPunc(ss[i]) && i < len(ss)-1 && ss[i+1] != ' ' {
			if ss[i+1] != '\'' {
				ss = append(ss[:i+1], append([]rune{' '}, ss[i+1:]...)...)
				i++
			}
		}
		i++
	}

	return string(ss)
}


func IsPunc(r rune) bool {
	switch r {
	case '.', ',', '!', '?', ':', ';':
		return true
	default:
		return false
	}
}
