package internal

import (
	"os"
	"strings"
)

type Options struct {
	Dirs                                []string
	All, Long, Reverse, Time, Recursive bool
}

func GetInputs() *Options {
	args := os.Args[1:]
	inputs := &Options{}

	for _, v := range args {

		if v == "-" {
			inputs.Dirs = append(inputs.Dirs, v)
			continue
		}
		// If it's a flag (starts with - and isn't just "-")
		if strings.HasPrefix(v, "-") && len(v) > 1 {
			for _, ch := range v[1:] {
				switch ch {
				case 'a': inputs.All = true
				case 'l': inputs.Long = true
				case 'r': inputs.Reverse = true
				case 't': inputs.Time = true
				case 'R': inputs.Recursive = true
				}
			}
			continue
		}

		// Treat EVERYTHING else as a path (including "-")
		inputs.Dirs = append(inputs.Dirs, v)
	}

	if len(inputs.Dirs) == 0 {
		inputs.Dirs = append(inputs.Dirs, ".")
	}
	return inputs
}