package main

type Stack struct {
	st []int
}

func NewStack() Stack {
	return Stack{
		st: make([]int, 0),
	}
}

func (st *Stack) Push(nd int) {
	st.st = append(st.st, nd)
}

func (st *Stack) Pop() int {
	r := st.st[len(st.st)-1]
	st.st = st.st[:len(st.st)-1]
	return r
}

func (st *Stack) Len() int {
	return len(st.st)
}
