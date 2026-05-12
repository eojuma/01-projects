package instructions

import "push_swap/stack"

func Rb(b *stack.Stack) {
	if b.Len() < 2 {
		return
	}

	first := b.Values[0]
	b.Values = append(b.Values[1:], first)

	// println("rb")
	// fmt.Printf("This is stack B: %v\n", *b)
}
