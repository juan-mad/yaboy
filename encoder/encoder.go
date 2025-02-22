package encoder

import (
	"fmt"
	"slices"
	"strings"
	"strconv"
)

// MSB returns the Most Significant Bits of an unsigned 16bit value as unsigned 8bit value
func MSB(n uint16) uint8 {
	return uint8(n >> 8)
}

// LSB returns the Least Significant Bits of an unsigned 16bit value as unsigned 8bit value
func LSB(n uint16) uint8 {
	return uint8(0x00FF & n)
}

var PREFIXES = []string{"CB", "DD", "ED", "FD"}
var BIT8REG = []string{"B", "C", "D", "E", "H", "L", "[HL]", "A"}
var REGP = []string{"BC", "DE", "HL", "SP"}
var REGPF = []string{"BC", "DE", "HL", "AF"}

var BIT16REG = []string{"BC", "DE", "HL", "SP", "AF"}

var REGMAPPER = map[string]uint8{
	"B":    0,
	"C":    1,
	"D":    2,
	"E":    3,
	"H":    4,
	"L":    5,
	"[HL]": 6,
	"A":    7,
	"BC": 0,
	"DE": 1,
	"HL": 2,
	"SP": 3,
	"AF": 3,
}

// IsPrefix checks that a string represents one of the prefixes
func IsPrefix(s string) bool {
	return slices.Contains(PREFIXES, s)
}

// Is8BitRegister checks if a string represents one of the 8 bit registries
func Is8BitRegister(s string) bool {
	return slices.Contains(BIT8REG, s)
}

// encodeLoadH produces the bytes for the LDH operation, related to a subset of load operations
// regarding the Accumulator registry A
func encodeLoadH(argStr string) ([]uint8, error) {
	// Function expects that the string can be split in two sides separated by a comma
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) != 2 {
		return []uint8{}, fmt.Errorf("could not parse LDH arguments, expected only one `,`: %s", argStr)
	}
	if args[0] == "[C]" && args[1] == "A" {
		return []uint8{0xE2}, nil
	}
	if args[0] == "A" && args[1] == "[C]" {
		return []uint8{0xF2}, nil
	}

	if args[0] == "A" {
		if args[1][:2] == "0x" && len(args[1]) >= 3 && len(args[1]) <= 4 {
			nn, err := strconv.ParseUint(args[1][2:], 16, 8)
			n := uint16(nn)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[1], err)
			}
			return []uint8{0xF0, uint8(n)}, nil
		}
	}

	if args[1] == "A" {
		if args[0][:2] == "0x" && len(args[0]) >= 3 && len(args[0]) <= 4 {
			nn, err := strconv.ParseUint(args[0][2:], 16, 8)
			n := uint16(nn)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[0], err)
			}
			return []uint8{0xE0, uint8(n)}, nil
		}
	}

	return []uint8{}, fmt.Errorf("Not implemented load instruction: LDH %s", argStr)
}

// encodeLoad produces the bytes for the LD operation, related to load operations
func encodeLoad(argStr string) ([]uint8, error) {
	var opcode uint8
	// Function expects that the string can be split in two sides separated by a comma
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) != 2 {
		return []uint8{}, fmt.Errorf("could not parse LD arguments, expected only one `,`: %s", argStr)
	}

	// illegal combination of registries
	if (args[0] == args[1]) && (args[0] == "[HL]") {
		return []uint8{}, fmt.Errorf("Illegal command: LD [HL],[HL]")
	}

	// programmatic LD operation between two 8bit registers
	if Is8BitRegister(args[0]) && Is8BitRegister(args[1]) {
		opcode = 0b01000000 | (REGMAPPER[args[0]] << 3) | (REGMAPPER[args[1]])
		return []uint8{opcode}, nil
	}

	// If the first operand is one of the 16bit registers BC, DE, HL, SP
	if slices.Contains(REGP, args[0]) { // try to parse 16bit value
		if args[1][:2] == "0x" && len(args[1]) >= 3 && len(args[1]) <= 6 {
			nn, err := strconv.ParseUint(args[1][2:], 16, 16)
			n := uint16(nn)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[1], err)
			}
			opcode = REGMAPPER[args[0]] << 3 | 0b110
			// little endian order, LSB before MSB
			return []uint8{opcode, LSB(n), MSB(n)}, nil
		}
	}

	if args[1] == "LP" && args[0][:2] == "0x" && len(args[0]) >= 3 && len(args[0]) <= 6 {
		// try parsing 16bit value
		nn, err := strconv.ParseUint(args[0][2:], 16, 16)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[0], err)
		}
		return []uint8{0x08, LSB(n), MSB(n)}, nil
	}

	if args[1] == "A" && args[0][:2] == "0x" && len(args[0]) >= 3 && len(args[0]) <= 6 {
		// try parsing 16bit value
		nn, err := strconv.ParseUint(args[0][2:], 16, 16)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[0], err)
		}
		return []uint8{0xEA, LSB(n), MSB(n)}, nil
	}

	if args[0] == "A" && args[1][:2] == "0x" && len(args[1]) >= 3 && len(args[1]) <= 6 {
		// try parsing 16bit value
		nn, err := strconv.ParseUint(args[1][2:], 16, 16)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[1], err)
		}
		return []uint8{0xFA, LSB(n), MSB(n)}, nil
	}

	if len(args[1]) == 3 || len(args[1]) == 4 {
		// try parsing 8bit value
		if args[1][:2] == "0x" {
			n, err := strconv.ParseUint(args[1][2:], 16, 8)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[1], err)
			}
			opcode = REGMAPPER[args[0]] << 3 | 0b110
			return []uint8{opcode, uint8(n)}, nil
		}
	}
	// LD HL,SP+e8 instruction: add signed e8 value to SP and copy the result in HL
	if args[0] == "HL" && args[1][:3] == "SP+" && len(args[1][3:]) >= 3 && len(args[1][3:]) <= 4 {
		//try
		if args[1][3:5] != "0x" {
			return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", args[0]+","+args[1])
		}
		n, err := strconv.ParseInt(args[1][5:], 16, 8)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[1], err)
		}
		return []uint8{0xF8, uint8(n)}, nil
	}
	
	arg := args[0] + "," + args[1]
	switch arg {
	case "[BC],A":
		return []uint8{0x02}, nil
	case "[DE],A":
		return []uint8{0x12}, nil
	case "[HL+],A":
		return []uint8{0x22}, nil
	case "[HL-],A":
		return []uint8{0x32}, nil
	case "A,[BC]":
		return []uint8{0x0A}, nil
	case "A,[DE]":
		return []uint8{0x1A}, nil
	case "A,[HL+]":
		return []uint8{0x2A}, nil
	case "A,[HL-]":
		return []uint8{0x3A}, nil
	case "SP,HL":
		return []uint8{0xF9}, nil
	}

	return []uint8{}, fmt.Errorf("Not implemented load instruction: LD %s", argStr)
}

