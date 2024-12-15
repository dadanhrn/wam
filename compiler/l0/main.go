package main

import (
	"encoding/json"
	"os"

	"github.com/dadanhrn/wam/common"
	"github.com/dadanhrn/wam/compiler/l0/lib"
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

	queryInst, queryRegisterMap := lib.CompileQuery(query)
	programInst, programRegisterMap := lib.CompileProgram(program)

	var out common.OutputFile
	out.Instructions = make([]common.Instruction, 0)
	out.Instructions = append(out.Instructions, queryInst...)
	out.Instructions = append(out.Instructions, programInst...)
	out.Instructions = append(out.Instructions, common.Instruction{
		Name: common.INST_PROCEED,
	})
	out.RegisterMaps.Query = queryRegisterMap
	out.RegisterMaps.Program = programRegisterMap

	err = json.NewEncoder(os.Stdout).Encode(out)
	if err != nil {
		panic(err)
	}
}
