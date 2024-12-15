package lib

import (
	"fmt"

	"github.com/dadanhrn/wam/common"
)

type Machine struct {
	heap      []HeapCell
	register  []int
	heapCtr   int
	registerS int
	registerP int
	writeMode bool
}

func New(heapSize int, registerN int) *Machine {
	return &Machine{
		heap:      make([]HeapCell, heapSize),
		register:  make([]int, registerN),
		heapCtr:   0,
		registerS: 0,
		registerP: 0,
		writeMode: false,
	}
}

func (m *Machine) Run(instructions []common.Instruction, labels map[string]int) {
Instructions_Run:
	for {
		inst := instructions[m.registerP]

		switch inst.Name {
		case common.INST_PUT_STRUCTURE:
			m.PutStructure(inst.Arguments.(common.InstPutStructure))
		case common.INST_PUT_VARIABLE:
			m.PutVariable(inst.Arguments.(common.InstPutVariable))
		case common.INST_PUT_VALUE:
			m.PutValue(inst.Arguments.(common.InstPutValue))
		case common.INST_SET_VARIABLE:
			m.SetVariable(inst.Arguments.(common.InstSetVariable))
		case common.INST_SET_VALUE:
			m.SetValue(inst.Arguments.(common.InstSetValue))
		case common.INST_GET_VARIABLE:
			m.GetVariable(inst.Arguments.(common.InstGetVariable))
		case common.INST_GET_VALUE:
			m.GetValue(inst.Arguments.(common.InstGetValue))
		case common.INST_GET_STRUCTURE:
			m.GetStructure(inst.Arguments.(common.InstGetStructure))
		case common.INST_UNIFY_VARIABLE:
			m.UnifyVariable(inst.Arguments.(common.InstUnifyVariable))
		case common.INST_UNIFY_VALUE:
			m.UnifyValue(inst.Arguments.(common.InstUnifyValue))
		case common.INST_CALL:
			arg := inst.Arguments.(common.InstCall)
			nextP, ok := labels[arg.Label]
			if !ok {
				panic(fmt.Sprintf("no call label %s", arg.Label))
			}
			m.registerP = nextP
		case common.INST_PROCEED:
			break Instructions_Run
		}
	}
}
