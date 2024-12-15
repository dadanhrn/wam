package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dadanhrn/wam/common"
	l1 "github.com/dadanhrn/wam/compiler/l1/lib"
	machine "github.com/dadanhrn/wam/machine/lib"
)

func main() {
	var f common.InputFile
	d := json.NewDecoder(os.Stdin)
	d.DisallowUnknownFields()
	err := d.Decode(&f)
	if err != nil {
		panic(err)
	}

	if len(f.Program) < 1 {
		panic("program has to contain at least one term")
	}

	program := make([]common.LFunctor, len(f.Program))
	var ok bool
	for i, node := range f.Program {
		p, err := node.Transform()
		if err != nil {
			panic(err)
		}

		program[i], ok = p.(common.LFunctor)
		if !ok {
			panic("root program term has to be a functor")
		}
	}

	q, err := f.Query.Transform()
	if err != nil {
		panic(err)
	}

	query, ok := q.(common.LFunctor)
	if !ok {
		panic("root query term has to be a functor")
	}

	queryInst := l1.CompileQuery(query)
	programInst, callLabels := l1.CompileProgram(program)

	for label, offset := range callLabels {
		callLabels[label] = offset + len(queryInst)
	}

	instructions := make([]common.Instruction, 0)
	instructions = append(instructions, queryInst...)
	instructions = append(instructions, programInst...)

	m := machine.New(1024, 1024)
	m.Run(instructions, callLabels)

	fmt.Println("### Registers")
	m.PrintRegisters(10)
	fmt.Println()
	fmt.Println("### Heap")
	m.PrintHeap()
}
