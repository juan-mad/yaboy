package encoder

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
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
	"BC":   0,
	"DE":   1,
	"HL":   2,
	"SP":   3,
	"AF":   3,
}

var ALUMAPPER = map[string]uint8{
	"ADD": 0,
	"ADC": 1,
	"SUB": 2,
	"SBC": 3,
	"AND": 4,
	"XOR": 5,
	"OR":  6,
	"CP":  7,
	"INC": 4,
	"DEC": 5,
}

// opcode = (2 << 6) | (ALUMAPPER[instr] << 3) | REGMAPPER[rightop]
// 0x84
// 0b100000100
var ALUUNARYOP = []string{"DAA", "SCF", "CPL", "CCF"}
var FLAGS = []string{"NZ", "Z", "NC", "C"}
var FLAGMAPPER = map[string]uint8{
	"NZ": 0,
	"Z":  1,
	"NC": 2,
	"C":  3,
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
			opcode = REGMAPPER[args[0]]<<3 | 0b110
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
		// if it fits in 8 bits, we encode as 0x3E, if not, as 0xFA
		if (n >> 8) == 0 {
			return []uint8{0x3E, uint8(n)}, nil
		} else {
			return []uint8{0xFA, LSB(n), MSB(n)}, nil
		}
	}

	if len(args[1]) == 3 || len(args[1]) == 4 {
		// try parsing 8bit value
		if args[1][:2] == "0x" {
			n, err := strconv.ParseUint(args[1][2:], 16, 8)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[1], err)
			}
			opcode = REGMAPPER[args[0]]<<3 | 0b110
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
	if len(chunks) > 2 || len(chunks) == 0 {
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

// encodeBitArithmetic produces the bytes for arithmetic instructions
func encodeBitArithmetic(instr string, argStr string) ([]uint8, error) {
	var opcode uint8

	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")

	if len(args) > 2 {
		return []uint8{}, fmt.Errorf("Unexpected number (%d) of arguments: %s", len(args), args)
	}

	// if 2 elems but first operand is not A, error
	if len(args) == 1 && ((args[0] == "A") && (args[0] != "HL")) {
		return []uint8{}, fmt.Errorf("Invalid ALU operation %s, register A or HL must be first operand: %s", instr, args)
	}

	// if no operands, it must be one of these
	if len(args) == 0 {
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
		default:
			// should never happen
			return []uint8{}, fmt.Errorf("Unknown use of encodeBitArithmetic: %s %s", instr, args)
		}
	}

	// if more than one element, check that it has no operands, or error
	if slices.Contains(ALUUNARYOP, instr) {
		return []uint8{}, fmt.Errorf("Invalid instruction %s with operands: %s", instr, args)
	}

	// INC and DEC instructions that take exactly one arguments
	if len(args) == 1 && (instr == "INC" || instr == "DEC") {
		if Is8BitRegister(args[0]) {
			opcode = (REGMAPPER[args[0]] << 3) | ALUMAPPER[instr]
			return []uint8{opcode}, nil
			// return []uint8{}, fmt.Errorf("Invalid instruction %s with no 8bit reg: %s", instr, chunks[1])
		} else if slices.Contains(REGP, args[0]) {
			if instr == "INC" {
				opcode = (REGMAPPER[args[0]] << 4) | 0x3
			} else if instr == "DEC" { // DEC
				opcode = (REGMAPPER[args[0]] << 4) | 0xB
			}
			return []uint8{opcode}, nil
		}
	}

	if len(args) == 2 && instr == "ADD" && args[0] == "SP" {
		if !(args[1][:2] == "0x" && len(args[1]) >= 3 && len(args[1]) <= 4) {
			return []uint8{}, fmt.Errorf("Invalid ADD argument: %s", args)
		}
		nn, err := strconv.ParseUint(args[1][2:], 16, 8)
		n := uint8(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args, err)
		}
		return []uint8{0xE8, n}, nil
	}

	// at this point, there should be no DEC or INC
	if instr == "DEC" || instr == "INC" {
		return []uint8{}, fmt.Errorf("Invalid %s instruction: %s", instr, args)
	}

	// rest of operations may ommit A as first operand, we just want the second one
	rightop := args[len(args)-1]

	// operation between operands
	if Is8BitRegister(rightop) {
		opcode = (2 << 6) | (ALUMAPPER[instr] << 3) | REGMAPPER[rightop]
		return []uint8{opcode}, nil
	} else if args[0] == "HL" {
		if !slices.Contains(REGP, args[1]) {
			return []uint8{}, fmt.Errorf("Invalid ADD arguments: %s", args)
		}
		opcode = (REGMAPPER[args[1]] << 4) | 0x9
		return []uint8{opcode}, nil
	} else {
		opcode = (3 << 6) | (ALUMAPPER[instr] << 3) | 0x6
		fmt.Println("aaaa")
		fmt.Printf("%#02x\n", opcode)
		if !(rightop[:2] == "0x" && len(rightop) >= 3 && len(rightop) <= 4) {
			return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", args)
		}
		nn, err := strconv.ParseUint(rightop[2:], 16, 8)
		n := uint8(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", rightop, err)
		}
		return []uint8{opcode, n}, nil
	}
}

