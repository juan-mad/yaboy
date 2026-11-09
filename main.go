package main

import (
	"yaboy/gameboy"
	"yaboy/tui"
	// "fmt"
	// "yaboy/decompiler"
	// "yaboy/encoder"
)

func main() {
	gb := gameboy.NewGameBoy()

	ui := tui.NewTui(&gb)
	ui.Run()
	// input, _ := encoder.ReadBinaryFile("test.asm")
	// decompiled, _ := decompiler.DecodeBytes(input)
	// for _, instr := range decompiled {
	// 	fmt.Println(instr)
	// }
}
