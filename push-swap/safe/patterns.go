package algorithm

import (
    "push_swap/stack"
)

// DetectAndSolvePattern checks for common patterns and solves them
// func DetectAndSolvePattern(a, b *stack.Stack, ins *InstructionSet) bool {
//     // Check if stack is already sorted
//     if IsSortedAscending(a) {
//         return true
//     }
    
//     // Check if only one rotation is needed
//     if IsOneRotationNeeded(a) {
//         SolveWithOneRotation(a, ins)
//         return true
//     }
    
//     // Check if only a swap is needed
//     if a.Len() == 2 && a.Values[0] > a.Values[1] {
//         ins.Sa(a)
//         return true
//     }
    
//     // Check for reverse sorted
//     if IsReverseSorted(a) {
//         SolveReverseSorted(a, ins)
//         return true
//     }
    
//     return false
// }

// IsSortedAscending checks if the stack is sorted in ascending order
func IsSortedAscending(s *stack.Stack) bool {
    for i := 0; i < s.Len()-1; i++ {
        if s.Values[i] > s.Values[i+1] {
            return false
        }
    }
    return true
}

// IsOneRotationNeeded checks if one rotation would sort the stack
func IsOneRotationNeeded(s *stack.Stack) bool {
    if s.Len() < 2 {
        return false
    }
    
    // Try rotating forward
    rotated := append(s.Values[1:], s.Values[0])
    sorted := true
    for i := 0; i < len(rotated)-1; i++ {
        if rotated[i] > rotated[i+1] {
            sorted = false
            break
        }
    }
    if sorted {
        return true
    }
    
    // Try rotating backward
    rotated = append([]int{s.Values[s.Len()-1]}, s.Values[:s.Len()-1]...)
    for i := 0; i < len(rotated)-1; i++ {
        if rotated[i] > rotated[i+1] {
            return false
        }
    }
    return true
}

// SolveWithOneRotation applies the minimal rotation to sort
func SolveWithOneRotation(a *stack.Stack, ins *InstructionSet) {
    if a.Len() < 2 {
        return
    }
    
    minPos := FindMinIndex(a)  // From helpers.go
    
    if minPos == 1 {
        ins.Ra(a)
    } else if minPos == a.Len()-1 {
        ins.Rra(a)
    }
}

// IsReverseSorted checks if the stack is in reverse order
func IsReverseSorted(s *stack.Stack) bool {
    for i := 0; i < s.Len()-1; i++ {
        if s.Values[i] < s.Values[i+1] {
            return false
        }
    }
    return true
}

// SolveReverseSorted handles reverse-sorted stacks
func SolveReverseSorted(a *stack.Stack, ins *InstructionSet) {
    if a.Len() == 3 {
        ins.Sa(a)
        ins.Rra(a)
    } else if a.Len() == 5 {
        ins.Rra(a)
        ins.Rra(a)
        ins.Sa(a)
        ins.Ra(a)
        ins.Ra(a)
        ins.Sa(a)
    }
}