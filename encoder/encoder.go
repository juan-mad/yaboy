package encoder

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"slices"
	"strings"
)

var PREFIXES = []string{"CB", "DD", "ED", "FD"}
var BIT8REG = []string{"B", "C", "D", "E", "H", "L", "[HL]", "A"}

// var BIT16REG = []string{"BC", "DE", "HL", "SP", "AF"}

var REGMAPPER = map[string]uint8{
	"B":    0,
	"C":    1,
	"D":    2,
	"E":    3,
	"H":    4,
	"L":    5,
	"[HL]": 6,
	"A":    7,
}

func IsPrefix(s string) bool {
	return slices.Contains(PREFIXES, s)
}

func Is8BitRegister(s string) bool {
	return slices.Contains(BIT8REG, s)
}

func encodeLoad(argStr string) ([]uint8, error) {
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) != 2 {
		return []uint8{}, fmt.Errorf("could not parse LD arguments, expected only one `,`: %s", argStr)
	}

	if Is8BitRegister(args[0]) && Is8BitRegister(args[1]) {
		var opcode uint8 = 0b01000000 | (REGMAPPER[args[0]] << 3) | (REGMAPPER[args[1]])
		return []uint8{opcode}, nil
	}

	return []uint8{}, fmt.Errorf("Not implemented load: %s", argStr)
}

func writeDataToFile(filename string, buffer []uint8) error {
	// Create a new file (or truncate existing one)
	// The 0644 permission means: owner can read/write, others can only read
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close() // Ensure file is closed when we're done

	// Write the actual data slice
	err = binary.Write(file, binary.LittleEndian, buffer)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

func EncodeTextFile(inputPath, outputPath string) error {
	readFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	buf := []uint8{}

	for scanner.Scan() {
		line := scanner.Text()
		chunks := strings.Fields(strings.TrimSpace(line))
		if IsPrefix(chunks[0]) {
			// TODO
		} else {
			opcodeString := chunks[0]

			switch opcodeString {
			case "LD":
				var inp string
				if len(chunks) == 3 {
					inp = chunks[1] + chunks[2]
				} else if len(chunks) == 2 {
					inp = chunks[1]
				} else {
					panic(fmt.Sprintf("Unexpected split of input line: %s\n", line))
				}
				bindata, err := encodeLoad(inp)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			}
		}
	}

	// Now we write into file
	writeDataToFile(outputPath, buf)

	return nil
}
