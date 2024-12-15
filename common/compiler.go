package common

import (
	"fmt"
)

const (
	NODE_FUNCTOR  string = "FUNCTOR"
	NODE_CONSTANT string = "CONSTANT"
	NODE_VARIABLE string = "VARIABLE"
)

type ASTNode struct {
	Type       string    `json:"type"`
	Identifier string    `json:"identifier"`
	Subterms   []ASTNode `json:"subterms"`
}

func (n ASTNode) Transform() (v interface{}, err error) {
	switch n.Type {
	case NODE_CONSTANT:
		fallthrough
	case NODE_FUNCTOR:
		f := LFunctor{
			Identifier: n.Identifier,
		}

		if len(n.Subterms) > 0 {
			f.Subterms = make([]interface{}, len(n.Subterms))
			for i, st := range n.Subterms {
				f.Subterms[i], err = st.Transform()
				if err != nil {
					return
				}
			}
		}

		v = f
	case NODE_VARIABLE:
		v = LConcreteVariable{
			Identifier: n.Identifier,
		}
	default:
		err = fmt.Errorf("unknown type %s", n.Type)
	}

	return
}

type LFunctor struct {
	Identifier string `json:"identifier"`
	Subterms   []interface{}
}

type LConcreteVariable struct {
	Identifier string
}

type InputFile struct {
	Program []ASTNode `json:"program"`
	Query   ASTNode   `json:"query"`
}

type OutputFile struct {
	Instructions []Instruction  `json:"instructions"`
	Labels       map[string]int `json:"labels,omitempty"`
	RegisterMaps RegisterMaps   `json:"register_maps,omitempty"`
}

type RegisterMaps struct {
	Query   []interface{} `json:"query"`
	Program []interface{} `json:"program"`
}
