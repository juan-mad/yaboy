package decompiler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	enc "yaboy/encoder"
)

var ENCODE_R8 = map[uint8]string{
	0: "B", 1: "C", 2: "D", 3: "E", 4: "H", 5: "L", 6: "[HL]", 7: "A",
}
var ENCODE_R16 = map[uint8]string{
	0: "BC", 1: "DE", 2: "HL", 3: "SP",
}

var ENCODE_R16STK = map[uint8]string{
	0: "BC", 1: "DE", 2: "HL", 3: "AF",
}
var ENCODE_R16MEM = map[uint8]string{
	0: "BC", 1: "DE", 2: "HL+", 3: "HL-",
}
var ENCODE_COND = map[uint8]string{
	0: "NZ", 1: "Z", 2: "NC", 3: "C",
}

var BLOCK2_BITES = map[uint8]string{
	0b000: "ADD",
	0b001: "ADC",
	0b010: "SUB",
	0b011: "SBC",
	0b100: "AND",
	0b101: "XOR",
	0b110: "OR",
	0b111: "CP",
}

// formatOutput formats the different pieces of an instruction into a single string
func formatOutput(ss []string, flip []int8) string {
	if len(flip) == 2 {
		aux := ss[flip[0]]
		ss[flip[0]] = ss[flip[1]]
		ss[flip[1]] = aux
	}
	if len(ss) == 2 {
		return fmt.Sprintf(ss[0], ss[1])
	}
	if len(ss) == 3 {
		return fmt.Sprintf(ss[0], ss[1], ss[2])
	}

	if len(ss) == 4 {
		return fmt.Sprintf(ss[0], ss[1], ss[2], ss[3])
	}
	return "MAL"
}

func decodePrefixOpcode(b uint8) (string, error) {
	reg := ENCODE_R8[b&0b111]
	top := b >> 6
	mid := (b >> 3) & 0b111

	switch top {
	case 0b00:
		switch mid {
		case 0b000:
			return fmt.Sprintf("RLC %s", reg), nil
		case 0b001:
			return fmt.Sprintf("RRC %s", reg), nil
		case 0b010:
			return fmt.Sprintf("RL %s", reg), nil
		case 0b011:
			return fmt.Sprintf("RR %s", reg), nil
		case 0b100:
			return fmt.Sprintf("SLA %s", reg), nil
		case 0b101:
			return fmt.Sprintf("SRA %s", reg), nil
		case 0b110:
			return fmt.Sprintf("SWAP %s", reg), nil
		case 0b111:
			return fmt.Sprintf("SRL %s", reg), nil
		}
	case 0b01:
		return fmt.Sprintf("BIT %d,%s", mid, reg), nil
	case 0b10:
		return fmt.Sprintf("RES %d,%s", mid, reg), nil
	case 0b11:
		return fmt.Sprintf("SET %d,%s", mid, reg), nil
	}

	return "", nil
}

