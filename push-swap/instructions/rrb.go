package instructions

import "push_swap/stack"

func Rrb(b *stack.Stack) {
	if b.Len() < 2 {
		return
	}

	last := b.Values[len(b.Values)-1]
	b.Values = append([]int{last}, b.Values[:len(b.Values)-1]...)

	// println("rrb")

}