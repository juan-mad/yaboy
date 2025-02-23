package encoder_test

var ADD_INPUT = `ADD A,B
ADD A,C
ADD A,D
ADD A,E
ADD A,H
ADD A,L
ADD A,[HL]
ADD A,A
ADD H
ADD E
ADD HL,BC
ADD  HL, DE
ADD HL ,HL
ADD HL , SP
ADD A, 0x7E
ADD SP ,0xDF
`

var ADC_INPUT = `ADC A,B
ADC A,C
ADC A,D
ADC A,E
ADC A,H
ADC A,L
ADC A,[HL]
ADC A,A
ADC A, 0x5A
ADC H
ADC D`

var SUB_INPUT = `SUB A,B
SUB A,C
SUB A,D
SUB A,E
SUB A,H
SUB A,L
SUB A,[HL]
SUB A,A
SUB A, 0xA6
SUB C
SUB [HL]`

var SBC_INPUT = `SBC A,B
SBC A,C
SBC A,D
SBC A,E
SBC A,H
SBC A,L
SBC A,[HL]
SBC A,A
SBC A, 0x55
SBC L
SBC [HL]`

var AND_INPUT = `AND A,B
AND A,C
AND A,D
AND A,E
AND A,H
AND A,L
AND A,[HL]
AND A,A
AND A, 0x12
AND C
AND D`

var XOR_INPUT = `XOR A,B
XOR A,C
XOR A,D
XOR A,E
XOR A,H
XOR A,L
XOR A,[HL]
XOR A,A
XOR A, 0xFE
XOR E
XOR H`

var OR_INPUT = `OR A,B
OR A,C
OR A,D
OR A,E
OR A,H
OR A,L
OR A,[HL]
OR A,A
OR A, 0x6F
OR H
OR C`

var CP_INPUT = `CP A,B
CP A,C
CP A,D
CP A,E
CP A,H
CP A,L
CP A,[HL]
CP A,A
CP A, 0x65
CP H
CP L`

var INC_INPUT = `INC B
INC D
INC H
INC [HL]
INC C
INC E
INC L
INC A`

var DEC_INPUT = `DEC B
DEC D
DEC H
DEC [HL]
DEC C
DEC E
DEC L
DEC A`

var OTHER_INPUT = `DAA
SCF
CPL
CCF`

var ADD_OUTPUT = []uint8{
	// ADD A,B
	0x80,
	// ADD A,C
	0x81,
	// ADD A,D
	0x82,
	// ADD A,E
	0x83,
	// ADD A,H
	0x84,
	// ADD A,L
	0x85,
	// ADD A,[HL]
	0x86,
	// ADD A,A
	0x87,
	// ADD H
	0x84,
	// ADD A,E
	0x83,
	// ADD HL,BC
	0x09,
	// ADD HL,DE
	0x19,
	// ADD HL,HL
	0x29,
	// ADD HL,SP
	0x39,
	// ADD A,0x7E
	0xC6,
	0x7E,
	// ADD SP,0xDF
	0xE8,
	0xDF,
}

var ADC_OUTPUT = []uint8{
	// ADC A,B
	0x88,
	// ADC A,C
	0x89,
	// ADC A,D
	0x8A,
	// ADC A,E
	0x8B,
	// ADC A,H
	0x8C,
	// ADC A,L
	0x8D,
	// ADC A,[HL]
	0x8E,
	// ADC A,A
	0x8F,
	// ADC A, 0x5A
	0xCE,
	0x5A,
	// ADC H
	0x8C,
	// ADC D
	0x8A,
}

var SUB_OUTPUT = []uint8{
	// SUB A,B
	0x90,
	// SUB A,C
	0x91,
	// SUB A,D
	0x92,
	// SUB A,E
	0x93,
	// SUB A,H
	0x94,
	// SUB A,L
	0x95,
	// SUB A,[HL]
	0x96,
	// SUB A,A
	0x97,
	// SUB A, 0xA6
	0xD6,
	0xA6,
	// SUB C
	0x91,
	// SUB [HL]
	0x96,
}

var SBC_OUTPUT = []uint8{
	// SBC A,B
	0x98,
	// SBC A,C
	0x99,
	// SBC A,D
	0x9A,
	// SBC A,E
	0x9B,
	// SBC A,H
	0x9C,
	// SBC A,L
	0x9D,
	// SBC A,[HL]
	0x9E,
	// SBC A,A
	0x9F,
	// SBC A, 0x55
	0xDE,
	0x55,
	// SBC L
	0x9D,
	// SBC [HL]
	0x9E,
}

var AND_OUTPUT = []uint8{
	// AND A,B
	0xA0,
	// AND A,C
	0xA1,
	// AND A,D
	0xA2,
	// AND A,E
	0xA3,
	// AND A,H
	0xA4,
	// AND A,L
	0xA5,
	// AND A,[HL]
	0xA6,
	// AND A,A
	0xA7,
	// AND A, 0x12
	0xE6,
	0x12,
	// AND C
	0xA1,
	// AND D
	0xA2,
}

var XOR_OUTPUT = []uint8{
	// XOR A,B
	0xA8,
	// XOR A,C
	0xA9,
	// XOR A,D
	0xAA,
	// XOR A,E
	0xAB,
	// XOR A,H
	0xAC,
	// XOR A,L
	0xAD,
	// XOR A,[HL]
	0xAE,
	// XOR A,A
	0xAF,
	// XOR A, 0xFE
	0xEE,
	0xFE,
	// XOR E
	0xAB,
	// XOR H
	0xAC,
}

var OR_OUTPUT = []uint8{
	// OR A,B
	0xB0,
	// OR A,C
	0xB1,
	// OR A,D
	0xB2,
	// OR A,E
	0xB3,
	// OR A,H
	0xB4,
	// OR A,L
	0xB5,
	// OR A,[HL]
	0xB6,
	// OR A,A
	0xB7,
	// OR A, 0x6F
	0xF6,
	0x6F,
	// OR H
	0xB4,
	// OR C
	0xB1,
}

var CP_OUTPUT = []uint8{
	// CP A,B
	0xB8,
	// CP A,C
	0xB9,
	// CP A,D
	0xBA,
	// CP A,E
	0xBB,
	// CP A,H
	0xBC,
	// CP A,L
	0xBD,
	// CP A,[HL]
	0xBE,
	// CP A,A
	0xBF,
	// CP A, 0x65
	0xFE,
	0x65,
	// CP H
	0xBC,
	// CP L
	0xBD,
}

var INC_OUTPUT = []uint8{
	// INC B
	0x04,
	// INC D
	0x14,
	// INC H
	0x24,
	// INC [HL]
	0x34,
	// INC C
	0x0C,
	// INC E
	0x1C,
	// INC L
	0x2C,
	// // INC A
	0x3C,
}

var DEC_OUTPUT = []uint8{
	// DEC B
	0x05,
	// DEC D
	0x15,
	// DEC H
	0x25,
	// DEC [HL]
	0x35,
	// DEC C
	0x0D,
	// DEC E
	0x1D,
	// DEC L
	0x2D,
	// DEC A
	0x3D,
}

var OTHER_OUTPUT = []uint8{
	// DAA
	0x27,
	// SCF
	0x37,
	// CPL
	0x2F,
	// CCF
	0x3F,
}
