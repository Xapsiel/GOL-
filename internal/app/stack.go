package app

type stack struct {
	size int
	head *node
}

type node struct {
	elem string
	next *node
}

func NewStack() *stack {
	return &stack{size: 0}
}

func (s *stack) Push(id string) {
	s.size++
	s.head = &node{elem: id, next: s.head}
}

func (s *stack) Pop() string {
	if s.size == 0 {
		return ""
	}
	s.size--
	elem := s.head.elem
	s.head = s.head.next
	return elem
}
