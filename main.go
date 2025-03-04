package main

import (
	"bufio"
	"fmt"
	"os"
	"yaboy/decompiler"
	enc "yaboy/encoder"
)

func main() {
	data, err := enc.ReadBinaryFile(os.Args[1])
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	f, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	res, err := decompiler.DecodeBytes(data)
	if err != nil {
		panic(err)
	}

	for _, s := range res {
		_, err = writer.WriteString(s + "\n")
		if err != nil {
			fmt.Printf("Error writing string output: %s", err)
			panic(err)
		}
		writer.Flush()
	}

}
