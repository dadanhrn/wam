package common

import (
	"encoding/json"
	"fmt"
)

const (
	INST_PUT_STRUCTURE  string = "PUT_STRUCTURE"
	INST_PUT_VARIABLE   string = "PUT_VARIABLE"
	INST_PUT_VALUE      string = "PUT_VALUE"
	INST_SET_VARIABLE   string = "SET_VARIABLE"
	INST_SET_VALUE      string = "SET_VALUE"
	INST_GET_STRUCTURE  string = "GET_STRUCTURE"
	INST_GET_VARIABLE   string = "GET_VARIABLE"
	INST_GET_VALUE      string = "GET_VALUE"
	INST_UNIFY_VARIABLE string = "UNIFY_VARIABLE"
	INST_UNIFY_VALUE    string = "UNIFY_VALUE"
	INST_CALL           string = "CALL"
	INST_PROCEED        string = "PROCEED"
)

type Instruction struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"args,omitempty"`
}

func (inst *Instruction) UnmarshalJSON(data []byte) (err error) {
	var instRaw struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"args"`
	}

	err = json.Unmarshal(data, &instRaw)
	if err != nil {
		return err
	}

	*inst = Instruction{
		Name: instRaw.Name,
	}

	switch instRaw.Name {
	case INST_PUT_STRUCTURE:
		var args InstPutStructure
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_PUT_VARIABLE:
		var args InstPutVariable
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_PUT_VALUE:
		var args InstPutValue
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_SET_VARIABLE:
		var args InstSetVariable
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_SET_VALUE:
		var args InstSetValue
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_GET_STRUCTURE:
		var args InstGetStructure
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_GET_VARIABLE:
		var args InstGetVariable
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_GET_VALUE:
		var args InstGetValue
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_UNIFY_VARIABLE:
		var args InstUnifyVariable
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_UNIFY_VALUE:
		var args InstUnifyValue
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_CALL:
		var args InstCall
		err = json.Unmarshal(instRaw.Arguments, &args)
		inst.Arguments = args
	case INST_PROCEED:
		return nil
	default:
		err = fmt.Errorf("unknown instruction %s", instRaw.Name)
	}

	return err
}

type InstPutStructure struct {
	Identifier   string `json:"identifier"`
	Arity        int    `json:"arity"`
	RegisterAddr int    `json:"register_addr"`
}

type InstPutVariable struct {
	Xn int `json:"Xn"`
	Ai int `json:"Ai"`
}

type InstPutValue struct {
	Xn int `json:"Xn"`
	Ai int `json:"Ai"`
}

type InstSetVariable struct {
	RegisterAddr int `json:"register_addr"`
}

type InstSetValue struct {
	Value int `json:"value"`
}

type InstGetStructure struct {
	Identifier   string `json:"identifier"`
	Arity        int    `json:"arity"`
	RegisterAddr int    `json:"reference"`
}

type InstGetVariable struct {
	Xn int `json:"Xn"`
	Ai int `json:"Ai"`
}

type InstGetValue struct {
	Xn int `json:"Xn"`
	Ai int `json:"Ai"`
}

type InstUnifyVariable struct {
	Reference int `json:"reference"`
}

type InstUnifyValue struct {
	Reference int `json:"reference"`
}

type InstCall struct {
	Label string `json:"label"`
}
