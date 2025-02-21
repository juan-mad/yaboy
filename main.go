package main

import (
	"os"
	enc "yaboy/encoder"
)

func main() {
	enc.EncodeTextFile(os.Args[1], os.Args[2])
}
