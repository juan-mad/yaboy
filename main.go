package main

import (
	// "os"
	"fmt"
	enc "yaboy/encoder"
)

var LOAD_INPUT string = `LD B,B
LD B,C
LD B,D`

func main() {
	data, err := enc.ReadBinaryFile("decoder/out.bin")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(len(data))
}
