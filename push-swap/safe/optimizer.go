package algorithm



// Combine instructions when possible
func OptimizeInstructions(instructions []string) []string {
    optimized := []string{}
    i := 0
    
    for i < len(instructions) {
        if i < len(instructions)-1 {
            // Check for combinable pairs
            if instructions[i] == "ra" && instructions[i+1] == "rb" {
                optimized = append(optimized, "rr")
                i += 2
                continue
            }
            if instructions[i] == "rra" && instructions[i+1] == "rrb" {
                optimized = append(optimized, "rrr")
                i += 2
                continue
            }
            if instructions[i] == "sa" && instructions[i+1] == "sb" {
                optimized = append(optimized, "ss")
                i += 2
                continue
            }
        }
        optimized = append(optimized, instructions[i])
        i++
    }
    
    return optimized
}

// Remove redundant instructions
func RemoveRedundancies(instructions []string) []string {
    result := []string{}
    
    for i := 0; i < len(instructions); i++ {
        if i < len(instructions)-1 {
            // Remove pairs that cancel each other
            if (instructions[i] == "ra" && instructions[i+1] == "rra") ||
               (instructions[i] == "rra" && instructions[i+1] == "ra") ||
               (instructions[i] == "sa" && instructions[i+1] == "sa") ||
               (instructions[i] == "pb" && instructions[i+1] == "pa") ||
               (instructions[i] == "pa" && instructions[i+1] == "pb") {
                i += 2
                continue
            }
        }
        result = append(result, instructions[i])
    }
    
    return result
}