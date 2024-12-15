package lib

import (
	"github.com/dadanhrn/wam/common"
)

func CompileQuery(query common.LFunctor) (instructions []common.Instruction, registerMap []interface{}) {
	registerMap = GenerateRegisterMap(query)

	deps := make(map[int][]int)
	roots := make([]int, 0)

	for i, v := range registerMap {

		// Check functor
		if f, ok := v.(common.RegFunctor); ok {
			allVars := true
			maxDep := -1

			// Iterate subterms
			for _, st := range f.Subterms {
				// Check for variables
				if _, ok2 := registerMap[st.Reference].(common.RegConcreteVariable); !ok2 {
					allVars = false
					if st.Reference > maxDep {
						maxDep = st.Reference
					}
				}
			}

			if _, ok := deps[maxDep]; !ok {
				deps[maxDep] = make([]int, 0)
			}

			deps[maxDep] = append(deps[maxDep], i)

			if allVars {
				roots = append(roots, i)
			}
		}
	}

	regSeq := make([]int, 0)
	addedRegSeq := make(map[int]struct{})
	regSeq = append(regSeq, roots...)

	for _, x := range roots {
		addedRegSeq[x] = struct{}{}
		traverse(x, deps, &addedRegSeq, &regSeq)
	}

	instructions = make([]common.Instruction, 0)
	addedInst := make(map[int]struct{})

	for _, x := range regSeq {
		f, ok := registerMap[x].(common.RegFunctor)
		if !ok {
			continue
		}

		instructions = append(instructions, common.Instruction{
			Name: common.INST_PUT_STRUCTURE,
			Arguments: common.InstPutStructure{
				Identifier:   f.Identifier,
				Arity:        len(f.Subterms),
				RegisterAddr: x,
			},
		})

		addedInst[x] = struct{}{}

		for _, stX := range f.Subterms {
			var inst common.Instruction
			if _, added := addedInst[stX.Reference]; added {
				inst = common.Instruction{
					Name: common.INST_SET_VALUE,
					Arguments: common.InstSetValue{
						Value: stX.Reference,
					},
				}
			} else {
				inst = common.Instruction{
					Name: common.INST_SET_VARIABLE,
					Arguments: common.InstSetVariable{
						RegisterAddr: stX.Reference,
					},
				}
				addedInst[stX.Reference] = struct{}{}
			}

			instructions = append(instructions, inst)
		}
	}

	return
}

func CompileProgram(program common.LFunctor) (instructions []common.Instruction, registerMap []interface{}) {
	registerMap = GenerateRegisterMap(program)
	addedInst := make(map[int]struct{})

	for i, x := range registerMap {
		f, ok := x.(common.RegFunctor)
		if !ok {
			// Skip if not functor
			continue
		}

		instructions = append(instructions, common.Instruction{
			Name: common.INST_GET_STRUCTURE,
			Arguments: common.InstGetStructure{
				Identifier:   f.Identifier,
				Arity:        len(f.Subterms),
				RegisterAddr: i,
			},
		})

		addedInst[i] = struct{}{}

		for _, stX := range f.Subterms {
			var inst common.Instruction
			if _, added := addedInst[stX.Reference]; added {
				inst = common.Instruction{
					Name: common.INST_UNIFY_VALUE,
					Arguments: common.InstUnifyValue{
						Reference: stX.Reference,
					},
				}
			} else {
				inst = common.Instruction{
					Name: common.INST_UNIFY_VARIABLE,
					Arguments: common.InstUnifyVariable{
						Reference: stX.Reference,
					},
				}
				addedInst[stX.Reference] = struct{}{}
			}

			instructions = append(instructions, inst)
		}
	}

	return
}