// encodeRelativeJump produces the bytes for relative jump instructions, JR
func encodeRelativeJump(argStr string) ([]uint8, error) {
	var opcode uint8
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) == 0 || len(args) > 2 {
		return []uint8{}, fmt.Errorf("could not parse JR arguments: %s", argStr)
	}

	hexStr := args[len(args)-1]
	if !(hexStr[:2] == "0x" && len(hexStr) >= 3 && len(hexStr) <= 4) {
		return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", argStr)
	}
	nn, err := strconv.ParseUint(hexStr[2:], 16, 8)
	n := uint8(nn)
	if err != nil {
		return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", argStr, err)
	}

	if len(args) == 1 {
		return []uint8{0x18, n}, nil
	}

	opcode = (1 << 5) | (FLAGMAPPER[args[1]] << 3)
	return []uint8{opcode, n}, nil

}

// encodeJump produces the bytes for jump instructions, JP
func encodeJump(argStr string) ([]uint8, error) {
	var opcode uint8
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) == 0 || len(args) > 2 {
		return []uint8{}, fmt.Errorf("could not parse JR arguments: %s", argStr)
	}

	if len(args) == 1 {
		if args[1] == "HL" {
			return []uint8{0xE9}, nil
		}
	}

	hexStr := args[len(args)-1]
	if !(hexStr[:2] == "0x" && len(hexStr) >= 3 && len(hexStr) <= 4) {
		return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", argStr)
	}
	nn, err := strconv.ParseUint(hexStr[2:], 16, 8)
	n := uint8(nn)
	if err != nil {
		return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", argStr, err)
	}

	if len(args) == 1 {
		return []uint8{0xC3, n}, nil
	}

	if !slices.Contains(FLAGS, args[1]) {
		return []uint8{}, fmt.Errorf("Invalid JP instruction: JP %s", argStr)
	}

	opcode = 0b11000010 | (FLAGMAPPER[args[1]] << 3)
	return []uint8{opcode, n}, nil

}

// encodeReturn produces the bytes for return instructions, RET and RETI
func encodeReturn(chunks []string) ([]uint8, error) {
	var opcode uint8
	instr := chunks[0]
	if len(chunks) > 2 {
		return []uint8{}, fmt.Errorf("too many inputs for encodeReturn function")
	}
	if instr != "RET" && instr != "RETI" {
		return []uint8{}, fmt.Errorf("invalid input for encodeReturn function: %s", chunks)
	}

	if instr == "RETI" {
		if len(chunks) > 1 {
			return []uint8{}, fmt.Errorf("Invalid RETI instruction, takes no operands, found: %s", chunks[1:])
		}
		return []uint8{0xD9}, nil
	}

	if instr == "RET" {
		if len(chunks) == 1 {
			return []uint8{0xC9}, nil
		}

		if !slices.Contains(FLAGS, chunks[1]) {
			return []uint8{}, fmt.Errorf("Invalid operand for RET, not a flag: %s", chunks[1])
		}

		opcode = 0b11000000 | (FLAGMAPPER[chunks[1]] << 3)
		return []uint8{opcode}, nil
	}

	return []uint8{}, fmt.Errorf("Unknown use of encodeReturn: %s", chunks)
}

// encodeRestart produces the bytes for restart instructions, RST
func encodeRestart(argStr string) ([]uint8, error) {

	if !(argStr[:2] == "0x" && len(argStr) >= 3 && len(argStr) <= 4) {
		return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", argStr)
	}
	nn, err := strconv.ParseUint(argStr, 16, 8)
	n := uint8(nn)
	if err != nil {
		return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", argStr, err)
	}

	switch n {
	case 0x00:
		return []uint8{0xC7}, nil
	case 0x08:
		return []uint8{0xCF}, nil
	case 0x10:
		return []uint8{0xD7}, nil
	case 0x18:
		return []uint8{0xDF}, nil
	case 0x20:
		return []uint8{0xE7}, nil
	case 0x28:
		return []uint8{0xEF}, nil
	case 0x30:
		return []uint8{0xF7}, nil
	case 0x38:
		return []uint8{0xFF}, nil
	default:
		return []uint8{}, fmt.Errorf("Invalid operand for RST instruction: %s", argStr)
	}

}

// encodeCall produces the bytes for call instructions, CALL
func encodeCall(argStr string) ([]uint8, error) {
	var opcode uint8
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) == 0 || len(args) > 2 {
		return []uint8{}, fmt.Errorf("could not parse JR arguments: %s", argStr)
	}

	hexStr := args[len(args)-1]
	if !(hexStr[:2] == "0x" && len(hexStr) >= 3 && len(hexStr) <= 6) {
		return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", argStr)
	}
	nn, err := strconv.ParseUint(hexStr[2:], 16, 16)
	n := uint16(nn)
	if err != nil {
		return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", argStr, err)
	}

	if len(args) == 1 {
		return []uint8{0xCD, LSB(n), MSB(n)}, nil
	}

	if !slices.Contains(FLAGS, args[1]) {
		return []uint8{}, fmt.Errorf("Invalid JP instruction: JP %s", argStr)
	}

	opcode = 0b11000100 | (FLAGMAPPER[args[1]] << 3)
	return []uint8{opcode, LSB(n), MSB(n)}, nil
}

func encodeRotateNoPrefix(instr string) ([]uint8, error) {
	switch instr {
	case "RLCA":
		return []uint8{0x07}, nil
	case "RLA":
		return []uint8{0x17}, nil
	case "RRCA":
		return []uint8{0x0F}, nil
	case "RRA":
		return []uint8{0x1F}, nil
	default:
		return []uint8{}, fmt.Errorf("unknown use of encodeRotateNoPrefix: %s", instr)
	}
}
