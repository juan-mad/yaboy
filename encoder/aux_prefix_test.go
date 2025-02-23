package encoder_test

var PREFIX_INPUT = `CB RLC B
CB RLC C
CB RLC D
CB RLC E
CB RLC H
CB RLC L
CB RLC [HL]
CB RLC A
CB RRC B
CB RRC C
CB RRC D
CB RRC E
CB RRC H
CB RRC L
CB RRC [HL]
CB RRC A
CB RL B
CB RL C
CB RL D
CB RL E
CB RL H
CB RL L
CB RL [HL]
CB RL A
CB RR B
CB RR C
CB RR D
CB RR E
CB RR H
CB RR L
CB RR [HL]
CB RR A
CB SLA B
CB SLA C
CB SLA D
CB SLA E
CB SLA H
CB SLA L
CB SLA [HL]
CB SLA A
CB SRA B
CB SRA C
CB SRA D
CB SRA E
CB SRA H
CB SRA L
CB SRA [HL]
CB SRA A
CB SWAP B
CB SWAP C
CB SWAP D
CB SWAP E
CB SWAP H
CB SWAP L
CB SWAP [HL]
CB SWAP A
CB SRL B
CB SRL C
CB SRL D
CB SRL E
CB SRL H
CB SRL L
CB SRL [HL]
CB SRL A
CB BIT 0, B
CB BIT 0, C
CB BIT 0, D
CB BIT 0, E
CB BIT 0, H
CB BIT 0, L
CB BIT 0, [HL]
CB BIT 0, A
CB BIT 1, B
CB BIT 1, C
CB BIT 1, D
CB BIT 1, E
CB BIT 1, H
CB BIT 1, L
CB BIT 1, [HL]
CB BIT 1, A
CB BIT 2, B
CB BIT 2, C
CB BIT 2, D
CB BIT 2, E
CB BIT 2, H
CB BIT 2, L
CB BIT 2, [HL]
CB BIT 2, A
CB BIT 3, B
CB BIT 3, C
CB BIT 3, D
CB BIT 3, E
CB BIT 3, H
CB BIT 3, L
CB BIT 3, [HL]
CB BIT 3, A
CB BIT 4, B
CB BIT 4, C
CB BIT 4, D
CB BIT 4, E
CB BIT 4, H
CB BIT 4, L
CB BIT 4, [HL]
CB BIT 4, A
CB BIT 5, B
CB BIT 5, C
CB BIT 5, D
CB BIT 5, E
CB BIT 5, H
CB BIT 5, L
CB BIT 5, [HL]
CB BIT 5, A
CB BIT 6, B
CB BIT 6, C
CB BIT 6, D
CB BIT 6, E
CB BIT 6, H
CB BIT 6, L
CB BIT 6, [HL]
CB BIT 6, A
CB BIT 7, B
CB BIT 7, C
CB BIT 7, D
CB BIT 7, E
CB BIT 7, H
CB BIT 7, L
CB BIT 7, [HL]
CB BIT 7, A
CB RES 0, B
CB RES 0, C
CB RES 0, D
CB RES 0, E
CB RES 0, H
CB RES 0, L
CB RES 0, [HL]
CB RES 0, A
CB RES 1, B
CB RES 1, C
CB RES 1, D
CB RES 1, E
CB RES 1, H
CB RES 1, L
CB RES 1, [HL]
CB RES 1, A
CB RES 2, B
CB RES 2, C
CB RES 2, D
CB RES 2, E
CB RES 2, H
CB RES 2, L
CB RES 2, [HL]
CB RES 2, A
CB RES 3, B
CB RES 3, C
CB RES 3, D
CB RES 3, E
CB RES 3, H
CB RES 3, L
CB RES 3, [HL]
CB RES 3, A
CB RES 4, B
CB RES 4, C
CB RES 4, D
CB RES 4, E
CB RES 4, H
CB RES 4, L
CB RES 4, [HL]
CB RES 4, A
CB RES 5, B
CB RES 5, C
CB RES 5, D
CB RES 5, E
CB RES 5, H
CB RES 5, L
CB RES 5, [HL]
CB RES 5, A
CB RES 6, B
CB RES 6, C
CB RES 6, D
CB RES 6, E
CB RES 6, H
CB RES 6, L
CB RES 6, [HL]
CB RES 6, A
CB RES 7, B
CB RES 7, C
CB RES 7, D
CB RES 7, E
CB RES 7, H
CB RES 7, L
CB RES 7, [HL]
CB RES 7, A
CB SET 0, B
CB SET 0, C
CB SET 0, D
CB SET 0, E
CB SET 0, H
CB SET 0, L
CB SET 0, [HL]
CB SET 0, A
CB SET 1, B
CB SET 1, C
CB SET 1, D
CB SET 1, E
CB SET 1, H
CB SET 1, L
CB SET 1, [HL]
CB SET 1, A
CB SET 2, B
CB SET 2, C
CB SET 2, D
CB SET 2, E
CB SET 2, H
CB SET 2, L
CB SET 2, [HL]
CB SET 2, A
CB SET 3, B
CB SET 3, C
CB SET 3, D
CB SET 3, E
CB SET 3, H
CB SET 3, L
CB SET 3, [HL]
CB SET 3, A
CB SET 4, B
CB SET 4, C
CB SET 4, D
CB SET 4, E
CB SET 4, H
CB SET 4, L
CB SET 4, [HL]
CB SET 4, A
CB SET 5, B
CB SET 5, C
CB SET 5, D
CB SET 5, E
CB SET 5, H
CB SET 5, L
CB SET 5, [HL]
CB SET 5, A
CB SET 6, B
CB SET 6, C
CB SET 6, D
CB SET 6, E
CB SET 6, H
CB SET 6, L
CB SET 6, [HL]
CB SET 6, A
CB SET 7, B
CB SET 7, C
CB SET 7, D
CB SET 7, E
CB SET 7, H
CB SET 7, L
CB SET 7, [HL]
CB SET 7, A
`

