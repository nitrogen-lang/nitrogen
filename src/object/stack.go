package object

type Stack struct {
	head   *stackElement
	length int
}

type stackElement struct {
	val  Object
	prev *stackElement
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(val Object) {
	s.head = &stackElement{
		val:  val,
		prev: s.head,
	}
	s.length++
}

func (s *Stack) GetFront() Object {
	if s.head == nil {
		return nil
	}
	return s.head.val
}

func (s *Stack) Pop() Object {
	if s.head == nil {
		return nil
	}
	r := s.head.val
	s.head = s.head.prev
	s.length--
	return r
}

func (s *Stack) Len() int {
	return s.length
}