// DecodeBytes takes a slice of bytes and returns a slice of strings corresponding to the
// bytes passed.
func DecodeBytes(data []uint8) ([]string, error) {
	output := []string{}
	unfInstr := []string{}
	getNextBytes := 0
	flipPos := []int8{}
	getNextPrefixByte := false
	for i, b := range data {

		if getNextPrefixByte {
			instr, err := decodePrefixOpcode(b)
			if err != nil {
				return []string{}, fmt.Errorf("Error decompiling prefix opcode %02x: %w", b, err)
			}
			output = append(output, instr)
			getNextPrefixByte = false
			continue
		}

		if getNextBytes > 0 {
			unfInstr = append(unfInstr, strconv.FormatUint(uint64(b), 16))
			getNextBytes = getNextBytes - 1

			if getNextBytes == 0 {
				output = append(output, formatOutput(unfInstr, flipPos))
				unfInstr = []string{}
				flipPos = []int8{}
			}
			continue
		}

		switch {
		case b == 0xCB:
			getNextPrefixByte = true
			break
		// MISC OPERATIONS
		case b == 0x00:
			output = append(output, "NOP")
			break
		case b == 0x10:
			unfInstr = []string{"STOP 0x%02s"}
			getNextBytes = 1
			break
		case b == 0x76:
			output = append(output, "HALT")
			break
		case b == 0xF3:
			output = append(output, "DI")
			break
		case b == 0xFB:
			output = append(output, "EI")
			break
			// BLOCK 1: 0b01xxxxxx, LD r8 r8
			// Notice we have removed LD [HL],[HL] from reaching this point, as it is HALT
		case (b >= 0x40) && (b <= 0x7F):
			dst := ENCODE_R8[(b>>3)&0b111]
			src := ENCODE_R8[b&0b111]
			output = append(output, fmt.Sprintf("LD %s,%s", dst, src))
			break
			// More LD operations
			// LD [r16],A
		case ((b >> 6) == 0) && ((b & 0b1111) == 0b0010):
			dst := ENCODE_R16MEM[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("LD [%s],A", dst))
			break
			// LD A,[r16]
		case ((b >> 6) == 0) && ((b & 0b1111) == 0b1010):
			dst := ENCODE_R16MEM[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("LD A,[%s]", dst))
			break
			// LD [n16],SP
		case b == 0b00001000:
			unfInstr = []string{"LD [0x%02s%02s],SP"}
			getNextBytes = 2
			flipPos = []int8{1, 2}
			break
			// LD r16,n16
		case ((b >> 6) == 0) && ((b & 0b1111) == 0b0001):
			dst := ENCODE_R16[(b>>4)&0b1111]
			unfInstr = []string{"LD %s,0x%02s%02s", dst}
			getNextBytes = 2
			flipPos = []int8{2, 3}
			break
			// INC r16
		case ((b >> 6) == 0) && ((b & 0b1111) == 0b0011):
			dst := ENCODE_R16[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("INC %s", dst))
			break
			// DEC r16
		case ((b >> 6) == 0) && ((b & 0b1111) == 0b1011):
			dst := ENCODE_R16[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("DEC %s", dst))
			break
			// ADD HL,r16
		case ((b >> 6) == 0) && ((b & 0b1111) == 0b1001):
			dst := ENCODE_R16[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("ADD HL,%s", dst))
			break
			// INC r8
		case ((b >> 6) == 0) && ((b & 0b111) == 0b100):
			dst := ENCODE_R8[(b>>3)&0b111]
			output = append(output, fmt.Sprintf("INC %s", dst))
			break
			// DEC r8
		case ((b >> 6) == 0) && ((b & 0b111) == 0b101):
			dst := ENCODE_R8[(b>>3)&0b111]
			output = append(output, fmt.Sprintf("DEC %s", dst))
			break
			// LD r8,n8
		case ((b >> 6) == 0) && ((b & 0b111) == 0b110):
			dst := ENCODE_R8[(b>>3)&0b111]
			unfInstr = []string{"LD %s,0x%02s", dst}
			getNextBytes = 1
			break
		case b == 0x07:
			output = append(output, "RLCA")
			break
		case b == 0x17:
			output = append(output, "RLA")
			break
		case b == 0x0F:
			output = append(output, "RRCA")
			break
		case b == 0x1F:
			output = append(output, "RRA")
			break
		case b == 0x27:
			output = append(output, "DAA")
			break
		case b == 0x37:
			output = append(output, "SCF")
			break
		case b == 0x2F:
			output = append(output, "CPL")
			break
		case b == 0x3F:
			output = append(output, "CCF")
			break
		case b == 0x18:
			unfInstr = []string{"JR 0x%02s"}
			getNextBytes = 1
			break
		case ((b >> 5) == 0b001) && ((b & 0b111) == 0):
			cond := ENCODE_COND[(b>>3)&0b11]
			unfInstr = []string{"JR %s,0x%02s", cond}
			getNextBytes = 1
			break
			// BLOCK 2
		case (b >> 6) == 0b10:
			instr := BLOCK2_BITES[(b>>3)&0b111]
			reg := ENCODE_R8[(b & 0b111)]
			output = append(output, fmt.Sprintf("%s A,%s", instr, reg))
			break
			// BLOCK 3
			// ADD A,n16
		case (b>>6) == 0b11 && ((b & 0b111) == 0b110):
			instr := BLOCK2_BITES[(b>>3)&0b111]
			unfInstr = []string{"%s A,0x%02s", instr}
			getNextBytes = 1
			break
			// RET COND
		case (b>>5) == 0b110 && ((b & 0b111) == 0b000):
			cond := ENCODE_COND[(b>>3)&0b11]
			output = append(output, fmt.Sprintf("RET %s", cond))
			break
			// RET
		case b == 0b11001001:
			output = append(output, "RET")
			break
			// RETI
		case b == 0b11011001:
			output = append(output, "RETI")
			break
			// JP COND,n16
		case ((b >> 5) == 0b110) && ((b & 0b111) == 0b010):
			cond := ENCODE_COND[(b>>3)&0b11]
			unfInstr = []string{"JP %s,0x%02s%02s", cond}
			getNextBytes = 2
			flipPos = []int8{2, 3}
			break
			// JP n16
		case b == 0b11000011:
			unfInstr = []string{"JP 0x%02s%02s"}
			getNextBytes = 2
			flipPos = []int8{1, 2}
			break
			// JP HL
		case b == 0b11101001:
			output = append(output, "JP HL")
			break
			// CALL COND,n16
		case (b>>5) == 0b110 && ((b & 0b111) == 0b100):
			cond := ENCODE_COND[(b>>3)&0b11]
			unfInstr = []string{"CALL %s,0x%02s%02s", cond}
			getNextBytes = 2
			flipPos = []int8{2, 3}
			break
			// CALL n16
		case b == 0b11001101:
			unfInstr = []string{"CALL 0x%02s%02s"}
			getNextBytes = 2
			flipPos = []int8{1, 2}
			break
			// RST
		case (b>>6) == 0b11 && ((b & 0b111) == 0b111):
			target := strconv.FormatUint((uint64((b>>3)&0b111) * 8), 16)
			output = append(output, fmt.Sprintf("RST 0x%02s", target))
			break
			// POP
		case (b>>6) == 0b11 && ((b & 0b1111) == 0b0001):
			reg := ENCODE_R16STK[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("POP %s", reg))
			break
			// PUSH
		case (b>>6) == 0b11 && ((b & 0b1111) == 0b0101):
			reg := ENCODE_R16STK[(b>>4)&0b11]
			output = append(output, fmt.Sprintf("PUSH %s", reg))
			break
		case b == 0b11100010:
			output = append(output, "LDH [C],A")
			break
		case b == 0b11100000:
			unfInstr = []string{"LDH [0x%02s],A"}
			getNextBytes = 1
			break
		case b == 0b11101010:
			unfInstr = []string{"LD [0x%02s%02s],A"}
			getNextBytes = 2
			flipPos = []int8{1, 2}
			break
		case b == 0b11110010:
			output = append(output, "LDH A,[C]")
			break
		case b == 0b11110000:
			unfInstr = []string{"LDH A,[0x%02s]"}
			getNextBytes = 1
			break
		case b == 0b11111010:
			unfInstr = []string{"LD A,[0x%02s%02s]"}
			getNextBytes = 2
			flipPos = []int8{1, 2}
			break
		case b == 0b11101000:
			unfInstr = []string{"ADD SP,0x%02s"}
			getNextBytes = 1
			break
		case b == 0b11111000:
			unfInstr = []string{"LD HL,SP+0x%02s"}
			getNextBytes = 1
			break
		case b == 0b11111001:
			output = append(output, "LD SP,HL")
			break
		case b == 0b11110011:
			output = append(output, "DI")
			break
		case b == 0b11111011:
			output = append(output, "EI")
			break
		default:
			fmt.Printf("Warning: Unmapped byte %#02x at pos %#04x\n", b, i)
			break
		}
	}
	return output, nil

}

func decompileBinary() {
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

	res, err := DecodeBytes(data)
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
