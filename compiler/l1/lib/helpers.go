package lib

import (
	"container/list"

	"github.com/dadanhrn/wam/common"
)

var (
	varEncountered map[string]int
	varOffset      int = 0
	vars           map[string]int
	strs           map[string]int
	deferredStr    *list.List
)

func trv_query(f common.LFunctor, level int) []common.Instruction {
	if level == 0 {
		varEncountered = make(map[string]int)
		varOffset = len(f.Subterms)
		vars = make(map[string]int)
	}

	instructions := make([]common.Instruction, 0)

	for i, s := range f.Subterms {
		switch subterm := s.(type) {
		case common.LFunctor:
			instructions = append(instructions, common.Instruction{
				Name: common.INST_PUT_STRUCTURE,
				Arguments: common.InstPutStructure{
					Identifier:   subterm.Identifier,
					Arity:        len(subterm.Subterms),
					RegisterAddr: i,
				},
			})
			stInst := trv_query(subterm, level+1)
			instructions = append(instructions, stInst...)
		case common.LConcreteVariable:
			if level == 0 {
				instructions = append(instructions, common.Instruction{
					Name: common.INST_PUT_VARIABLE,
					Arguments: common.InstPutVariable{
						Xn: varOffset + vars[subterm.Identifier],
						Ai: i,
					},
				})
				vars[subterm.Identifier] = len(vars)
			} else {
				if _, ok := vars[subterm.Identifier]; ok {
					instructions = append(instructions, common.Instruction{
						Name: common.INST_SET_VALUE,
						Arguments: common.InstSetValue{
							Value: varOffset + vars[subterm.Identifier],
						},
					})
				} else {
					vars[subterm.Identifier] = len(vars)
					instructions = append(instructions, common.Instruction{
						Name: common.INST_SET_VARIABLE,
						Arguments: common.InstSetVariable{
							RegisterAddr: varOffset + vars[subterm.Identifier],
						},
					})
				}
			}
			varEncountered[subterm.Identifier] = len(varEncountered) + 1
		}
	}

	return instructions
}

func trv_program(f common.LFunctor, level int) ([]common.Instruction, []common.Instruction) {
	if level == 0 {
		varOffset = len(f.Subterms)
		vars = make(map[string]int)
		strs = make(map[string]int)
		deferredStr = list.New()
	}

	deferredInst := make([]common.Instruction, 0)
	instructions := make([]common.Instruction, 0)
	for i, s := range f.Subterms {
		switch subterm := s.(type) {
		case common.LFunctor:
			var ref int
			// id := fmt.Sprintf("%s/%d", subterm.Identifier, len(subterm.Subterms))
			if level == 0 {
				instructions = append(instructions, common.Instruction{
					Name: common.INST_GET_STRUCTURE,
					Arguments: common.InstGetStructure{
						Identifier:   subterm.Identifier,
						Arity:        len(subterm.Subterms),
						RegisterAddr: i,
					},
				})
				stInst, stDef := trv_program(subterm, level+1)
				instructions = append(instructions, stInst...)
				deferredInst = append(deferredInst, stDef...)
			} else {
				x, ok := strs[subterm.Identifier]
				if !ok {
					x = len(strs)
					strs[subterm.Identifier] = x
				}
				// len(vars): does not cover the case when some variables are not encountered before the first nested functor
				ref = varOffset + len(vars) + x

				instructions = append(instructions, common.Instruction{
					Name: common.INST_UNIFY_VARIABLE,
					Arguments: common.InstUnifyVariable{
						Reference: ref,
					},
				})

				// deferredStr.PushBack()

				deferredInst = append(deferredInst, common.Instruction{
					Name: common.INST_GET_STRUCTURE,
					Arguments: common.InstGetStructure{
						Identifier:   subterm.Identifier,
						Arity:        len(subterm.Subterms),
						RegisterAddr: ref,
					},
				})
				stInst, stDef := trv_program(subterm, level+1)
				deferredInst = append(deferredInst, stInst...)
				deferredInst = append(deferredInst, stDef...)
			}
		case common.LConcreteVariable:
			// ASSUMING VARIABLES ARE ONLY IN LEVEL 0
			if level == 0 {
				if x, ok := vars[subterm.Identifier]; !ok {
					// first encounter
					x = len(vars)
					vars[subterm.Identifier] = x
					instructions = append(instructions, common.Instruction{
						Name: common.INST_GET_VARIABLE,
						Arguments: common.InstGetVariable{
							Xn: varOffset + x,
							Ai: i,
						},
					})
				} else {
					// encountered before
					instructions = append(instructions, common.Instruction{
						Name: common.INST_GET_VALUE,
						Arguments: common.InstGetValue{
							Xn: varOffset + x,
							Ai: i,
						},
					})
				}
			} else {
				x, ok := vars[subterm.Identifier]
				if !ok {
					x = len(vars)
					vars[subterm.Identifier] = x
				}

				instructions = append(instructions, common.Instruction{
					Name: common.INST_UNIFY_VARIABLE,
					Arguments: common.InstUnifyVariable{
						Reference: varOffset + x,
					},
				})
			}

		}
	}

	if level == 0 {
		instructions = append(instructions, deferredInst...)
		instructions = append(instructions, common.Instruction{
			Name: common.INST_PROCEED,
		})
	}

	return instructions, deferredInst
}
