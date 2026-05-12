package algorithm

import (
    "push_swap/stack"
    "sort"
    "math"
)

// CalculateMedian returns the median value of the stack
func CalculateMedian(s *stack.Stack) int {
    if s.Len() == 0 {
        return 0
    }
    temp := make([]int, s.Len())
    copy(temp, s.Values)
    sort.Ints(temp)
    return temp[len(temp)/2]
}

// FindMinIndex returns the index of the minimum value
func FindMinIndex(a *stack.Stack) int {
    if a.Len() == 0 {
        return 0
    }
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

// FindMaxIndex returns the index of the maximum value
func FindMaxIndex(a *stack.Stack) int {
    if a.Len() == 0 {
        return 0
    }
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

// FindInsertIndex finds where x should be inserted in a sorted stack
func FindInsertIndex(a *stack.Stack, x int) int {
    if a.Len() == 0 {
        return 0
    }
    
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

// CalculateCost calculates the cost to insert element at bIndex
func CalculateCost(a, b *stack.Stack, bIndex int) int {
    bCost := bIndex
    if bIndex > b.Len()/2 {
        bCost = b.Len() - bIndex
    }
    
    insertIndex := FindInsertIndex(a, b.Values[bIndex])
    aCost := insertIndex
    if insertIndex > a.Len()/2 {
        aCost = a.Len() - insertIndex
    }
    
    return bCost + aCost
}

// FindBestElement finds the element in B with minimal insertion cost
func FindBestElement(a, b *stack.Stack) int {
    if b.Len() == 0 {
        return 0
    }
    
    bestIndex := 0
    minCost := math.MaxInt32
    
    for i := range b.Values {
        cost := CalculateCost(a, b, i)
        if cost < minCost {
            minCost = cost
            bestIndex = i
        }
    }
    
    return bestIndex
}