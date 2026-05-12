package instructions

import "push_swap/stack"

func Ra(a *stack.Stack) {
	if a.Len() < 2 {
		return
	}

	first := a.Values[0]
	a.Values = append(a.Values[1:], first)

	// fmt.Println("ra")
	// fmt.Printf("This is stack A: %v\n", *a)

}
