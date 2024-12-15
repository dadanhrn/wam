package lib

import (
	"container/list"
	"fmt"
)

func (m *Machine) unify(a1 int, a2 int) {
	pdl := list.New()
	pdl.PushBack(a1)
	pdl.PushBack(a2)
	fail := false

	for !(pdl.Len() == 0 || fail) {
		d1 := m.deref(pdl.Remove(pdl.Back()).(int))
		d2 := m.deref(pdl.Remove(pdl.Back()).(int))

		if d1 != d2 {
			cell1 := m.heap[d1]
			t1 := cell1.Type
			v1 := cell1.Data.(int)

			cell2 := m.heap[d2]
			t2 := cell2.Type
			v2 := cell2.Data.(int)

			if t1 == HEAP_REF || t2 == HEAP_REF {
				m.bind(d1, d2)
			} else {
				str1 := m.heap[v1].Data.(HeapFunctor)
				f1 := str1.Identifier
				n1 := str1.Arity

				str2 := m.heap[v2].Data.(HeapFunctor)
				f2 := str2.Identifier
				n2 := str2.Arity

				if f1 == f2 && n1 == n2 {
					for i := 1; i <= n1; i++ {
						pdl.PushBack(v1 + i)
						pdl.PushBack(v2 + i)
					}
				} else {
					fail = true
					panic("unify fail")
				}
			}
		}
	}
}

func (m *Machine) deref(addr int) int {
	c := m.heap[addr]
	if c.Type == HEAP_REF && c.Data != addr {
		return m.deref(c.Data.(int))
	}

	return addr
}

func (m *Machine) bind(a1 int, a2 int) {
	t1 := m.heap[a1].Type
	t2 := m.heap[a2].Type

	if t1 == HEAP_REF && ((t2 != HEAP_REF) || a2 < a1) {
		m.heap[a1] = m.heap[a2]
		m.trail(a1)
	} else {
		m.heap[a2] = m.heap[a1]
		m.trail(a2)
	}
}

func (m *Machine) PrintHeap() {
	for i := 0; i < m.heapCtr; i++ {
		c := m.heap[i]
		fmt.Printf("%d: %s %#v\n", i, c.Type, c.Data)
	}
}

func (m *Machine) PrintRegisters(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("%d: %d\n", i, m.register[i])
	}
}

func (m *Machine) trail(addr int) {}
