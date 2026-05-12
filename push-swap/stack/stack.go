package stack

type Stack struct {
    Values []int
}


func NewStack(values []int) *Stack {
	return &Stack{
		Values: values,
	}
}

func (s *Stack) Push(val int) {
    s.Values = append([]int{val}, s.Values...)
}

func (s *Stack) Pop() (int, bool) {
    if len(s.Values) == 0 {
        return 0, false
    }
    top := s.Values[0]
    s.Values = s.Values[1:]
    return top, true
}

func (s *Stack) Len() int {
    return len(s.Values)
}

func (s *Stack) IsSorted() bool {
    for i := 0; i < len(s.Values)-1; i++ {
        // For stack with top at index 0, sorted means decreasing order
        if s.Values[i] > s.Values[i+1] {
            return false
        }
    }
    return true
}



