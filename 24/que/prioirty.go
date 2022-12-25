package que

type Priorityer interface {
	Priority() int
}

type Priority[T Priorityer] struct {
	s []T
}

func NewPriority[T Priorityer]() *Priority[T] {
	return &Priority[T]{s: make([]T, 0)}
}

func (p *Priority[T]) Len() int {
	return len(p.s)
}

func (p *Priority[T]) children(idx int) []T {
	var c []T
	if idx*2+1 < len(p.s) {
		c = append(c, p.s[idx*2+1])
	}
	if idx*2+2 < len(p.s) {
		c = append(c, p.s[idx*2+2])
	}
	return c
}

func (p *Priority[T]) childrenIdx(idx int) []int {
	var c []int
	if idx*2+1 < len(p.s) {
		c = append(c, idx*2+1)
	}
	if idx*2+2 < len(p.s) {
		c = append(c, idx*2+2)
	}
	return c
}

func (p *Priority[T]) parent(idx int) T {
	if idx == 0 {
		return *new(T)
	}
	return p.s[(idx-1)/2]
}

func (p *Priority[T]) swap(i, j int) {
	p.s[i], p.s[j] = p.s[j], p.s[i]
}

func (p *Priority[T]) bubbleUp(idx int) {
	if idx == 0 {
		return
	}
	parent := p.parent(idx)
	if parent.Priority() > p.s[idx].Priority() {
		p.swap(idx, (idx-1)/2)
		p.bubbleUp((idx - 1) / 2)
	}
}

func (p *Priority[T]) bubbleDown(idx int) {
	if idx*2+1 >= len(p.s) {
		return
	}
	min := idx
	children := p.children(idx)
	childrenIdx := p.childrenIdx(idx)

	for i, c := range children {
		if c.Priority() < p.s[min].Priority() {
			min = childrenIdx[i]
		}
	}

	if min != idx {
		p.swap(idx, min)
		p.bubbleDown(min)
	}
}

func (p *Priority[T]) Push(t T) {
	p.s = append(p.s, t)
	p.bubbleUp(len(p.s) - 1)
}

func (p *Priority[T]) Pop() T {
	if len(p.s) == 0 {
		return *new(T)
	}
	t := p.s[0]
	p.s[0] = p.s[len(p.s)-1]
	p.s = p.s[:len(p.s)-1]
	p.bubbleDown(0)
	return t
}
