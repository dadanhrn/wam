package main

import (
	"encoding/json"
	"os"

	"github.com/dadanhrn/wam/common"
	"github.com/dadanhrn/wam/compiler/l1/lib"
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

	queryInst := lib.CompileQuery(query)
	programInst, callLabels := lib.CompileProgram(program)

	for label, offset := range callLabels {
		callLabels[label] = offset + len(queryInst)
	}

	var out common.OutputFile
	out.Instructions = make([]common.Instruction, 0)
	out.Instructions = append(out.Instructions, queryInst...)
	out.Instructions = append(out.Instructions, programInst...)
	out.Labels = callLabels

	err = json.NewEncoder(os.Stdout).Encode(out)
	if err != nil {
		panic(err)
	}
}