var PREFIX_OUTPUT = []uint8{
	// RLC B
	0xCB,
	0x00,
	// RLC C
	0xCB,
	0x01,
	// RLC D
	0xCB,
	0x02,
	// RLC E
	0xCB,
	0x03,
	// RLC H
	0xCB,
	0x04,
	// RLC L
	0xCB,
	0x05,
	// RLC [HL]
	0xCB,
	0x06,
	// RLC A
	0xCB,
	0x07,
	// CB RRC B
	0xCB,
	0x08,
	// CB RRC C
	0xCB,
	0x09,
	// CB RRC D
	0xCB,
	0x0A,
	// CB RRC E
	0xCB,
	0x0B,
	// CB RRC H
	0xCB,
	0x0C,
	// CB RRC L
	0xCB,
	0x0D,
	// CB RRC [HL]
	0xCB,
	0x0E,
	// CB RRC A
	0xCB,
	0x0F,
	// CB RL B
	0xCB,
	0x10,
	// CB RL C
	0xCB,
	0x11,
	// CB RL D
	0xCB,
	0x12,
	// CB RL E
	0xCB,
	0x13,
	// CB RL H
	0xCB,
	0x14,
	// CB RL L
	0xCB,
	0x15,
	// CB RL [HL]
	0xCB,
	0x16,
	// CB RL A
	0xCB,
	0x17,
	// CB RR B
	0xCB,
	0x18,
	// CB RR C
	0xCB,
	0x19,
	// CB RR D
	0xCB,
	0x1A,
	// CB RR E
	0xCB,
	0x1B,
	// CB RR H
	0xCB,
	0x1C,
	// CB RR L
	0xCB,
	0x1D,
	// CB RR [HL]
	0xCB,
	0x1E,
	// CB RR A
	0xCB,
	0x1F,
	// CB SLA B
	0xCB,
	0x20,
	// CB SLA C
	0xCB,
	0x21,
	// CB SLA D
	0xCB,
	0x22,
	// CB SLA E
	0xCB,
	0x23,
	// CB SLA H
	0xCB,
	0x24,
	// CB SLA L
	0xCB,
	0x25,
	// CB SLA [HL]
	0xCB,
	0x26,
	// CB SLA A
	0xCB,
	0x27,
	// CB SRA B
	0xCB,
	0x28,
	// CB SRA C
	0xCB,
	0x29,
	// CB SRA D
	0xCB,
	0x2A,
	// CB SRA E
	0xCB,
	0x2B,
	// CB SRA H
	0xCB,
	0x2C,
	// CB SRA L
	0xCB,
	0x2D,
	// CB SRA [HL]
	0xCB,
	0x2E,
	// CB SRA A
	0xCB,
	0x2F,
	// CB SWAP B
	0xCB,
	0x30,
	// CB SWAP C
	0xCB,
	0x31,
	// CB SWAP D
	0xCB,
	0x32,
	// CB SWAP E
	0xCB,
	0x33,
	// CB SWAP H
	0xCB,
	0x34,
	// CB SWAP L
	0xCB,
	0x35,
	// CB SWAP [HL]
	0xCB,
	0x36,
	// CB SWAP A
	0xCB,
	0x37,
	// CB SRL B
	0xCB,
	0x38,
	// CB SRL C
	0xCB,
	0x39,
	// CB SRL D
	0xCB,
	0x3A,
	// CB SRL E
	0xCB,
	0x3B,
	// CB SRL H
	0xCB,
	0x3C,
	// CB SRL L
	0xCB,
	0x3D,
	// CB SRL [HL]
	0xCB,
	0x3E,
	// CB SRL A`
	0xCB,
	0x3F,
	// CB BIT 0, B
	0xCB,
	0x40,
	// CB BIT 0, C
	0xCB,
	0x41,
	// CB BIT 0, D
	0xCB,
	0x42,
	// CB BIT 0, E
	0xCB,
	0x43,
	// CB BIT 0, H
	0xCB,
	0x44,
	// CB BIT 0, L
	0xCB,
	0x45,
	// CB BIT 0, [HL]
	0xCB,
	0x46,
	// CB BIT 0, A
	0xCB,
	0x47,
	// CB BIT 1, B
	0xCB,
	0x48,
	// CB BIT 1, C
	0xCB,
	0x49,
	// CB BIT 1, D
	0xCB,
	0x4A,
	// CB BIT 1, E
	0xCB,
	0x4B,
	// CB BIT 1, H
	0xCB,
	0x4C,
	// CB BIT 1, L
	0xCB,
	0x4D,
	// CB BIT 1, [HL]
	0xCB,
	0x4E,
	// CB BIT 1, A
	0xCB,
	0x4F,
	// CB BIT 2, B
	0xCB,
	0x50,
	// CB BIT 2, C
	0xCB,
	0x51,
	// CB BIT 2, D
	0xCB,
	0x52,
	// CB BIT 2, E
	0xCB,
	0x53,
	// CB BIT 2, H
	0xCB,
	0x54,
	// CB BIT 2, L
	0xCB,
	0x55,
	// CB BIT 2, [HL]
	0xCB,
	0x56,
	// CB BIT 2, A
	0xCB,
	0x57,
	// CB BIT 3, B
	0xCB,
	0x58,
	// CB BIT 3, C
	0xCB,
	0x59,
	// CB BIT 3, D
	0xCB,
	0x5A,
	// CB BIT 3, E
	0xCB,
	0x5B,
	// CB BIT 3, H
	0xCB,
	0x5C,
	// CB BIT 3, L
	0xCB,
	0x5D,
	// CB BIT 3, [HL]
	0xCB,
	0x5E,
	// CB BIT 3, A
	0xCB,
	0x5F,
	// CB BIT 4, B
	0xCB,
	0x60,
	// CB BIT 4, C
	0xCB,
	0x61,
	// CB BIT 4, D
	0xCB,
	0x62,
	// CB BIT 4, E
	0xCB,
	0x63,
	// CB BIT 4, H
	0xCB,
	0x64,
	// CB BIT 4, L
	0xCB,
	0x65,
	// CB BIT 4, [HL]
	0xCB,
	0x66,
	// CB BIT 4, A
	0xCB,
	0x67,
	// CB BIT 5, B
	0xCB,
	0x68,
	// CB BIT 5, C
	0xCB,
	0x69,
	// CB BIT 5, D
	0xCB,
	0x6A,
	// CB BIT 5, E
	0xCB,
	0x6B,
	// CB BIT 5, H
	0xCB,
	0x6C,
	// CB BIT 5, L
	0xCB,
	0x6D,
	// CB BIT 5, [HL]
	0xCB,
	0x6E,
	// CB BIT 5, A
	0xCB,
	0x6F,
	// CB BIT 6, B
	0xCB,
	0x70,
	// CB BIT 6, C
	0xCB,
	0x71,
	// CB BIT 6, D
	0xCB,
	0x72,
	// CB BIT 6, E
	0xCB,
	0x73,
	// CB BIT 6, H
	0xCB,
	0x74,
	// CB BIT 6, L
	0xCB,
	0x75,
	// CB BIT 6, [HL]
	0xCB,
	0x76,
	// CB BIT 6, A
	0xCB,
	0x77,
	// CB BIT 7, B
	0xCB,
	0x78,
	// CB BIT 7, C
	0xCB,
	0x79,
	// CB BIT 7, D
	0xCB,
	0x7A,
	// CB BIT 7, E
	0xCB,
	0x7B,
	// CB BIT 7, H
	0xCB,
	0x7C,
	// CB BIT 7, L
	0xCB,
	0x7D,
	// CB BIT 7, [HL]
	0xCB,
	0x7E,
	// CB BIT 7, A
	0xCB,
	0x7F,
	// CB RES 0, B
	0xCB,
	0x80,
	// CB RES 0, C
	0xCB,
	0x81,
	// CB RES 0, D
	0xCB,
	0x82,
	// CB RES 0, E
	0xCB,
	0x83,
	// CB RES 0, H
	0xCB,
	0x84,
	// CB RES 0, L
	0xCB,
	0x85,
	// CB RES 0, [HL]
	0xCB,
	0x86,
	// CB RES 0, A
	0xCB,
	0x87,
	// CB RES 1, B
	0xCB,
	0x88,
	// CB RES 1, C
	0xCB,
	0x89,
	// CB RES 1, D
	0xCB,
	0x8A,
	// CB RES 1, E
	0xCB,
	0x8B,
	// CB RES 1, H
	0xCB,
	0x8C,
	// CB RES 1, L
	0xCB,
	0x8D,
	// CB RES 1, [HL]
	0xCB,
	0x8E,
	// CB RES 1, A
	0xCB,
	0x8F,
	// CB RES 2, B
	0xCB,
	0x90,
	// CB RES 2, C
	0xCB,
	0x91,
	// CB RES 2, D
	0xCB,
	0x92,
	// CB RES 2, E
	0xCB,
	0x93,
	// CB RES 2, H
	0xCB,
	0x94,
	// CB RES 2, L
	0xCB,
	0x95,
	// CB RES 2, [HL]
	0xCB,
	0x96,
	// CB RES 2, A
	0xCB,
	0x97,
	// CB RES 3, B
	0xCB,
	0x98,
	// CB RES 3, C
	0xCB,
	0x99,
	// CB RES 3, D
	0xCB,
	0x9A,
	// CB RES 3, E
	0xCB,
	0x9B,
	// CB RES 3, H
	0xCB,
	0x9C,
	// CB RES 3, L
	0xCB,
	0x9D,
	// CB RES 3, [HL]
	0xCB,
	0x9E,
	// CB RES 3, A
	0xCB,
	0x9F,
	// CB RES 4, B
	0xCB,
	0xA0,
	// CB RES 4, C
	0xCB,
	0xA1,
	// CB RES 4, D
	0xCB,
	0xA2,
	// CB RES 4, E
	0xCB,
	0xA3,
	// CB RES 4, H
	0xCB,
	0xA4,
	// CB RES 4, L
	0xCB,
	0xA5,
	// CB RES 4, [HL]
	0xCB,
	0xA6,
	// CB RES 4, A
	0xCB,
	0xA7,
	// CB RES 5, B
	0xCB,
	0xA8,
	// CB RES 5, C
	0xCB,
	0xA9,
	// CB RES 5, D
	0xCB,
	0xAA,
	// CB RES 5, E
	0xCB,
	0xAB,
	// CB RES 5, H
	0xCB,
	0xAC,
	// CB RES 5, L
	0xCB,
	0xAD,
	// CB RES 5, [HL]
	0xCB,
	0xAE,
	// CB RES 5, A
	0xCB,
	0xAF,
	// CB RES 6, B
	0xCB,
	0xB0,
	// CB RES 6, C
	0xCB,
	0xB1,
	// CB RES 6, D
	0xCB,
	0xB2,
	// CB RES 6, E
	0xCB,
	0xB3,
	// CB RES 6, H
	0xCB,
	0xB4,
	// CB RES 6, L
	0xCB,
	0xB5,
	// CB RES 6, [HL]
	0xCB,
	0xB6,
	// CB RES 6, A
	0xCB,
	0xB7,
	// CB RES 7, B
	0xCB,
	0xB8,
	// CB RES 7, C
	0xCB,
	0xB9,
	// CB RES 7, D
	0xCB,
	0xBA,
	// CB RES 7, E
	0xCB,
	0xBB,
	// CB RES 7, H
	0xCB,
	0xBC,
	// CB RES 7, L
	0xCB,
	0xBD,
	// CB RES 7, [HL]
	0xCB,
	0xBE,
	// CB RES 7, A
	0xCB,
	0xBF,
	// CB SET 0, B
	0xCB,
	0xC0,
	// CB SET 0, C
	0xCB,
	0xC1,
	// CB SET 0, D
	0xCB,
	0xC2,
	// CB SET 0, E
	0xCB,
	0xC3,
	// CB SET 0, H
	0xCB,
	0xC4,
	// CB SET 0, L
	0xCB,
	0xC5,
	// CB SET 0, [HL]
	0xCB,
	0xC6,
	// CB SET 0, A
	0xCB,
	0xC7,
	// CB SET 1, B
	0xCB,
	0xC8,
	// CB SET 1, C
	0xCB,
	0xC9,
	// CB SET 1, D
	0xCB,
	0xCA,
	// CB SET 1, E
	0xCB,
	0xCB,
	// CB SET 1, H
	0xCB,
	0xCC,
	// CB SET 1, L
	0xCB,
	0xCD,
	// CB SET 1, [HL]
	0xCB,
	0xCE,
	// CB SET 1, A
	0xCB,
	0xCF,
	// CB SET 2, B
	0xCB,
	0xD0,
	// CB SET 2, C
	0xCB,
	0xD1,
	// CB SET 2, D
	0xCB,
	0xD2,
	// CB SET 2, E
	0xCB,
	0xD3,
	// CB SET 2, H
	0xCB,
	0xD4,
	// CB SET 2, L
	0xCB,
	0xD5,
	// CB SET 2, [HL]
	0xCB,
	0xD6,
	// CB SET 2, A
	0xCB,
	0xD7,
	// CB SET 3, B
	0xCB,
	0xD8,
	// CB SET 3, C
	0xCB,
	0xD9,
	// CB SET 3, D
	0xCB,
	0xDA,
	// CB SET 3, E
	0xCB,
	0xDB,
	// CB SET 3, H
	0xCB,
	0xDC,
	// CB SET 3, L
	0xCB,
	0xDD,
	// CB SET 3, [HL]
	0xCB,
	0xDE,
	// CB SET 3, A
	0xCB,
	0xDF,
	// CB SET 4, B
	0xCB,
	0xE0,
	// CB SET 4, C
	0xCB,
	0xE1,
	// CB SET 4, D
	0xCB,
	0xE2,
	// CB SET 4, E
	0xCB,
	0xE3,
	// CB SET 4, H
	0xCB,
	0xE4,
	// CB SET 4, L
	0xCB,
	0xE5,
	// CB SET 4, [HL]
	0xCB,
	0xE6,
	// CB SET 4, A
	0xCB,
	0xE7,
	// CB SET 5, B
	0xCB,
	0xE8,
	// CB SET 5, C
	0xCB,
	0xE9,
	// CB SET 5, D
	0xCB,
	0xEA,
	// CB SET 5, E
	0xCB,
	0xEB,
	// CB SET 5, H
	0xCB,
	0xEC,
	// CB SET 5, L
	0xCB,
	0xED,
	// CB SET 5, [HL]
	0xCB,
	0xEE,
	// CB SET 5, A
	0xCB,
	0xEF,
	// CB SET 6, B
	0xCB,
	0xF0,
	// CB SET 6, C
	0xCB,
	0xF1,
	// CB SET 6, D
	0xCB,
	0xF2,
	// CB SET 6, E
	0xCB,
	0xF3,
	// CB SET 6, H
	0xCB,
	0xF4,
	// CB SET 6, L
	0xCB,
	0xF5,
	// CB SET 6, [HL]
	0xCB,
	0xF6,
	// CB SET 6, A
	0xCB,
	0xF7,
	// CB SET 7, B
	0xCB,
	0xF8,
	// CB SET 7, C
	0xCB,
	0xF9,
	// CB SET 7, D
	0xCB,
	0xFA,
	// CB SET 7, E
	0xCB,
	0xFB,
	// CB SET 7, H
	0xCB,
	0xFC,
	// CB SET 7, L
	0xCB,
	0xFD,
	// CB SET 7, [HL]
	0xCB,
	0xFE,
	// CB SET 7, A
	0xCB,
	0xFF,
}
