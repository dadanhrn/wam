package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dadanhrn/wam/common"
	l0 "github.com/dadanhrn/wam/compiler/l0/lib"
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

	if len(f.Program) != 1 {
		panic("program has to contain exactly one term")
	}

	p, err := f.Program[0].Transform()
	if err != nil {
		panic(err)
	}

	program, ok := p.(common.LFunctor)
	if !ok {
		panic("root program term has to be a functor")
	}

	q, err := f.Query.Transform()
	if err != nil {
		panic(err)
	}

	query, ok := q.(common.LFunctor)
	if !ok {
		panic("root query term has to be a functor")
	}

	queryInst, _ := l0.CompileQuery(query)
	programInst, _ := l0.CompileProgram(program)

	instructions := make([]common.Instruction, 0)
	instructions = append(instructions, queryInst...)
	instructions = append(instructions, programInst...)
	instructions = append(instructions, common.Instruction{
		Name: common.INST_PROCEED,
	})

	m := machine.New(1024, 1024)
	m.Run(instructions, nil)

	fmt.Println("### Registers")
	m.PrintRegisters(10)
	fmt.Println()
	fmt.Println("### Heap")
	m.PrintHeap()
}
