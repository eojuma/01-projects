package instructions

import "push_swap/stack" 

func Pa(a, b *stack.Stack) {
    if val, ok := b.Pop(); ok {
        a.Push(val)
    }
    // fmt.Println("pa") 
	// fmt.Printf("This is stack A: %v\n", *a)
}
