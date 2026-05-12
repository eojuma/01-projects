package algorithm

import (
    "push_swap/stack"
    "testing"
    // "fmt"
)

// Mock instruction set for testing
func newMockInstructionSet() *InstructionSet {
    return &InstructionSet{
        Sa: func(a *stack.Stack) { 
            if a.Len() >= 2 {
                a.Values[0], a.Values[1] = a.Values[1], a.Values[0]
            }
        },
        Sb: func(b *stack.Stack) { 
            if b.Len() >= 2 {
                b.Values[0], b.Values[1] = b.Values[1], b.Values[0]
            }
        },
        Ss: func(a, b *stack.Stack) { 
            if a.Len() >= 2 {
                a.Values[0], a.Values[1] = a.Values[1], a.Values[0]
            }
            if b.Len() >= 2 {
                b.Values[0], b.Values[1] = b.Values[1], b.Values[0]
            }
        },
        Pa: func(a, b *stack.Stack) { 
            if b.Len() > 0 {
                val := b.Values[0]
                b.Values = b.Values[1:]
                a.Values = append([]int{val}, a.Values...)
            }
        },
        Pb: func(a, b *stack.Stack) { 
            if a.Len() > 0 {
                val := a.Values[0]
                a.Values = a.Values[1:]
                b.Values = append([]int{val}, b.Values...)
            }
        },
        Ra: func(a *stack.Stack) { 
            if a.Len() > 1 {
                a.Values = append(a.Values[1:], a.Values[0])
            }
        },
        Rb: func(b *stack.Stack) { 
            if b.Len() > 1 {
                b.Values = append(b.Values[1:], b.Values[0])
            }
        },
        Rr: func(a, b *stack.Stack) { 
            if a.Len() > 1 {
                a.Values = append(a.Values[1:], a.Values[0])
            }
            if b.Len() > 1 {
                b.Values = append(b.Values[1:], b.Values[0])
            }
        },
        Rra: func(a *stack.Stack) { 
            if a.Len() > 1 {
                a.Values = append([]int{a.Values[a.Len()-1]}, a.Values[:a.Len()-1]...)
            }
        },
        Rrb: func(b *stack.Stack) { 
            if b.Len() > 1 {
                b.Values = append([]int{b.Values[b.Len()-1]}, b.Values[:b.Len()-1]...)
            }
        },
        Rrr: func(a, b *stack.Stack) { 
            if a.Len() > 1 {
                a.Values = append([]int{a.Values[a.Len()-1]}, a.Values[:a.Len()-1]...)
            }
            if b.Len() > 1 {
                b.Values = append([]int{b.Values[b.Len()-1]}, b.Values[:b.Len()-1]...)
            }
        },
    }
}

// Helper function to create a stack from slice
func newStack(values []int) *stack.Stack {
    s := &stack.Stack{Values: make([]int, len(values))}
    copy(s.Values, values)
    return s
}

// Helper to check if stack is sorted
func isStackSorted(s *stack.Stack) bool {
    for i := 0; i < s.Len()-1; i++ {
        if s.Values[i] > s.Values[i+1] {
            return false
        }
    }
    return true
}

// ==================== TEST CASES ====================

func TestSortTwo(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        {"Already sorted", []int{1, 2}},
        {"Needs swap", []int{2, 1}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            
            SortWithInstructions(a, b, ins)
            
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v, want sorted", a.Values)
            }
            if b.Len() != 0 {
                t.Errorf("Stack B should be empty, got %v", b.Values)
            }
        })
    }
}

func TestSortThree(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        {"Already sorted", []int{1, 2, 3}},
        {"2 1 3", []int{2, 1, 3}},
        {"3 1 2", []int{3, 1, 2}},
        {"1 3 2", []int{1, 3, 2}},
        {"2 3 1", []int{2, 3, 1}},
        {"3 2 1", []int{3, 2, 1}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            
            SortWithInstructions(a, b, ins)
            
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v, want sorted", a.Values)
            }
            if b.Len() != 0 {
                t.Errorf("Stack B should be empty, got %v", b.Values)
            }
        })
    }
}

