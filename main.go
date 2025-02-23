package main

import (
	// "os"
	"fmt"
	"yaboy/decompiler"
	enc "yaboy/encoder"
)

func main() {
	data, err := enc.ReadBinaryFile("myrom.gb")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(len(data))
	decompiler.DecodeBytes(data)
}
