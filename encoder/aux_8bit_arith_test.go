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
