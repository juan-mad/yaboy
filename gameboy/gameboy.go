package gameboy

type Ram [1 << 16]uint8

type GameBoy struct {
	RAM     Ram
	Reg     CPURegisters
	Started bool
}

func (gb GameBoy) ReadByAddress(addr uint16) byte {
	return gb.RAM[addr]
}

func (gb *GameBoy) WriteOnAddress(value byte, addr uint16) {
	gb.RAM[addr] = value
}

func (gb *GameBoy) FetchCycle(addr uint16) (byte, byte) {
	return gb.ReadByAddress(addr), 0
}

func (gb *GameBoy) ExecuteInstruction() {
	if !gb.Started {
		// Only fetch cycle
		// M1
		gb.Started = true

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
		return
	}

	opcode := gb.Reg.IR
	switch opcode {
	case 0x00:
		return
	case 0x01:
		return
	case 0x02:
		return
	case 0x03:
		return
	case 0x04:
		return
	case 0x05:
		return
	case 0x06: // LD B, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(Z)

	case 0x07:
		return
	case 0x08:
		return
	case 0x09:
		return
	case 0x0A:
		return
	case 0x0B:
		return
	case 0x0C:
		return
	case 0x0D:
		return
	case 0x0E: // LD C, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteC(Z)

	case 0x0F:
		return
	case 0x10:
		return
	case 0x11:
		return
	case 0x12:
		return
	case 0x13:
		return
	case 0x14:
		return
	case 0x15:
		return
	case 0x16: // LD D, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(Z)

	case 0x17:
		return
	case 0x18:
		return
	case 0x19:
		return
	case 0x1A:
		return
	case 0x1B:
		return
	case 0x1C:
		return
	case 0x1D:
		return
	case 0x1E: // LD E, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(Z)

	case 0x1F:
		return
	case 0x20:
		return
	case 0x21:
		return
	case 0x22:
		return
	case 0x23:
		return
	case 0x24:
		return
	case 0x25:
		return
	case 0x26: // LD H, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(Z)

	case 0x27:
		return
	case 0x28:
		return
	case 0x29:
		return
	case 0x2A:
		return
	case 0x2B:
		return
	case 0x2C:
		return
	case 0x2D:
		return
	case 0x2E: // LD L, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(Z)

	case 0x2F:
		return
	case 0x30:
		return
	case 0x31:
		return
	case 0x32:
		return
	case 0x33:
		return
	case 0x34:
		return
	case 0x35:
		return
	case 0x36: // LD [HL], n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3
		gb.WriteOnAddress(Z, gb.Reg.ReadHL())

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(Z)

	case 0x37:
		return
	case 0x38:
		return
	case 0x39:
		return
	case 0x3A:
		return
	case 0x3B:
		return
	case 0x3C:
		return
	case 0x3D:
		return
	case 0x3E: // LD A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.PC)
		gb.Reg.IncreasePC()

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(Z)

	case 0x3F:
		return

	case 0x40: // LD B, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(gb.Reg.ReadB())

	case 0x41: // LD B, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(gb.Reg.ReadC())
	case 0x42: // LD B, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(gb.Reg.ReadD())

	case 0x43: // LD B, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(gb.Reg.ReadE())

	case 0x44: // LD B, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(gb.Reg.ReadH())
	case 0x45: // LD B, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(gb.Reg.ReadL())
	case 0x46: // LD B, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteB(Z)

	case 0x47: //  LD B, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadA())

	case 0x48: // LD A, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadB())

	case 0x49: // LD A, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadC())

	case 0x4A: // LD A, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadD())

	case 0x4B: // LD A, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadE())

	case 0x4C: // LD A, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadH())

	case 0x4D: // LD A, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadL())

	case 0x4E: // LD A, B
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(Z)

	case 0x4F: //  LD A, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(gb.Reg.ReadA())

	case 0x50: // LD D, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadB())

	case 0x51: // LD D, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadC())

	case 0x52: // LD D, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadD())

	case 0x53: // LD D, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadE())

	case 0x54: // LD D, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadH())

	case 0x55: // LD D, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadL())

	case 0x56: // LD D, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(Z)

	case 0x57: //  LD D, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteD(gb.Reg.ReadA())

	case 0x58: // LD E, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadB())

	case 0x59: // LD E, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadC())

	case 0x5A: // LD E, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadD())

	case 0x5B: // LD E, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadE())

	case 0x5C: // LD E, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadH())

	case 0x5D: // LD E, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadL())

	case 0x5E: // LD E, B
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(Z)

	case 0x5F: //  LD E, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteE(gb.Reg.ReadA())

	case 0x60: // LD H, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadB())

	case 0x61: // LD H, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadC())

	case 0x62: // LD H, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadD())

	case 0x63: // LD H, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadE())

	case 0x64: // LD H, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadH())

	case 0x65: // LD H, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadL())

	case 0x66: // LD H, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(Z)

	case 0x67: //  LD H, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(gb.Reg.ReadA())

	case 0x68: // LD L, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadB())

	case 0x69: // LD L, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadC())

	case 0x6A: // LD L, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadD())

	case 0x6B: // LD L, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadE())

	case 0x6C: // LD L, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadH())

	case 0x6D: // LD L, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadL())

	case 0x6E: // LD L, B
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(Z)

	case 0x6F: //  LD L, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadA())

	case 0x70: // LD [HL], B
		// M2
		gb.WriteOnAddress(gb.Reg.ReadB(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x71: // LD [HL], C
		// M2
		gb.WriteOnAddress(gb.Reg.ReadC(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x72: // LD [HL], D
		// M2
		gb.WriteOnAddress(gb.Reg.ReadD(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x73: // LD [HL], E
		// M2
		gb.WriteOnAddress(gb.Reg.ReadE(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x74: // LD [HL], H
		// M2
		gb.WriteOnAddress(gb.Reg.ReadH(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x75: // LD [HL], L
		// M2
		gb.WriteOnAddress(gb.Reg.ReadL(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x76: // LD [HL], [HL]
		// TO BE IMPLEMENTED
		break

	case 0x77: //  LD [HL], A
		// M2
		gb.WriteOnAddress(gb.Reg.ReadA(), gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x78: // LD L, B
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadB())

	case 0x79: // LD L, C
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadC())

	case 0x7A: // LD L, D
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadD())

	case 0x7B: // LD L, E
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadE())

	case 0x7C: // LD L, H
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadH())

	case 0x7D: // LD L, L
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadL())

	case 0x7E: // LD L, B
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(Z)

	case 0x7F: //  LD L, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadA())
	}

}

func NewGameBoy() GameBoy {
	return GameBoy{}
}
