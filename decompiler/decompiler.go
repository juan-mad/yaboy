package decompiler

import "fmt"
import "strconv"

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

// DecodeBytes takes a slice of bytes and returns a slice of strins corresponding to the
// bytes passed.
func DecodeBytes(data []uint8) ([]string, error) {
	output := []string{}
	unfInstr := []string{}
	getNextBytes := 0
	flipPos := []int8{}

	for _, b := range data {
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
		case b == 0xCB:
			output = append(output, "PREFIX")
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

		}

	}
	return output, nil

}
