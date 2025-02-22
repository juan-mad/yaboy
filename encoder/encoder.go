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

var ALUMAPPER = map[string]uint8 {
	"ADD": 0,
	"ADC": 1,
	"SUB": 2,
	"SBC": 3,
	"AND": 4,
	"XOR": 5,
	"OR": 6,
	"CP": 7,
	"INC": 4,
	"DEC": 5,
}

var ALUUNARYOP = []string{"DAA", "SCF", "CPL", "CCF",}

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
			n := uint8(nn)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[1], err)
			}
			return []uint8{0xF0, n}, nil
		}
	}

	if args[1] == "A" {
		if args[0][:2] == "0x" && len(args[0]) >= 3 && len(args[0]) <= 4 {
			nn, err := strconv.ParseUint(args[0][2:], 16, 8)
			n := uint8(nn)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[0], err)
			}
			return []uint8{0xE0, n}, nil
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
		n := uint8(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", chunks[1], err)
		}
		return []uint8{0x10, n}, nil
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


	return []uint8{}, fmt.Errorf("Invalid push instruction: %s", chunks)
}
func encodeBitArithmetic(chunks []string) ([]uint8, error) {
	var opcode uint8
	var instr string = chunks[0]
	var rightop string
	
	// check that chunks has 1, 2, or 3 elements
	if len(chunks) < 1 || len(chunks) > 4 {
		return []uint8{}, fmt.Errorf("Unexpected number of arguments: %s", chunks)
	}
	// if 3 elems but first operand is not A, error
	if len(chunks) == 3 && chunks[1] != "A" {
		return []uint8{}, fmt.Errorf("Invalid 8bit ALU operation, A register must be first operand: %s", chunks)
	}

	// if only one element, it must be an instruction with no operands
	if len(chunks) == 1 {
		if !slices.Contains(ALUUNARYOP, instr) {
			return []uint8{}, fmt.Errorf("Invalid %s instruction with no operands", instr)
		}
		switch instr {
			case "DAA":
				return []uint8{0x27}, nil
			case "SCF":
				return []uint8{0x37}, nil
			case "CPL":
				return []uint8{0x2F}, nil
			case "CCF":
				return []uint8{0x3F}, nil
		}
	}

	// if more than one element, check that it has no operands, or error
	if slices.Contains(ALUUNARYOP, instr) {
		return []uint8{}, fmt.Errorf("Invalid instruction %s with operands: %s", instr, chunks[1:])
	}

	// INC and DEC instructions that take exactly one arguments
	if len(chunks) == 2 && (instr == "INC" || instr == "DEC" || instr == "ADD") {
		if Is8BitRegister(chunks[1]) {
			opcode = (REGMAPPER[chunks[1]] << 3) | ALUMAPPER[instr]
			return []uint8{opcode}, nil
			// return []uint8{}, fmt.Errorf("Invalid instruction %s with no 8bit reg: %s", instr, chunks[1])
		} else if slices.Contains(REGP, chunks[1]) {
			if instr == "INC" {
				opcode = (REGMAPPER[chunks[1]] << 4) | 0x3
			} else if instr == "DEC" { // DEC
				opcode = (REGMAPPER[chunks[1]] << 4) | 0xB
			} else { // ADD
				opcode = (REGMAPPER[chunks[1]] << 4) | 0x9
			}
			return []uint8{opcode}, nil
		}
	}
	if len(chunks) == 3 && instr == "ADD" && chunks[1] == "SP" {
		if !(chunks[2][:2] == "0x" && len(chunks[2]) >= 3 && len(chunks[2]) <= 4) {
			return []uint8{}, fmt.Errorf("Invalid ADD argument: %s", chunks)
		}
		nn, err := strconv.ParseInt(rightop[2:], 16, 8)
		n := uint8(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", rightop, err)
		}
		return []uint8{0xE8, n}, nil
	}

	// at this point, there should be no DEC or INC
	if instr == "DEC" || instr == "INC" {
		return []uint8{}, fmt.Errorf("Invalid %s instruction: %s", instr, chunks)
	}
	
	// rest of operations may ommit A as first operand, we just want the second one
	rightop = chunks[len(chunks)-1]

	// operation between operands
	if Is8BitRegister(rightop) {
		 opcode = (2 << 6) | (ALUMAPPER[instr] << 3) | REGMAPPER[rightop]
		 return []uint8{opcode}, nil
	} else {
		opcode = (3 << 6) | ALUMAPPER[instr] | 0x6
		if !(rightop[:2] == "0x" && len(rightop) >= 3 && len(rightop) <= 4) {
			return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", chunks)
		}
		nn, err := strconv.ParseUint(rightop[2:], 16, 8)
		n := uint8(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", rightop, err)
		}
		return []uint8{opcode, n}, nil
	}
}