// encodeIncrement produces the bytes for the INC instructions
func encodeIncrement(argStr string) ([]uint8, error) {
	var opcode uint8

	if Is8BitRegister(argStr) {
		opcode = 0b00000100 | (REGMAPPER[argStr] << 3)
		return []uint8{opcode}, nil
	}

	if slices.Contains(REGP, argStr) {
		opcode = 0b0000011 | (REGMAPPER[argStr] << 4)
	}
	
	return []uint8{}, fmt.Errorf("Invalid increment instruction: INC %s", argStr)
}

// encodePush produces the bytes for the PUSH instructions
func encodePush(argStr string) ([]uint8, error) {
	var opcode uint8
	if slices.Contains(REGPF, argStr) {
		opcode = 0b11000101 | (REGMAPPER[argStr] << 4)
		return []uint8{opcode}, nil
	}

	return []uint8{}, fmt.Errorf("Invalid push instruction: PUSH %s", argStr)
}

// encodePop produces the bytes for the POP instructions
func encodePop(argStr string) ([]uint8, error) {
	var opcode uint8
	if slices.Contains(REGPF, argStr) {
		opcode = 0b11000001 | (REGMAPPER[argStr] << 4)
		return []uint8{opcode}, nil
	}

	return []uint8{}, fmt.Errorf("Invalid push instruction: POP %s", argStr)
}

// encodeMisc produces the bytes for miscellaneous instructions
// in particular, STOP, HALT, PREFIX, NOP, EI, DEI
func encodeMisc(chunks []string) ([]uint8, error) {
	if len(chunks) > 2  || len(chunks) == 0 {
		return []uint8{}, fmt.Errorf("chunks argument must have length 1 or 2, found %d", len(chunks))
	}
	if len(chunks) == 2 {
		if chunks[0] != "STOP" {
			return []uint8{}, fmt.Errorf("Invalid use of %s with arguments: %s", chunks[0], chunks[1])
		}
		if !(chunks[1][:2] == "0x" && len(chunks[1]) >= 3 && len(chunks[1]) <= 4) {
			return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", chunks[0]+","+chunks[1])
		}
		nn, err := strconv.ParseUint(chunks[1][2:], 16, 8)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", chunks[1], err)
		}
		return []uint8{0x10, uint8(n)}, nil
		// case "NOP", "STOP", "HALT", "PREFIX", "EI", "DI":
	}

	instrStr := chunks[0]
	switch instrStr {
		case "NOP":
			return []uint8{0x00}, nil
		case "STOP":
			return []uint8{0x10}, nil
		case "HALT":
			return []uint8{0x76}, nil
		case "PREFIX":
			return []uint8{0xCB}, nil
		case "DI":
			return []uint8{0xF3}, nil
		case "EI":
			return []uint8{0xFB}, nil
	}


	return []uint8{}, fmt.Errorf("Invalid push instruction: ", chunks)
}
func encode8bitArithmetic(chunks []string) ([]uint8, error) {
	var opcode uint8

	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) != 2 {
		return []uint8{}, fmt.Errorf("could not parse LDH arguments, expected only one `,`: %s", argStr)
	}
}

