package lib

import (
	"container/list"

	"github.com/dadanhrn/wam/common"
)

func GenerateRegisterMap(root common.LFunctor) []interface{} {
	q := list.New()
	stq := list.New()
	encountered := make(map[string]int)
	mapping := make([]interface{}, 0)

	q.PushFront(root)
	stq.PushFront(-1)

	for ref := 0; q.Len() > 0; ref++ {
		el := q.Remove(q.Back())
		prevOccurrence := -1

		switch item := el.(type) {
		case common.LFunctor:
			mapping = append(mapping, &common.RegFunctor{
				Identifier: item.Identifier,
				Subterms:   make([]common.RegSubterm, 0),
			})

			for _, st := range item.Subterms {
				parent := len(mapping) - 1
				stq.PushFront(parent)
				q.PushFront(st)
			}
		case common.LConcreteVariable:
			label := item.Identifier
			if i, ok := encountered[label]; !ok {
				mapping = append(mapping, common.RegConcreteVariable{
					Identifier: item.Identifier,
				})

				encountered[label] = len(mapping) - 1
			} else {
				prevOccurrence = i
			}
		}

		possibleParent := stq.Remove(stq.Back()).(int)
		if possibleParent > -1 {
			candidateRef := len(mapping) - 1
			if prevOccurrence > -1 {
				candidateRef = prevOccurrence
			}

			f := mapping[possibleParent].(*common.RegFunctor)
			f.Subterms = append(f.Subterms, common.RegSubterm{
				Reference: candidateRef,
			})
		}
	}

	for i, item := range mapping {
		if p, ok := item.(*common.RegFunctor); ok {
			mapping[i] = *p
		}
	}

	return mapping
}

func traverse(current int, deps map[int][]int, addedRegSeq *map[int]struct{}, regSeq *[]int) {
	if dependents, ok := deps[current]; ok {
		for _, x := range dependents {
			if _, isAdded := (*addedRegSeq)[x]; !isAdded {
				*regSeq = append(*regSeq, x)
				(*addedRegSeq)[x] = struct{}{}
			}
			traverse(x, deps, addedRegSeq, regSeq)
		}
	}
}
