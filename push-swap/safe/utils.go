package algorithm

import (
    "push_swap/stack"
    "sort"
	"push_swap/instructions"
    "math"
)

// Median calculation
func calculateMedian(s *stack.Stack) int {
    if s.Len() == 0 {
        return 0
    }
    temp := make([]int, s.Len())
    copy(temp, s.Values)
    sort.Ints(temp)
    return temp[len(temp)/2]
}
// Find minimum index
func findMinIndex(a *stack.Stack) int {
    minVal := a.Values[0]
    minIndex := 0
    
    for i, v := range a.Values {
        if v < minVal {
            minVal = v
            minIndex = i
        }
    }
    return minIndex
}

// Find maximum index
func findMaxIndex(a *stack.Stack) int {
    maxVal := a.Values[0]
    maxIndex := 0
    
    for i, v := range a.Values {
        if v > maxVal {
            maxVal = v
            maxIndex = i
        }
    }
    return maxIndex
}

// Rotate to top
func rotateToTop(a *stack.Stack, index int) {
    if index == 0 {
        return
    }
    
    if index <= a.Len()/2 {
        for i := 0; i < index; i++ {
            instructions.Ra(a)
        }
    } else {
        for i := 0; i < a.Len()-index; i++ {
            instructions.Rra(a)
        }
    }
}

// Find insert index
func findInsertIndex(a *stack.Stack, x int) int {
    if x < a.Values[0] {
        return 0
    }
    if x > a.Values[a.Len()-1] {
        return a.Len()
    }
    
    for i := 0; i < a.Len()-1; i++ {
        if a.Values[i] < x && x < a.Values[i+1] {
            return i + 1
        }
    }
    return 0
}

// Cost calculation for insertion
func calculateCost(a, b *stack.Stack, bIndex int) int {
    // Cost to bring b element to top
    bCost := bIndex
    if bIndex > b.Len()/2 {
        bCost = b.Len() - bIndex
    }
    
    // Cost to prepare A
    insertIndex := findInsertIndex(a, b.Values[bIndex])
    aCost := insertIndex
    if insertIndex > a.Len()/2 {
        aCost = a.Len() - insertIndex
    }
    
    return bCost + aCost
}

// Find best element to insert
func findBestElement(a, b *stack.Stack) int {
    bestIndex := 0
    minCost := math.MaxInt32
    
    for i := range b.Values {
        cost := calculateCost(a, b, i)
        if cost < minCost {
            minCost = cost
            bestIndex = i
        }
    }
    
    return bestIndex
}