package main

type Node struct {
	Value string
	Next  *Node
}

type Stack struct {
	top  *Node
	size int
}

func (s *Stack) Push(value string) {
	n := &Node{Value: value}
	if s.top == nil {
		s.top = n
	} else {
		n.Next = s.top
		s.top = n
	}
	s.size++
}

func (s *Stack) Pop() string {
	if s.top == nil {
		s.size = 0
		return ""
	} else {
		v := s.top.Value
		s.top = s.top.Next
		s.size--
		return v
	}
}

func (s *Stack) Empty() bool {
	return s.top == nil
}

func (s *Stack) Size() int {
	if s.top == nil {
		return 0
	}
	return s.size
}
