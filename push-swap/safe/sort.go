package algorithm

import (
    "push_swap/stack"
)

// InstructionSet holds all the instruction functions
type InstructionSet struct {
    Sa  func(*stack.Stack)
    Sb  func(*stack.Stack)
    Ss  func(*stack.Stack, *stack.Stack)
    Pa  func(*stack.Stack, *stack.Stack)
    Pb  func(*stack.Stack, *stack.Stack)
    Ra  func(*stack.Stack)
    Rb  func(*stack.Stack)
    Rr  func(*stack.Stack, *stack.Stack)
    Rra func(*stack.Stack)
    Rrb func(*stack.Stack)
    Rrr func(*stack.Stack, *stack.Stack)
}

// SortWithInstructions uses the provided instruction set
func SortWithInstructions(a, b *stack.Stack, ins *InstructionSet) {
    if a.Len() <= 1 {
        return
    }
    
    // First check if we can solve with simple patterns
    if DetectAndSolvePattern(a, b, ins) {
        return
    }
    
    // If not, use the regular algorithm
    if a.Len() == 2 {
        sortTwoWithIns(a, ins)
        return
    }
    
    if a.Len() == 3 {
        sortThreeWithIns(a, ins)
        return
    }
    
    if a.Len() <= 5 {
        SortSmall(a, b, ins)
        return
    }
    
    // For larger stacks
    sortLargeWithIns(a, b, ins)
}

func sortTwoWithIns(a *stack.Stack, ins *InstructionSet) {
    if a.Values[0] > a.Values[1] {
        ins.Sa(a)
    }
}

func sortThreeWithIns(a *stack.Stack, ins *InstructionSet) {
    v := a.Values
    if v[0] < v[1] && v[1] < v[2] {
        return
    }
    if v[0] > v[1] && v[1] > v[2] {
        ins.Sa(a)
        ins.Rra(a)
    } else if v[0] > v[1] && v[0] < v[2] {
        ins.Sa(a)
    } else if v[0] < v[1] && v[0] > v[2] {
        ins.Rra(a)
    } else if v[0] > v[2] && v[1] < v[2] {
        ins.Ra(a)
    } else if v[0] < v[2] && v[1] > v[2] {
        ins.Sa(a)
        ins.Ra(a)
    }
}

func sortLargeWithIns(a, b *stack.Stack, ins *InstructionSet) {
    // Push to B with strategy
    pushToBWithIns(a, b, ins)
    
    // Insert back with cost calculation
    insertBackWithIns(a, b, ins)
    
    // Final rotation
    finalRotateWithIns(a, ins)
}

func pushToBWithIns(a, b *stack.Stack, ins *InstructionSet) {
    median := CalculateMedian(a)
    
    for a.Len() > 3 {
        if a.Values[0] > median {
            ins.Pb(a, b)
            
            if b.Len() > 1 && b.Values[0] < b.Values[1] {
                ins.Rb(b)
            }
        } else {
            ins.Ra(a)
        }
        
        if a.Len()%10 == 0 {
            median = CalculateMedian(a)
        }
    }
    
    sortThreeWithIns(a, ins)
}

func insertBackWithIns(a, b *stack.Stack, ins *InstructionSet) {
    InsertBackWithIns(a, b, ins)
}