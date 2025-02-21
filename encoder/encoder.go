package encoder

import (
	"fmt"
	"slices"
	"strings"
	"strconv"
)

func MSB(n uint16) uint8 {
	return uint8(n >> 8)
}

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

func IsPrefix(s string) bool {
	return slices.Contains(PREFIXES, s)
}

func Is8BitRegister(s string) bool {
	return slices.Contains(BIT8REG, s)
}

func encodeLoadH(argStr string) ([]uint8, error) {
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) != 2 {
		return []uint8{}, fmt.Errorf("could not parse LD arguments, expected only one `,`: %s", argStr)
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

func encodeLoad(argStr string) ([]uint8, error) {
	var opcode uint8
	args := strings.Split(strings.ReplaceAll(strings.TrimSpace(argStr), " ", ""), ",")
	if len(args) != 2 {
		return []uint8{}, fmt.Errorf("could not parse LDH arguments, expected only one `,`: %s", argStr)
	}

	if (args[0] == args[1]) && (args[0] == "[HL]") {
		return []uint8{}, fmt.Errorf("Illegal command: LD [HL],[HL]")
	}

	if Is8BitRegister(args[0]) && Is8BitRegister(args[1]) {
		opcode = 0b01000000 | (REGMAPPER[args[0]] << 3) | (REGMAPPER[args[1]])
		return []uint8{opcode}, nil
	}

	if slices.Contains(REGP, args[0]) {
		if args[1][:2] == "0x" && len(args[1]) >= 3 && len(args[1]) <= 6 {
			nn, err := strconv.ParseUint(args[1][2:], 16, 16)
			n := uint16(nn)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[1], err)
			}
			opcode = REGMAPPER[args[0]] << 3 | 0b110
			return []uint8{opcode, LSB(n), MSB(n)}, nil
		}
	}

	if args[1] == "LP" && args[0][:2] == "0x" && len(args[0]) >= 3 && len(args[0]) <= 6 {
		nn, err := strconv.ParseUint(args[0][2:], 16, 16)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[0], err)
		}
		return []uint8{0x08, LSB(n), MSB(n)}, nil
	}

	if args[1] == "A" && args[0][:2] == "0x" && len(args[0]) >= 3 && len(args[0]) <= 6 {
		nn, err := strconv.ParseUint(args[0][2:], 16, 16)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[0], err)
		}
		return []uint8{0xEA, LSB(n), MSB(n)}, nil
	}

	if args[0] == "A" && args[1][:2] == "0x" && len(args[1]) >= 3 && len(args[1]) <= 6 {
		nn, err := strconv.ParseUint(args[1][2:], 16, 16)
		n := uint16(nn)
		if err != nil {
			return []uint8{}, fmt.Errorf("Could not parse immediate 16bit hex value %s: %w", args[1], err)
		}
		return []uint8{0xFA, LSB(n), MSB(n)}, nil
	}

	if len(args[1]) == 3 || len(args[1]) == 4 {
		if args[1][:2] == "0x" {
			n, err := strconv.ParseUint(args[1][2:], 16, 8)
			if err != nil {
				return []uint8{}, fmt.Errorf("Could not parse immediate 8bit hex value %s: %w", args[1], err)
			}
			opcode = REGMAPPER[args[0]] << 3 | 0b110
			return []uint8{opcode, uint8(n)}, nil
		}
	}
	if args[0] == "HL" && args[1][:3] == "SP+" && len(args[1][3:]) >= 3 && len(args[1][3:]) <= 4 {
		if args[1][3:5] != "0x" {
			return []uint8{}, fmt.Errorf("Invalid arguments passed: %s", args[0]+","+args[1])
		}
		n, err := strconv.ParseUint(args[1][5:], 16, 8)
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

func encodePush(argStr string) ([]uint8, error) {
	var opcode uint8
	if slices.Contains(REGPF, argStr) {
		opcode = 0b11000101 | (REGMAPPER[argStr] << 4)
		return []uint8{opcode}, nil
	}

	return []uint8{}, fmt.Errorf("Invalid push instruction: PUSH %s", argStr)
}

func encodePop(argStr string) ([]uint8, error) {
	var opcode uint8
	if slices.Contains(REGPF, argStr) {
		opcode = 0b11000001 | (REGMAPPER[argStr] << 4)
		return []uint8{opcode}, nil
	}

	return []uint8{}, fmt.Errorf("Invalid push instruction: POP %s", argStr)
}

