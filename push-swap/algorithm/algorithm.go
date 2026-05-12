package algorithm

import (
    "push_swap/stack"
    "sort"
    "math"
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

// ==================== MAIN SORT FUNCTION ====================

func SortWithInstructions(a, b *stack.Stack, ins *InstructionSet) {
    if a.Len() <= 1 {
        return
    }
    
    // Check if already sorted
    if isSorted(a) {
        return
    }
    
    // Use circular sort for small sets (≤6)
    if a.Len() <= 6 {
        circularSort(a, b, ins)
        return
    }
    
    // Use optimized algorithm for large sets
    sortLarge(a, b, ins)
}

// ==================== CIRCULAR SORT FOR SMALL SETS (≤6) ====================

func circularSort(a, b *stack.Stack, ins *InstructionSet) {
    switch a.Len() {
    case 2:
        if a.Values[0] > a.Values[1] {
            ins.Sa(a)
        }
    case 3:
        sortThree(a, ins)
    case 4:
        sortFour(a, ins)
    case 5:
        sortFive(a, b, ins)
    case 6:
        sortSix(a, b, ins)
    }
}
func sortFour(a *stack.Stack, ins *InstructionSet) {
    // Find the smallest element
    minIdx := findMinIndex(a)
    
    // Bring smallest to top
    switch minIdx {
    case 1:
        ins.Ra(a)
    case 2:
        ins.Ra(a)
        ins.Ra(a)
    case 3:
        ins.Rra(a)
    }
    
    // Now we have [smallest, x, y, z]
    // Sort the remaining 3 using sortThree logic
    
    // Check if the remaining 3 are sorted
    if a.Values[1] > a.Values[2] || a.Values[2] > a.Values[3] {
        // Apply sortThree patterns to the remaining 3
        // We need to map sortThree operations to the full stack
        
        // Create a temporary slice of the remaining 3
        remaining := []int{a.Values[1], a.Values[2], a.Values[3]}
        
        // Apply the correct pattern based on the order
        if remaining[0] > remaining[1] && remaining[0] < remaining[2] {
            // Pattern: [2,1,3]
            ins.Sa(a)
        } else if remaining[0] > remaining[1] && remaining[1] > remaining[2] {
            // Pattern: [3,2,1]
            ins.Sa(a)
            ins.Rra(a)
        } else if remaining[0] > remaining[2] && remaining[1] < remaining[2] {
            // Pattern: [3,1,2]
            ins.Ra(a)
        } else if remaining[0] < remaining[2] && remaining[1] > remaining[2] {
            // Pattern: [1,3,2]
            ins.Sa(a)
            ins.Ra(a)
        } else if remaining[0] < remaining[1] && remaining[0] > remaining[2] {
            // Pattern: [2,3,1]
            ins.Rra(a)
        }
    }
}

func sortFive(a, b *stack.Stack, ins *InstructionSet) {
    // Push 2 to B
    ins.Pb(a, b)
    ins.Pb(a, b)
    
    // Sort remaining 3
    sortThree(a, ins)
    
    // Insert back from B
    insertBackFromB(a, b, ins)
}

func sortSix(a, b *stack.Stack, ins *InstructionSet) {
    // Save the original state for debugging
    
    // Push all but 3 to B
    for a.Len() > 3 {
        ins.Pb(a, b)
    }
    
    // Sort the remaining 3
    sortThree(a, ins)
    
    // Insert back from B - but do it intelligently
    // First, collect all B values
    bValues := make([]int, b.Len())
    for i := 0; i < len(bValues); i++ {
        bValues[i] = b.Values[i]
    }
    
    // Insert in descending order (largest first)
    // Find the largest in B
    for len(bValues) > 0 {
        // Find largest value in remaining B values
        largestIdx := 0
        largestVal := bValues[0]
        for i, val := range bValues {
            if val > largestVal {
                largestVal = val
                largestIdx = i
            }
        }
        
        // Rotate B to bring largest to top
        rotateBToTop(b, largestIdx, ins)
        
        // Insert it
        x := b.Values[0]
        index := findInsertIndexForCircular(a, x)
        rotateAToTop(a, index, ins)
        ins.Pa(a, b)
        
        // Remove from our slice
        bValues = append(bValues[:largestIdx], bValues[largestIdx+1:]...)
    }
}


func insertBackFromB(a, b *stack.Stack, ins *InstructionSet) {
    // First, optimize B order if needed
    if b.Len() == 3 {
        // We want to insert largest first, then middle, then smallest
        // Find the largest in B and bring it to top
        largestIdx := findLargestIndex(b)
        rotateBToTop(b, largestIdx, ins)
    }
    
    for b.Len() > 0 {
        x := b.Values[0]
        index := findInsertIndexForCircular(a, x)
        rotateAToTop(a, index, ins)
        ins.Pa(a, b)
    }
    
    // Final rotation if needed
    if !isSorted(a) {
        minIdx := findMinIndex(a)
        if minIdx != 0 {
            rotateAToTop(a, minIdx, ins)
        }
    }
}

func findLargestIndex(b *stack.Stack) int {
    if b.Len() == 0 {
        return 0
    }
    maxVal := b.Values[0]
    maxIdx := 0
    for i, v := range b.Values {
        if v > maxVal {
            maxVal = v
            maxIdx = i
        }
    }
    return maxIdx
}
func findInsertIndexForCircular(a *stack.Stack, x int) int {
    if a.Len() == 0 {
        return 0
    }
    
    bestIndex := -1
    bestValue := math.MaxInt32
    
    // Find the smallest value greater than x
    for i, value := range a.Values {
        if value > x && value < bestValue {
            bestValue = value
            bestIndex = i
        }
    }
    
    if bestIndex == -1 {
        // x is biggest or smallest - put at position of smallest element
        return findMinIndex(a)
    }
    
    return bestIndex
}

// ==================== LARGE SET SORT ====================

func sortLarge(a, b *stack.Stack, ins *InstructionSet) {
    pushToB(a, b, ins)
    insertBack(a, b, ins)
    finalRotate(a, ins)
}

func pushToB(a, b *stack.Stack, ins *InstructionSet) {
    median := calculateMedian(a)
    rotated := 0
    maxIterations := a.Len() * 10
    iterations := 0
    
    for a.Len() > 3 && iterations < maxIterations {
        iterations++
        
        if a.Values[0] > median {
            ins.Pb(a, b)
            rotated = 0
        } else {
            ins.Ra(a)
            rotated++
            
            if rotated >= a.Len() {
                median = calculateMedian(a)
                rotated = 0
            }
        }
        
        if a.Len()%10 == 0 && a.Len() > 0 {
            median = calculateMedian(a)
            rotated = 0
        }
    }
    
    if iterations >= maxIterations {
        for a.Len() > 3 {
            ins.Pb(a, b)
        }
    }
    
    sortThree(a, ins)
}

func insertBack(a, b *stack.Stack, ins *InstructionSet) {
    for b.Len() > 0 {
        bestIdx := findBestElement(a, b)
        rotateBToTop(b, bestIdx, ins)
        targetIdx := findInsertIndex(a, b.Values[0])
        rotateAToTop(a, targetIdx, ins)
        ins.Pa(a, b)
    }
}

// ==================== COMMON SORTING FUNCTIONS ====================

func sortThree(a *stack.Stack, ins *InstructionSet) {
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

// ==================== ROTATION HELPERS ====================

func rotateAToTop(a *stack.Stack, idx int, ins *InstructionSet) {
    if idx == 0 || a.Len() == 0 {
        return
    }
    if idx <= a.Len()/2 {
        for i := 0; i < idx; i++ {
            ins.Ra(a)
        }
    } else {
        for i := 0; i < a.Len()-idx; i++ {
            ins.Rra(a)
        }
    }
}

func rotateBToTop(b *stack.Stack, idx int, ins *InstructionSet) {
    if idx == 0 || b.Len() == 0 {
        return
    }
    if idx <= b.Len()/2 {
        for i := 0; i < idx; i++ {
            ins.Rb(b)
        }
    } else {
        for i := 0; i < b.Len()-idx; i++ {
            ins.Rrb(b)
        }
    }
}

func finalRotate(a *stack.Stack, ins *InstructionSet) {
    minIdx := findMinIndex(a)
    rotateAToTop(a, minIdx, ins)
}

// ==================== HELPER FUNCTIONS ====================

func calculateMedian(s *stack.Stack) int {
    if s.Len() == 0 {
        return 0
    }
    temp := make([]int, s.Len())
    copy(temp, s.Values)
    sort.Ints(temp)
    return temp[len(temp)/2]
}

func findMinIndex(a *stack.Stack) int {
    if a.Len() == 0 {
        return 0
    }
    minVal := a.Values[0]
    minIdx := 0
    for i, v := range a.Values {
        if v < minVal {
            minVal = v
            minIdx = i
        }
    }
    return minIdx
}

func findInsertIndex(a *stack.Stack, x int) int {
    if a.Len() == 0 {
        return 0
    }
    
    bestIndex := -1
    bestValue := math.MaxInt32
    
    for i, value := range a.Values {
        if value > x && value < bestValue {
            bestValue = value
            bestIndex = i
        }
    }
    
    if bestIndex == -1 {
        return findMinIndex(a)
    }
    
    return bestIndex
}

func calculateCost(a, b *stack.Stack, bIdx int) int {
    bCost := bIdx
    if bIdx > b.Len()/2 {
        bCost = b.Len() - bIdx
    }
    
    insertIdx := findInsertIndex(a, b.Values[bIdx])
    aCost := insertIdx
    if insertIdx > a.Len()/2 {
        aCost = a.Len() - insertIdx
    }
    return bCost + aCost
}

func findBestElement(a, b *stack.Stack) int {
    if b.Len() == 0 {
        return 0
    }
    bestIdx := 0
    minCost := math.MaxInt32
    for i := range b.Values {
        cost := calculateCost(a, b, i)
        if cost < minCost {
            minCost = cost
            bestIdx = i
        }
    }
    return bestIdx
}

// ==================== UTILITY FUNCTIONS ====================

func isSorted(s *stack.Stack) bool {
    for i := 0; i < s.Len()-1; i++ {
        if s.Values[i] > s.Values[i+1] {
            return false
        }
    }
    return true
}

// ==================== OPTIMIZER ====================

func OptimizeInstructions(instructions []string) []string {
    optimized := []string{}
    i := 0
    for i < len(instructions) {
        if i < len(instructions)-1 {
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

func RemoveRedundancies(instructions []string) []string {
    result := []string{}
    for i := 0; i < len(instructions); i++ {
        if i < len(instructions)-1 {
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
