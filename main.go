package main

import (
	// "os"
	"fmt"
	dec "yaboy/decoder"
	// enc "yaboy/encoder"
)

var LOAD_INPUT string = `LD B,B
LD B,C
LD B,D`

func main() {
	data, err := dec.ReadBinaryFile("decoder/out.bin")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(len(data))
}
