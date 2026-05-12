package instructions

import "push_swap/stack"

func Rra(a *stack.Stack) {
	if a.Len() < 2 {
		return
	}

	last := a.Values[len(a.Values)-1]
	a.Values = append([]int{last}, a.Values[:len(a.Values)-1]...)

	// fmt.Println("rra")
	// fmt.Printf("This is stack A: %v\n", *a)
}
