package encoder

import "os"
import "fmt"
import "encoding/binary"
import "bufio"

import "strings"

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
		if len(chunks) == 0 || (len(chunks) == 1 && chunks[0] == "") {
			return fmt.Errorf("empty line encountered!")
		}
		if IsPrefix(chunks[0]) {
			fmt.Println(chunks[0])
			if len(chunks) == 1 {
				return fmt.Errorf("Invalid instruction, only CB prefix!")
			}
			instruction := chunks[1]
			if len(chunks) == 2 {
				return fmt.Errorf("Invalid line: %s", line)
			}
			inp := strings.Join(chunks[2:], "")
			bindata, err := EncodePrefixInstructions(instruction, inp)
			if err != nil {
				return err
			}
			buf = append(buf, bindata...)
			continue
		} else {
			opcodeString := chunks[0]

			switch opcodeString {
			case "LD":
				inp := strings.Join(chunks[1:], "")
				bindata, err := EncodeLoad(inp)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "LDH":
				inp := strings.Join(chunks[1:], "")
				bindata, err := EncodeLoadH(inp)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "PUSH":
				if len(chunks) != 2 {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodePush(chunks[1])
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "POP":
				if len(chunks) != 2 {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodePop(chunks[1])
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "NOP", "STOP", "HALT", "PREFIX", "EI", "DI":
				if len(chunks) > 2 {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodeMisc(chunks)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "ADD", "ADC", "SUB", "SBC", "AND", "XOR", "OR", "CP", "DAA", "SCF", "CPL", "CCF", "INC", "DEC":
				inp := strings.Join(chunks[1:], "")
				bindata, err := EncodeBitArithmetic(chunks[0], inp)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
			case "JR":
				var inp string
				if len(chunks) == 3 {
					inp = chunks[1] + chunks[2]
				} else if len(chunks) == 2 {
					inp = chunks[1]
				} else {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodeRelativeJump(inp)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "JP":
				var inp string
				if len(chunks) == 3 {
					inp = chunks[1] + chunks[2]
				} else if len(chunks) == 2 {
					inp = chunks[1]
				} else {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodeJump(inp)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "RET", "RETI":
				if len(chunks) > 2 {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodeReturn(chunks)
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "RST":
				if len(chunks) != 2 {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodeRestart(chunks[1])
				if err != nil {
					return err
				}
				buf = append(buf, bindata...)
				break
			case "RLCA", "RLA", "RRCA", "RRA":
				if len(chunks) != 1 {
					return fmt.Errorf("Unexpected split of input line: %s\n", line)
				}
				bindata, err := EncodeRotateNoPrefix(chunks[0])
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
