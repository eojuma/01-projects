package algorithm

import (
    "push_swap/stack"
)

// InsertBackWithIns inserts elements from B back to A optimally
func InsertBackWithIns(a, b *stack.Stack, ins *InstructionSet) {
    for b.Len() > 0 {
        bestIndex := FindBestElement(a, b)
        
        rotateBToTopWithIns(b, bestIndex, ins)
        
        targetIndex := FindInsertIndex(a, b.Values[0])
        
        rotateToTopWithIns(a, targetIndex, ins)
        
        ins.Pa(a, b)
    }
}

func rotateBToTopWithIns(b *stack.Stack, index int, ins *InstructionSet) {
    if index == 0 {
        return
    }
    
    if index <= b.Len()/2 {
        for i := 0; i < index; i++ {
            ins.Rb(b)
        }
    } else {
        for i := 0; i < b.Len()-index; i++ {
            ins.Rrb(b)
        }
    }
}

func rotateToTopWithIns(a *stack.Stack, index int, ins *InstructionSet) {
    if index == 0 {
        return
    }
    
    if index <= a.Len()/2 {
        for i := 0; i < index; i++ {
            ins.Ra(a)
        }
    } else {
        for i := 0; i < a.Len()-index; i++ {
            ins.Rra(a)
        }
    }
}

func finalRotateWithIns(a *stack.Stack, ins *InstructionSet) {
    minIndex := FindMinIndex(a)
    rotateToTopWithIns(a, minIndex, ins)
}