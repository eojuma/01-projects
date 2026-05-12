package algorithm

import (
    "push_swap/stack"
)

// SortSmall handles stacks of size 4-5
func SortSmall(a, b *stack.Stack, ins *InstructionSet) {
    if a.Len() <= 1 {
        return
    }
    
    if a.Len() == 2 {
        sortTwoWithIns(a, ins)
        return
    }
    
    if a.Len() == 3 {
        sortThreeWithIns(a, ins)
        return
    }
    
    if a.Len() == 4 {
        sortFourOptimized(a, b, ins)
        return
    }
    
    if a.Len() == 5 {
        sortFiveOptimized(a, b, ins)
        return
    }
}


func sortFourOptimized(a, b *stack.Stack, ins *InstructionSet) {
    minPos := FindMinIndex(a)  // From helpers.go
    
    if minPos == 1 {
        ins.Sa(a)
    } else if minPos == 2 {
        ins.Rra(a)
        ins.Rra(a)
    } else if minPos == 3 {
        ins.Rra(a)
    }
    
    ins.Pb(a, b)
    sortThreeWithIns(a, ins)
    ins.Pa(a, b)
}

func sortFiveOptimized(a, b *stack.Stack, ins *InstructionSet) {
    smallest, secondSmallest := findTwoSmallest(a)
    
    bringToTopWithMinOps(a, smallest, ins)
    ins.Pb(a, b)
    
    bringToTopWithMinOps(a, secondSmallest, ins)
    ins.Pb(a, b)
    
    sortThreeWithIns(a, ins)
    
    if b.Values[0] > b.Values[1] {
        ins.Sb(b)
    }
    ins.Pa(a, b)
    ins.Pa(a, b)
}

func bringToTopWithMinOps(s *stack.Stack, value int, ins *InstructionSet) {
    pos := -1
    for i, v := range s.Values {
        if v == value {
            pos = i
            break
        }
    }
    
    if pos == -1 {
        return
    }
    
    if pos <= s.Len()/2 {
        for i := 0; i < pos; i++ {
            ins.Ra(s)
        }
    } else {
        for i := 0; i < s.Len()-pos; i++ {
            ins.Rra(s)
        }
    }
}

func findTwoSmallest(s *stack.Stack) (int, int) {
    smallest := s.Values[0]
    secondSmallest := s.Values[1]
    
    if secondSmallest < smallest {
        smallest, secondSmallest = secondSmallest, smallest
    }
    
    for i := 2; i < s.Len(); i++ {
        val := s.Values[i]
        if val < smallest {
            secondSmallest = smallest
            smallest = val
        } else if val < secondSmallest {
            secondSmallest = val
        }
    }
    
    return smallest, secondSmallest
}

func sortSmallWithIns(a, b *stack.Stack, ins *InstructionSet) {
    SortSmall(a, b, ins)
}