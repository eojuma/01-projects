package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"push_swap/instructions"
	"push_swap/stack"
)

func parseInput(args []string) ([]int, error) {
	input := strings.Join(args, " ")
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

func executeInstruction(cmd string, a, b *stack.Stack) error {
	switch cmd {

	case "sa":
		instructions.Sa(a)
	case "sb":
		instructions.Sb(b)
	case "ss":
		instructions.Sa(a)
		instructions.Sb(b)
	case "pa":
		instructions.Pa(a, b)
	case "pb":
		instructions.Pb(a, b)

	case "ra":
		instructions.Ra(a)
	case "rb":
		instructions.Rb(b)
	case "rr":
		instructions.Ra(a)
		instructions.Rb(b)

	case "rra":
		instructions.Rra(a)
	case "rrb":
		instructions.Rrb(b)
	case "rrr":
		instructions.Rra(a)
		instructions.Rrb(b)

	default:
		return fmt.Errorf("invalid instruction")
	}

	return nil
}



func main() {
	if len(os.Args) < 2 {
		return
	}

	values, err := parseInput(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	a := stack.NewStack(values)
	b := stack.NewStack([]int{})

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if err := executeInstruction(line, a, b); err != nil {
			fmt.Fprintln(os.Stderr, "Error")
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}



if a.IsSorted() && b.Len() == 0 {
	fmt.Println("OK")
	} else {
	fmt.Println("KO")
}
}
