package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"

    "push_swap/algorithm"
    "push_swap/stack"
    "push_swap/instructions"
)

var instructionsList []string

// Wrapper functions that both execute and record
func sa(a *stack.Stack) {
    instructions.Sa(a)
    instructionsList = append(instructionsList, "sa")
}

func sb(b *stack.Stack) {
    instructions.Sb(b)
    instructionsList = append(instructionsList, "sb")
}

func ss(a *stack.Stack, b *stack.Stack) {
    instructions.Sa(a)
    instructions.Sb(b)
    instructionsList = append(instructionsList, "ss")
}

func pa(a *stack.Stack, b *stack.Stack) {
    instructions.Pa(a, b)
    instructionsList = append(instructionsList, "pa")
}

func pb(a *stack.Stack, b *stack.Stack) {
    instructions.Pb(a, b)
    instructionsList = append(instructionsList, "pb")
}

func ra(a *stack.Stack) {
    instructions.Ra(a)
    instructionsList = append(instructionsList, "ra")
}

func rb(b *stack.Stack) {
    instructions.Rb(b)
    instructionsList = append(instructionsList, "rb")
}

func rr(a *stack.Stack, b *stack.Stack) {
    instructions.Ra(a)
    instructions.Rb(b)
    instructionsList = append(instructionsList, "rr")
}

func rra(a *stack.Stack) {
    instructions.Rra(a)
    instructionsList = append(instructionsList, "rra")
}

func rrb(b *stack.Stack) {
    instructions.Rrb(b)
    instructionsList = append(instructionsList, "rrb")
}

func rrr(a *stack.Stack, b *stack.Stack) {
    instructions.Rra(a)
    instructions.Rrb(b)
    instructionsList = append(instructionsList, "rrr")
}

func parseInput(input string) ([]int, error) {
    fields := strings.Fields(input)
    values := []int{}
    seen := make(map[int]bool)

    for _, f := range fields {
        num, err := strconv.Atoi(f)
        if err != nil {
            return nil, err
        }
        if seen[num] {
            return nil, fmt.Errorf("duplicate")
        }
        seen[num] = true
        values = append(values, num)
    }
    return values, nil
}

func main() {
    if len(os.Args) < 2 {
        return
    }

    input := strings.Join(os.Args[1:], " ")
    values, err := parseInput(input)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error")
        return
    }

    // Check if already sorted
    if isSorted(values) {
        return
    }

    a := stack.NewStack(values)
    b := stack.NewStack([]int{})
    
    // Clear instructions list
    instructionsList = []string{}
    
    // Create instruction set
    ins := &algorithm.InstructionSet{
        Sa:  sa,
        Sb:  sb,
        Ss:  ss,
        Pa:  pa,
        Pb:  pb,
        Ra:  ra,
        Rb:  rb,
        Rr:  rr,
        Rra: rra,
        Rrb: rrb,
        Rrr: rrr,
    }
    
    // Run sorting algorithm with instruction set
    algorithm.SortWithInstructions(a, b, ins)
    
    // Print instructions (already recorded)
    for _, instr := range instructionsList {
        fmt.Println(instr)
    }
    // fmt.Println(*b)

    
}

// Helper to check if slice is sorted
func isSorted(values []int) bool {
    for i := 0; i < len(values)-1; i++ {
        if values[i] > values[i+1] {
            return false
        }
    }
    return true
}