package decompiler_test

var INPUT_0x = []uint8{
	// NOP
	0x00,
	// LD BC,0x1234
	0x01,
	0x34,
	0x12,
	// LD [BC],A
	0x02,
	// INC BC
	0x03,
	// INC B
	0x04,
	// DEC B
	0x05,
	// LD B,0xff
	0x06,
	0xFF,
	// RLCA
	0x07,
	// LD [0xfaef],SP
	0x08,
	0xEF,
	0xFA,
	// ADD HL,BC
	0x09,
	// LD A,[BC]
	0x0A,
	// DEC BC
	0x0B,
	// INC C
	0x0C,
	// DEC C
	0x0D,
	// LD C,0x67
	0x0E,
	0x67,
	// RRCA
	0x0F,
}

var OUTPUT_0x = []string{
	"NOP",
	"LD BC,0x1234",
	"LD [BC],A",
	"INC BC",
	"INC B",
	"DEC B",
	"LD B,0xff",
	"RLCA",
	"LD [0xfaef],SP",
	"ADD HL,BC",
	"LD A,[BC]",
	"DEC BC",
	"INC C",
	"DEC C",
	"LD C,0x67",
	"RRCA",
}

var INPUT_1x = []uint8{
	// STOP 0x5E
	0x10,
	0x5E,
	// LD DE,0x56BC
	0x11,
	0xBC,
	0x56,
	// LD [DE],A
	0x12,
	// INC DE
	0x13,
	// INC D
	0x14,
	// DEC D
	0x15,
	// LD D,0xA4
	0x16,
	0xA4,
	// RLA
	0x17,
	// JR 0x14
	0x18,
	0x14,
	// ADD HL,DE
	0x19,
	// LD A,[DE]
	0x1A,
	// DEC DE
	0x1B,
	// INC E
	0x1C,
	// DEC E
	0x1D,
	// LD E,0x12
	0x1E,
	0x12,
	// RRA
	0x1F,
}

var OUTPUT_1x = []string{
	"STOP 0x5e",
	"LD DE,0x56bc",
	"LD [DE],A",
	"INC DE",
	"INC D",
	"DEC D",
	"LD D,0xa4",
	"RLA",
	"JR 0x14",
	"ADD HL,DE",
	"LD A,[DE]",
	"DEC DE",
	"INC E",
	"DEC E",
	"LD E,0x12",
	"RRA",
}

var INPUT_2x = []uint8{
	// JR NZ,0x12
	0x20,
	0x12,
	// LD HL,0x1234
	0x21,
	0x34,
	0x12,
	// LD [HL+],A
	0x22,
	// INC HL
	0x23,
	// INC H
	0x24,
	// DEC H
	0x25,
	// LD H,0x12
	0x26,
	0x12,
	// DAA
	0x27,
	// JR Z,0xaa
	0x28,
	0xAA,
	// ADD HL,HL
	0x29,
	// LD A,[HL+]
	0x2A,
	// DEC HL
	0x2B,
	// INC L
	0x2C,
	// DEC L
	0x2D,
	// LD L,0x12
	0x2E,
	0x12,
	// CPL
	0x2F,
}

var OUTPUT_2x = []string{
	"JR NZ,0x12",
	"LD HL,0x1234",
	"LD [HL+],A",
	"INC HL",
	"INC H",
	"DEC H",
	"LD H,0x12",
	"DAA",
	"JR Z,0xaa",
	"ADD HL,HL",
	"LD A,[HL+]",
	"DEC HL",
	"INC L",
	"DEC L",
	"LD L,0x12",
	"CPL",
}

