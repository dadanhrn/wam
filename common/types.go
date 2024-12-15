package common

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RegFunctor struct {
	Identifier string
	Subterms   []RegSubterm
}

func (r RegFunctor) MarshalJSON() ([]byte, error) {
	if len(r.Subterms) == 0 {
		return json.Marshal(r.Identifier)
	}

	stRefs := make([]string, len(r.Subterms))
	for i, st := range r.Subterms {
		stRef, err := json.Marshal(st)
		if err != nil {
			return nil, err
		}

		stRefs[i] = string(stRef)
	}

	stRefsStr := strings.Join(stRefs, ", ")
	data := fmt.Sprintf("%s(%s)", r.Identifier, stRefsStr)
	return json.Marshal(data)
}

type RegSubterm struct {
	Reference int `json:"reference"`
}

func (r RegSubterm) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf("X%d", r.Reference)
	return json.Marshal(data)
}

type RegConcreteVariable struct {
	Identifier string `json:"identifier"`
}

func (r RegConcreteVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Identifier)
}
