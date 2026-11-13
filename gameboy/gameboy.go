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

// Adds up two unsigned 8 bit variables. Returns the value of the
// sum and the carry array
func AddUnsigned8(a, b uint8) (result uint8, carry_per_bit uint8) {
	return a + b, a & b
}

func AddUnsigned8Carry(a, b uint8, carry bool) (result uint8, halfCarry, carryOut bool) {
	sum := uint16(a) + uint16(b)
	lowNibbleSum := (a & 0x0F) + (b & 0x0F)
	if carry {
		sum = sum + 1
		lowNibbleSum = lowNibbleSum + 1
	}

	return uint8(sum), lowNibbleSum > 0x0F, sum > 0xFF
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
	case 0x01: // LD BC, n16
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		W := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteBC((uint16(W) << 8) | uint16(Z))

	case 0x02: // LD [BC], A
		// M2
		gb.WriteOnAddress(gb.Reg.ReadA(), gb.Reg.ReadBC())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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
	case 0x08: // LD (nn), SP
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		W := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		target := (uint16(W) << 8) | uint16(Z)

		// M4
		gb.WriteOnAddress(byte(gb.Reg.ReadSP()&0xF), target)
		target = target + 1

		// M5
		gb.WriteOnAddress(byte(gb.Reg.ReadSP()>>8), target)

		// M6/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x09:
		return
	case 0x0A: // LD A, [BC]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadBC())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(Z)

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
	case 0x11: // LD DE, n16
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		W := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteDE((uint16(W) << 8) | uint16(Z))

	case 0x12: // LD [DE], A
		// M2
		gb.WriteOnAddress(gb.Reg.ReadA(), gb.Reg.ReadDE())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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
	case 0x1A: // LD A, [DE]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadDE())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(Z)

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
	case 0x21: // LD HL, n16
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		W := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteHL((uint16(W) << 8) | uint16(Z))

	case 0x22: // LD [HL+], A
		// M2
		gb.WriteOnAddress(gb.Reg.ReadA(), gb.Reg.ReadHL())
		gb.Reg.WriteHL(gb.Reg.ReadHL() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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
	case 0x2A: // LD A, [HL+]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())
		gb.Reg.WriteHL(gb.Reg.ReadHL() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
		gb.Reg.WriteA(Z)

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
	case 0x31: // LD SP, n16
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		W := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteSP((uint16(W) << 8) | uint16(Z))

	case 0x32: // LD [HL-], A
		// M2
		gb.WriteOnAddress(gb.Reg.ReadA(), gb.Reg.ReadHL())
		gb.Reg.WriteHL(gb.Reg.ReadHL() - 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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
	case 0x3A: // LD A, [HL-]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())
		gb.Reg.WriteHL(gb.Reg.ReadHL() - 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
		gb.Reg.WriteA(Z)

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

	case 0x76: // HALT
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

	case 0x7F: // LD L, A
		// M2/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteL(gb.Reg.ReadA())

	case 0x80: // ADD A, B
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadB())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x81: // ADD A, C
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadC())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x82: // ADD A, D
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadD())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x83: // ADD A, E
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadE())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x84: // ADD A, H
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadH())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x85: // ADD A, L
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadL())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x86: // ADD A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x87: // ADD A, A
		// M2/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadA())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x88: // ADC A, B
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadB(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x89: // ADC A, C
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadC(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x8A: // ADC A, D
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadD(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x8B: // ADC A, E
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadE(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x8C: // ADC A, H
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadH(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x8D: // ADC A, L
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadL(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x8E: // ADC A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), Z, gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x8F: // ADC A, A
		// M2/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x90:
	case 0x91:
	case 0x92:
	case 0x93:
	case 0x94:
	case 0x95:
	case 0x96:
	case 0x97:
	case 0x98:
	case 0x99:
	case 0x9A:
	case 0x9B:
	case 0x9C:
	case 0x9D:
	case 0x9E:
	case 0x9F:
	case 0xA0:
	case 0xA1:
	case 0xA2:
	case 0xA3:
	case 0xA4:
	case 0xA5:
	case 0xA6:
	case 0xA7:
	case 0xA8:
	case 0xA9:
	case 0xAA:
	case 0xAB:
	case 0xAC:
	case 0xAD:
	case 0xAE:
	case 0xAF:
	case 0xB0:
	case 0xB1:
	case 0xB2:
	case 0xB3:
	case 0xB4:
	case 0xB5:
	case 0xB6:
	case 0xB7:
	case 0xB8:
	case 0xB9:
	case 0xBA:
	case 0xBB:
	case 0xBC:
	case 0xBD:
	case 0xBE:
	case 0xBF:
	case 0xC0:
	case 0xC1: // POP BC
		// this instruction does not do calculations that affect flags
		Z := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		W := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteBC((uint16(W) << 8) | uint16(Z))

	case 0xC2:
	case 0xC3:
	case 0xC4:
	case 0xC5: // PUSH BC
		// M2
		gb.Reg.DecreaseSP()

		// M3
		gb.WriteOnAddress(byte(gb.Reg.ReadBC()&0xF), gb.Reg.ReadSP())
		gb.Reg.DecreaseSP()

		// M4
		gb.WriteOnAddress(byte(gb.Reg.ReadBC()>>8), gb.Reg.ReadSP())

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xC6: // ADD A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result, carry := AddUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 1)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xC7:
	case 0xC8:
	case 0xC9:
	case 0xCA:
	case 0xCB:
	case 0xCC:
	case 0xCD:
	case 0xCE: // ADC A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result, halfCarry, carryOut := AddUnsigned8Carry(gb.Reg.ReadA(), Z, gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carryOut)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xCF:
	case 0xD0:
	case 0xD1: // POP DE
		// this instruction does not do calculations that affect flags
		Z := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		W := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteDE((uint16(W) << 8) | uint16(Z))

	case 0xD2:
	case 0xD3:
	case 0xD4:
	case 0xD5: // PUSH DE
		// M2
		gb.Reg.DecreaseSP()

		// M3
		gb.WriteOnAddress(byte(gb.Reg.ReadDE()&0xF), gb.Reg.ReadSP())
		gb.Reg.DecreaseSP()

		// M4
		gb.WriteOnAddress(byte(gb.Reg.ReadDE()>>8), gb.Reg.ReadSP())

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xD6:
	case 0xD7:
	case 0xD8:
	case 0xD9:
	case 0xDA:
	case 0xDB:
	case 0xDC:
	case 0xDD:
	case 0xDE:
	case 0xDF:
	case 0xE0: // LDH (n), A
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		gb.WriteOnAddress(gb.Reg.ReadA(), 0xFFFF|uint16(Z))

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xE1: // POP HL
		// this instruction does not do calculations that affect flags
		Z := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		W := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteHL((uint16(W) << 8) | uint16(Z))

	case 0xE2: // LDH [C], A
		// M2
		gb.WriteOnAddress(gb.Reg.ReadA(), 0xFF00|uint16(gb.Reg.ReadC()))

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xE3:
	case 0xE4:
	case 0xE5: // PUSH HL
		// M2
		gb.Reg.DecreaseSP()

		// M3
		gb.WriteOnAddress(byte(gb.Reg.ReadHL()&0xF), gb.Reg.ReadSP())
		gb.Reg.DecreaseSP()

		// M4
		gb.WriteOnAddress(byte(gb.Reg.ReadHL()>>8), gb.Reg.ReadSP())

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xE6:
	case 0xE7:
	case 0xE8:
	case 0xE9:
	case 0xEA: // LD (nn), A
		// M2
		Z := uint16(gb.ReadByAddress(gb.Reg.ReadPC()))
		gb.Reg.IncreasePC()

		// M3
		W := uint16(gb.ReadByAddress(gb.Reg.ReadPC()))
		gb.Reg.IncreasePC()

		// M4
		gb.WriteOnAddress(gb.Reg.ReadA(), (W<<8)|Z)

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xEB:
	case 0xEC:
	case 0xED:
	case 0xEE:
	case 0xEF:
	case 0xF0: // LDH A, [C]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		Z = gb.ReadByAddress(0xFF00 | uint16(Z))

		// M4/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
		gb.Reg.WriteA(Z)

	case 0xF1: // POP AF
		// this instruction does not do calculations that affect flags, BUT
		// replaces all flags from what is read from memory
		Z := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		W := gb.ReadByAddress(gb.Reg.ReadSP())
		gb.Reg.IncreaseSP()

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteAF((uint16(W) << 8) | uint16(Z))

	case 0xF2: // LDH A, [C]
		// M2
		Z := gb.ReadByAddress(0xFF00 | uint16(gb.Reg.ReadC()))

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
		gb.Reg.WriteA(Z)

	case 0xF3:
	case 0xF4:
	case 0xF5: // PUSH AF
		// M2
		gb.Reg.DecreaseSP()

		// M3
		gb.WriteOnAddress(byte(gb.Reg.ReadAF()&0xF), gb.Reg.ReadSP())
		gb.Reg.DecreaseSP()

		// M4
		gb.WriteOnAddress(byte(gb.Reg.ReadAF()>>8), gb.Reg.ReadSP())

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xF6:
	case 0xF7:
	case 0xF8: // LD HL, SP+e
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		result, carry := AddUnsigned8(uint8(gb.Reg.ReadSP()>>8), Z)

		gb.Reg.WriteL(result)

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag((carry>>3)&0b1 == 0b1)
		gb.Reg.SetCarryFlag((carry>>7)&0b1 == 0b1)

		// M4/M1
		result = uint8(gb.Reg.ReadHL())
		if (Z >> 7) == 0b1 {
			result = result + 0xFF
		}
		if gb.Reg.ReadCarryFlag() {
			result = result + 0b1
		}

		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteH(result)

	case 0xF9: // LD SP, HL
		// M2
		gb.Reg.WriteSP(gb.Reg.ReadHL())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xFA: // LD A, (nn)
		// M2
		Z := uint16(gb.ReadByAddress(gb.Reg.ReadPC()))
		gb.Reg.IncreasePC()

		// M3
		W := uint16(gb.ReadByAddress(gb.Reg.ReadPC()))
		gb.Reg.IncreasePC()

		// M4
		value := gb.ReadByAddress((W << 8) | Z)

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
		gb.Reg.WriteA(value)

	case 0xFB:
	case 0xFC:
	case 0xFD:
	case 0xFE:
	case 0xFF:

	}

}

func NewGameBoy() GameBoy {
	return GameBoy{}
}