var INPUT_3x = []uint8{
	//. JR NC,0x12
	0x30,
	0x12,
	// LD SP,0x1234
	0x31,
	0x34,
	0x12,
	// LD [HL-],A
	0x32,
	// INC SP
	0x33,
	// INC [HL]
	0x34,
	// DEC [HL]
	0x35,
	// LD [HL],0x12
	0x36,
	0x12,
	// SCF
	0x37,
	// JR C,0x12
	0x38,
	0x12,
	// ADD HL,SP
	0x39,
	// LD A,[HL-]
	0x3A,
	// DEC SP
	0x3B,
	// INC A
	0x3C,
	// DEC A
	0x3D,
	// LD A,0x12
	0x3E,
	0x12,
	// CCF
	0x3F,
}

var OUTPUT_3x = []string{
	"JR NC,0x12",
	"LD SP,0x1234",
	"LD [HL-],A",
	"INC SP",
	"INC [HL]",
	"DEC [HL]",
	"LD [HL],0x12",
	"SCF",
	"JR C,0x12",
	"ADD HL,SP",
	"LD A,[HL-]",
	"DEC SP",
	"INC A",
	"DEC A",
	"LD A,0x12",
	"CCF",
}

var INPUT_4x = []uint8{
	// LD B,B
	0x40,
	// LD B,C
	0x41,
	// LD B,D
	0x42,
	// LD B,E
	0x43,
	// LD B,H
	0x44,
	// LD B,L
	0x45,
	// LD B,[HL]
	0x46,
	// LD B,A
	0x47,
	// LD C,B
	0x48,
	// LD C,C
	0x49,
	// LD C,D
	0x4A,
	// LD C,E
	0x4B,
	// LD C,H
	0x4C,
	// LD C,L
	0x4D,
	// LD C,[HL]
	0x4E,
	// LD C,A
	0x4F,
}

var OUTPUT_4x = []string{
	"LD B,B",
	"LD B,C",
	"LD B,D",
	"LD B,E",
	"LD B,H",
	"LD B,L",
	"LD B,[HL]",
	"LD B,A",
	"LD C,B",
	"LD C,C",
	"LD C,D",
	"LD C,E",
	"LD C,H",
	"LD C,L",
	"LD C,[HL]",
	"LD C,A",
}

var INPUT_5x = []uint8{
	// LD D,B
	0x50,
	// LD D,C
	0x51,
	// LD D,D
	0x52,
	// LD D,E
	0x53,
	// LD D,H
	0x54,
	// LD D,L
	0x55,
	// LD D,[HL]
	0x56,
	// LD D,A
	0x57,
	// LD E,B
	0x58,
	// LD E,C
	0x59,
	// LD E,D
	0x5A,
	// LD E,E
	0x5B,
	// LD E,H
	0x5C,
	// LD E,L
	0x5D,
	// LD E,[HL]
	0x5E,
	// LD E,A
	0x5F,
}

var OUTPUT_5x = []string{
	"LD D,B",
	"LD D,C",
	"LD D,D",
	"LD D,E",
	"LD D,H",
	"LD D,L",
	"LD D,[HL]",
	"LD D,A",
	"LD E,B",
	"LD E,C",
	"LD E,D",
	"LD E,E",
	"LD E,H",
	"LD E,L",
	"LD E,[HL]",
	"LD E,A",
}

var INPUT_6x = []uint8{
	// LD H,B
	0x60,
	// LD H,C
	0x61,
	// LD H,D
	0x62,
	// LD H,E
	0x63,
	// LD H,H
	0x64,
	// LD H,L
	0x65,
	// LD H,[HL]
	0x66,
	// LD H,A
	0x67,
	// LD L,B
	0x68,
	// LD L,C
	0x69,
	// LD L,D
	0x6A,
	// LD L,E
	0x6B,
	// LD L,H
	0x6C,
	// LD L,L
	0x6D,
	// LD L,[HL]
	0x6E,
	// LD L,A
	0x6F,
}

var OUTPUT_6x = []string{
	"LD H,B",
	"LD H,C",
	"LD H,D",
	"LD H,E",
	"LD H,H",
	"LD H,L",
	"LD H,[HL]",
	"LD H,A",
	"LD L,B",
	"LD L,C",
	"LD L,D",
	"LD L,E",
	"LD L,H",
	"LD L,L",
	"LD L,[HL]",
	"LD L,A",
}

