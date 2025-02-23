package encoder_test

var INC16_INPUT = `INC BC
INC DE
INC HL
INC SP`

var DEC16_INPUT = `DEC BC
DEC DE
DEC HL
DEC SP`

var ADD16_INPUT = `ADD HL,BC
ADD HL,DE
ADD HL,HL
ADD HL,SP
ADD SP,0x66`

var INC16_OUTPUT = []uint8{
	// INC BC
	0x03,
	// INC DE
	0x13,
	// INC HL
	0x23,
	// INC SP
	0x33,
}

var DEC16_OUTPUT = []uint8{
	// DEC BC
	0x0B,
	// DEC DE
	0x1B,
	// DEC HL
	0x2B,
	// DEC SP
	0x3B,
}

var ADD16_OUTPUT = []uint8{
	// ADD HL,BC
	0x09,
	// ADD HL,DE
	0x19,
	// ADD HL,HL
	0x29,
	// ADD HL,SP
	0x39,
	// ADD SP,0x66
	0xE8,
	0x66,
}
