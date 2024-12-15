package lib

import (
	"fmt"

	"github.com/dadanhrn/wam/common"
)

func CompileQuery(query common.LFunctor) (instructions []common.Instruction) {
	instructions = make([]common.Instruction, 0)
	instructions = trv_query(query, 0)
	instructions = append(instructions, common.Instruction{
		Name: common.INST_CALL,
		Arguments: common.InstCall{
			Label: fmt.Sprintf("%s/%d", query.Identifier, len(query.Subterms)),
		},
	})

	return
}

func CompileProgram(program []common.LFunctor) (instructions []common.Instruction, callLabel map[string]int) {
	instructions = make([]common.Instruction, 0)
	callLabel = make(map[string]int)

	for _, fact := range program {
		label := fmt.Sprintf("%s/%d", fact.Identifier, len(fact.Subterms))
		callLabel[label] = len(instructions)

		inst, _ := trv_program(fact, 0)
		instructions = append(instructions, inst...)
	}

	return
}
