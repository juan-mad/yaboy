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