var INPUT_7x = []uint8{
	// LD [HL],B
	0x70,
	// LD [HL],C
	0x71,
	// LD [HL],D
	0x72,
	// LD [HL],E
	0x73,
	// LD [HL],H
	0x74,
	// LD [HL],L
	0x75,
	// HALT
	0x76,
	// LD [HL],A
	0x77,
	// LD A,B
	0x78,
	// LD A,C
	0x79,
	// LD A,D
	0x7A,
	// LD A,E
	0x7B,
	// LD A,H
	0x7C,
	// LD A,L
	0x7D,
	// LD A,[HL]
	0x7E,
	// LD A,A
	0x7F,
}

var OUTPUT_7x = []string{
	"LD [HL],B",
	"LD [HL],C",
	"LD [HL],D",
	"LD [HL],E",
	"LD [HL],H",
	"LD [HL],L",
	"HALT",
	"LD [HL],A",
	"LD A,B",
	"LD A,C",
	"LD A,D",
	"LD A,E",
	"LD A,H",
	"LD A,L",
	"LD A,[HL]",
	"LD A,A",
}

var INPUT_8x = []uint8{
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
}

var OUTPUT_8x = []string{
	"ADD A,B",
	"ADD A,C",
	"ADD A,D",
	"ADD A,E",
	"ADD A,H",
	"ADD A,L",
	"ADD A,[HL]",
	"ADD A,A",
	"ADC A,B",
	"ADC A,C",
	"ADC A,D",
	"ADC A,E",
	"ADC A,H",
	"ADC A,L",
	"ADC A,[HL]",
	"ADC A,A",
}

var INPUT_9x = []uint8{
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
}

var OUTPUT_9x = []string{
	"SUB A,B",
	"SUB A,C",
	"SUB A,D",
	"SUB A,E",
	"SUB A,H",
	"SUB A,L",
	"SUB A,[HL]",
	"SUB A,A",
	"SBC A,B",
	"SBC A,C",
	"SBC A,D",
	"SBC A,E",
	"SBC A,H",
	"SBC A,L",
	"SBC A,[HL]",
	"SBC A,A",
}

var INPUT_Ax = []uint8{
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
}

var OUTPUT_Ax = []string{
	"AND A,B",
	"AND A,C",
	"AND A,D",
	"AND A,E",
	"AND A,H",
	"AND A,L",
	"AND A,[HL]",
	"AND A,A",
	"XOR A,B",
	"XOR A,C",
	"XOR A,D",
	"XOR A,E",
	"XOR A,H",
	"XOR A,L",
	"XOR A,[HL]",
	"XOR A,A",
}

var INPUT_Bx = []uint8{
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
}

var OUTPUT_Bx = []string{
	"OR A,B",
	"OR A,C",
	"OR A,D",
	"OR A,E",
	"OR A,H",
	"OR A,L",
	"OR A,[HL]",
	"OR A,A",
	"CP A,B",
	"CP A,C",
	"CP A,D",
	"CP A,E",
	"CP A,H",
	"CP A,L",
	"CP A,[HL]",
	"CP A,A",
}

var INPUT_Cx = []uint8{
	// "RET NZ",
	0xC0,
	// "POP BC",
	0xC1,
	// "JP NZ,0x1234",
	0xC2,
	0x34,
	0x12,
	// "JP 0xABCD",
	0xC3,
	0xCD,
	0xAB,
	// "CALL NZ,0x1234",
	0xC4,
	0x34,
	0x12,
	// "PUSH BC",
	0xC5,
	// "ADD A,0x12",
	0xC6,
	0x12,
	// "RST 0x00",
	0xC7,
	// "RET Z",
	0xC8,
	// "RET",
	0xC9,
	// "JP Z,0x1234",
	0xCA,
	0x34,
	0x12,
	// "CALL Z,0x1234",
	0xCC,
	0x34,
	0x12,
	// "CALL 0x1234",
	0xCD,
	0x34,
	0x12,
	// "ADC A,0x12",
	0xCE,
	0x12,
	// "RST 0x08",
	0xCF,
}

