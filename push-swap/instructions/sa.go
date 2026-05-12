package instructions

import "push_swap/stack"

func Sa(a *stack.Stack) {
	if a.Len() < 2 {
		return
	}

	a.Values[0], a.Values[1] = a.Values[1], a.Values[0]

	// fmt.Println("sa")
	// fmt.Printf("%v this is stack a", *a)

}