func TestSortFour(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        {"Already sorted", []int{1, 2, 3, 4}},
        // {"Reverse sorted", []int{4, 3, 2, 1}},
        // {"2 1 3 4", []int{2, 1, 3, 4}},
        // {"3 2 1 4", []int{3, 2, 1, 4}},
        {"4 1 2 3", []int{4, 1, 2, 3}},
        // {"1 3 2 4", []int{1, 3, 2, 4}},
        {"2 3 4 1", []int{2, 3, 4, 1}},
        // {"1 4 3 2", []int{1, 4, 3, 2}},
        // {"3 1 4 2", []int{3, 1, 4, 2}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            
            SortWithInstructions(a, b, ins)
            
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v, want sorted", a.Values)
            }
            if b.Len() != 0 {
                t.Errorf("Stack B should be empty, got %v", b.Values)
            }
        })
    }
}

func TestSortFive(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        {"Already sorted", []int{1, 2, 3, 4, 5}},
        {"Reverse sorted", []int{5, 4, 3, 2, 1}},
        {"2 1 3 4 5", []int{2, 1, 3, 4, 5}},
        {"3 2 1 4 5", []int{3, 2, 1, 4, 5}},
        {"5 1 2 3 4", []int{5, 1, 2, 3, 4}},
        {"1 3 2 4 5", []int{1, 3, 2, 4, 5}},
        {"2 3 4 5 1", []int{2, 3, 4, 5, 1}},
        {"3 4 5 1 2", []int{3, 4, 5, 1, 2}},
        {"4 5 1 2 3", []int{4, 5, 1, 2, 3}},
        {"1 2 4 3 5", []int{1, 2, 4, 3, 5}},
        {"5 3 4 2 1", []int{5, 3, 4, 2, 1}},
        {"2 5 3 1 4", []int{2, 5, 3, 1, 4}},
        {"3 5 2 4 1", []int{3, 5, 2, 4, 1}},
        {"4 2 5 1 3", []int{4, 2, 5, 1, 3}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            SortWithInstructions(a, b, ins)
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v, want sorted", a.Values)
            }
            if b.Len() != 0 {
                t.Errorf("Stack B should be empty, got %v", b.Values)
            }
        })
    }
}

func TestSortLarge(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        // {"6 numbers random", []int{3, 6, 2, 1, 5, 4}},
        // {"6 numbers reverse", []int{6, 5, 4, 3, 2, 1}},
        {"7 numbers", []int{4, 7, 3, 6, 2, 5, 1}},
        {"8 numbers", []int{8, 7, 6, 5, 4, 3, 2, 1}},
        {"9 numbers", []int{5, 9, 4, 8, 3, 7, 2, 6, 1}},
        {"10 numbers reverse", []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}},
        {"10 numbers random", []int{3, 8, 1, 9, 4, 7, 2, 10, 5, 6}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            
            SortWithInstructions(a, b, ins)
            
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v", a.Values)
            }
            if b.Len() != 0 {
                t.Errorf("Stack B should be empty, got %v", b.Values)
            }
        })
    }
}

// Test edge cases
func TestSortEdgeCases(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        {"Single element", []int{42}},
        {"Two elements sorted", []int{1, 2}},
        {"Two elements unsorted", []int{2, 1}},
        {"Already sorted 5", []int{1, 2, 3, 4, 5}},
        {"Already sorted 10", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            
            SortWithInstructions(a, b, ins)
            
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v", a.Values)
            }
        })
    }
}

// Test specific problematic cases
func TestSortProblematicCases(t *testing.T) {
    tests := []struct {
        name  string
        input []int
    }{
        {"1 3 2 7 5", []int{1, 3, 2, 7, 5}},
        {"1 3 5 7 2", []int{1, 3, 5, 7, 2}},
        {"2 1 3 6 5 4", []int{2, 1, 3, 6, 5, 4}},
        {"3 2 1 6 5 4", []int{3, 2, 1, 6, 5, 4}},
        {"4 3 2 1 8 7 6 5", []int{4, 3, 2, 1, 8, 7, 6, 5}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            a := newStack(tt.input)
            b := &stack.Stack{Values: []int{}}
            ins := newMockInstructionSet()
            
            SortWithInstructions(a, b, ins)
            
            if !isStackSorted(a) {
                t.Errorf("Stack not sorted: got %v", a.Values)
            }
            if b.Len() != 0 {
                t.Errorf("Stack B should be empty, got %v", b.Values)
            }
        })
    }
}