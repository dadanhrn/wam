package lib

import (
	"github.com/dadanhrn/wam/common"
)

func (m *Machine) PutStructure(arg common.InstPutStructure) {
	m.register[arg.RegisterAddr] = m.heapCtr
	m.heap[m.heapCtr] = HeapCell{
		Type: HEAP_STR,
		Data: m.heapCtr + 1,
	}
	m.heapCtr += 1

	m.heap[m.heapCtr] = HeapCell{
		Type: HEAP_STR_DATA,
		Data: HeapFunctor{
			Identifier: arg.Identifier,
			Arity:      arg.Arity,
		},
	}
	m.heapCtr += 1
	m.registerP = m.registerP + 1
}

func (m *Machine) PutVariable(arg common.InstPutVariable) {
	m.heap[m.heapCtr] = HeapCell{
		Type: HEAP_REF,
		Data: m.heapCtr,
	}
	m.register[arg.Xn] = m.heapCtr
	m.register[arg.Ai] = m.heapCtr
	m.heapCtr += 1
	m.registerP += 1
}

func (m *Machine) PutValue(arg common.InstPutValue) {
	m.register[arg.Ai] = m.register[arg.Xn]
	m.registerP += 1
}

func (m *Machine) SetVariable(arg common.InstSetVariable) {
	m.register[arg.RegisterAddr] = m.heapCtr
	m.heap[m.heapCtr] = HeapCell{
		Type: HEAP_REF,
		Data: m.heapCtr,
	}
	m.heapCtr += 1
	m.registerP += 1
}

func (m *Machine) SetValue(arg common.InstSetValue) {
	data := m.heap[m.register[arg.Value]]
	m.heap[m.heapCtr] = HeapCell{
		Type: data.Type,
		Data: data.Data,
	}
	m.heapCtr += 1
	m.registerP += 1
}

func (m *Machine) GetVariable(arg common.InstGetVariable) {
	m.register[arg.Xn] = m.register[arg.Ai]
	m.registerP += 1
}

func (m *Machine) GetValue(arg common.InstGetValue) {
	m.unify(m.register[arg.Xn], m.register[arg.Ai])
	m.registerP += 1
}

func (m *Machine) GetStructure(arg common.InstGetStructure) {
	addr := m.deref(m.register[arg.RegisterAddr])
	cell := m.heap[addr]
	switch cell.Type {
	case HEAP_REF:
		m.heap[m.heapCtr] = HeapCell{
			Type: HEAP_STR,
			Data: m.heapCtr + 1,
		}
		m.heap[m.heapCtr+1] = HeapCell{
			Type: HEAP_STR_DATA,
			Data: HeapFunctor{
				Identifier: arg.Identifier,
				Arity:      arg.Arity,
			},
		}
		m.bind(addr, m.heapCtr)
		m.heapCtr = m.heapCtr + 2
		m.writeMode = true
	case HEAP_STR:
		fAddr := cell.Data.(int)
		cellB := m.heap[fAddr]
		cellF := cellB.Data.(HeapFunctor)
		if cellF.Identifier == arg.Identifier && cellF.Arity == arg.Arity {
			m.registerS = fAddr + 1
			m.writeMode = false
		}
	}
	m.registerP += 1
}

func (m *Machine) UnifyVariable(arg common.InstUnifyVariable) {
	if m.writeMode {
		m.heap[m.heapCtr] = HeapCell{
			Type: HEAP_REF,
			Data: m.heapCtr,
		}
		m.register[arg.Reference] = m.heapCtr
		m.heapCtr = m.heapCtr + 1
	} else {
		m.register[arg.Reference] = m.registerS
	}
	m.registerS = m.registerS + 1
	m.registerP = m.registerP + 1
}

func (m *Machine) UnifyValue(arg common.InstUnifyValue) {
	if m.writeMode {
		m.heap[m.heapCtr] = m.heap[m.register[arg.Reference]]
		m.heapCtr = m.heapCtr + 1
	} else {
		m.unify(m.register[arg.Reference], m.registerS)
	}
	m.registerS = m.registerS + 1
	m.registerP = m.registerP + 1
}
