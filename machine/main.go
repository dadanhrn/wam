package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dadanhrn/wam/common"
	"github.com/dadanhrn/wam/machine/lib"
)

func main() {
	var f common.OutputFile
	d := json.NewDecoder(os.Stdin)
	d.DisallowUnknownFields()
	err := d.Decode(&f)
	if err != nil {
		panic(err)
	}

	m := lib.New(1024, 1024)
	m.Run(f.Instructions, f.Labels)

	fmt.Println("### Registers")
	m.PrintRegisters(10)
	fmt.Println()
	fmt.Println("### Heap")
	m.PrintHeap()
}
