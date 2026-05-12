package instructions

import "push_swap/stack" 


func Pb(a, b *stack.Stack) {
    val, ok := a.Pop()
    if !ok { return }
    b.Push(val)
    // fmt.Println("pb") // optional
}