var OUTPUT_Cx = []string{
	"RET NZ",
	"POP BC",
	"JP NZ,0x1234",
	"JP 0xabcd",
	"CALL NZ,0x1234",
	"PUSH BC",
	"ADD A,0x12",
	"RST 0x00",
	"RET Z",
	"RET",
	"JP Z,0x1234",
	"CALL Z,0x1234",
	"CALL 0x1234",
	"ADC A,0x12",
	"RST 0x08",
}

var INPUT_Dx = []uint8{
	// "RET NC",
	0xD0,
	// "POP DE",
	0xD1,
	// "JP NC,0x1234",
	0xD2,
	0x34,
	0x12,
	// "CALL NC,0x1234",
	0xD4,
	0x34,
	0x12,
	// "PUSH DE",
	0xD5,
	// "SUB A,0x12",
	0xD6,
	0x12,
	// "RES 0x10",
	0xD7,
	// "RET C",
	0xD8,
	// "RETI",
	0xD9,
	// "JP C,0x1234",
	0xDA,
	0x34,
	0x12,
	// "CALL C,0x1234",
	0xDC,
	0x34,
	0x12,
	// "SBC A,0x12",
	0xDE,
	0x12,
	// "RST 0x18",
	0xDF,
}

var OUTPUT_Dx = []string{
	"RET NC",
	"POP DE",
	"JP NC,0x1234",
	"CALL NC,0x1234",
	"PUSH DE",
	"SUB A,0x12",
	"RST 0x10",
	"RET C",
	"RETI",
	"JP C,0x1234",
	"CALL C,0x1234",
	"SBC A,0x12",
	"RST 0x18",
}

var INPUT_Ex = []uint8{
	// "LDH [0x12],A",
	0xE0,
	0x12,
	// "POP HL",
	0xE1,
	// "LDH [C],A",
	0xE2,
	// "PUSH HL",
	0xE5,
	// "AND A,0x12",
	0xE6,
	0x12,
	// "RST 0x20",
	0xE7,
	// "ADD SP,0x12",
	0xE8,
	0x12,
	// "JP HL",
	0xE9,
	// "LD [0x1234],A",
	0xEA,
	0x34,
	0x12,
	// "XOR A,0x12",
	0xEE,
	0x12,
	// "RST 0x28",
	0xEF,
}

var OUTPUT_Ex = []string{
	"LDH [0x12],A",
	"POP HL",
	"LDH [C],A",
	"PUSH HL",
	"AND A,0x12",
	"RST 0x20",
	"ADD SP,0x12",
	"JP HL",
	"LD [0x1234],A",
	"XOR A,0x12",
	"RST 0x28",
}

var INPUT_Fx = []uint8{
	// "LDH A,[0x12]",
	0xF0,
	0x12,
	// "POP AF",
	0xF1,
	// "LDH A,[C]",
	0xF2,
	// "DI",
	0xF3,
	// "PUSH AF",
	0xF5,
	// "OR A,0x12",
	0xF6,
	0x12,
	// "RST 0x30",
	0xF7,
	// "LD HL,SP+0x12",
	0xF8,
	0x12,
	// "LD SP,HL",
	0xF9,
	// "LD A,[0x1234]",
	0xFA,
	0x34,
	0x12,
	// "EI",
	0xFB,
	// "CP A,0x12",
	0xFE,
	0x12,
	// "RST 0x38",
	0xFF,
}

var OUTPUT_Fx = []string{
	"LDH A,[0x12]",
	"POP AF",
	"LDH A,[C]",
	"DI",
	"PUSH AF",
	"OR A,0x12",
	"RST 0x30",
	"LD HL,SP+0x12",
	"LD SP,HL",
	"LD A,[0x1234]",
	"EI",
	"CP A,0x12",
	"RST 0x38",
}
