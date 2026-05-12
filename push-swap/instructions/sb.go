package instructions

import "push_swap/stack"

func Sb(b *stack.Stack) {
	if b.Len() < 2 {
		return
	}

	b.Values[0], b.Values[1] = b.Values[1], b.Values[0]

	// println("sb")
	// fmt.Printf("%v this is stack b", *b)

}