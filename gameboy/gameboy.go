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
func AddUnsigned8(a, b uint8) (result uint8, halfCarry, carry bool) {
	carry_per_bit := a & b
	return a + b, (carry_per_bit>>3)&0b1 == 1, (carry_per_bit >> 7 & 0b1) == 1
}

func SubstractUnsigned8(a, b uint8) (result uint8, halfCarry, carry bool) {
	return a - b, (a & 0x0F) < (b & 0x0F), a < b
}

func SubstractUnsigned8Carry(a, b uint8, carry bool) (result uint8, halfCarry, carryOut bool) {
	subtrahend := uint16(b)
	lowSubtrahend := (b & 0x0F)
	if carry {
		subtrahend = subtrahend + 1
		lowSubtrahend = lowSubtrahend + 1
	}

	return uint8(uint16(a) - subtrahend), (a & 0x0F) < lowSubtrahend, uint16(a) < subtrahend
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

// RotateLeft shifts the input value 1 bit to the left, and puts the carry flag
// into the bit 0 position
func RotateLeft(a uint8, carry bool) (result uint8, top_bit bool) {
	top_bit = (a >> 7) == 0b1
	result = a << 1
	if carry {
		result = result | 0b1
	}

	return result, top_bit
}

// RotateRight shifts the input value 1 bit to the right, and puts the carry flag
// into the bit 0 position
func RotateRight(a uint8, carry bool) (result uint8, bottom_bit bool) {
	bottom_bit = (a & 0b1) == 0b1
	result = a >> 1
	if carry {
		result = result | (0b1 << 7)
	}

	return result, bottom_bit
}

// RotateLeftCircular shifts the input value 1 bit to the left, and puts the original
// 7th bit in the place of the new 0th bit
func RotateLeftCircular(a uint8) (result uint8, top_bit bool) {
	top_bit = (a >> 7) == 0b1
	result = (a << 1) | (a >> 7)

	return result, top_bit
}

// RotateRightCircular shifts the input value 1 bit to the right, and puts the original
// 0th bit in the place of the new 7th bit
func RotateRightCircular(a uint8) (result uint8, bottom_bit bool) {
	bottom_bit = (a & 0b1) == 0b1
	result = (a >> 1) | (a << 7)

	return result, bottom_bit
}

// ShiftLeftArithmetic shifts the input value 1 bit to the left. The 7th bit becomes
// the carry flag; the 0th bit is always 0.
func ShiftLeftArithmetic(a uint8) (result uint8, carry bool) {
	carry = (a >> 7) == 0b1
	result = a << 1

	return result, carry
}

// ShiftRightArithmetic shifts the input value 1 bit to the right. The 0th bit becomes
// the carry flag; the 7th bit remains unchanged.
func ShiftRightArithmetic(a uint8) (result uint8, carry bool) {
	carry = (a & 0b1) == 0b1
	result = (a >> 1) | (a & 0b10000000)

	return result, carry
}

// ShiftRightLogically shifts the input value 1 bit to the right. The 0th bit becomes
// the carry flag; the 7th bit is always 0.
func ShiftRightLogically(a uint8) (result uint8, carry bool) {
	carry = (a & 0b1) == 0b1
	result = (a >> 1)

	return result, carry
}

func SwapNibbles(a uint8) uint8 {
	return (a >> 4) | ((a & 0xF) << 4)
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

	case 0x03: // INC BC
		// M2
		gb.Reg.WriteBC(gb.Reg.ReadBC() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x04: // INC B
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadB(), 1)
		gb.Reg.WriteB(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x05: // DEC B
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadB(), 1)
		gb.Reg.WriteB(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x07: // RLCA
		// M2/M1
		result, carry := RotateLeftCircular(gb.Reg.ReadA())

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(result)

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

	case 0x09: // ADD HL, BC
		// M2
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadL(), gb.Reg.ReadC())
		gb.Reg.WriteL(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		// M3/M1
		result, halfCarry, carry = AddUnsigned8Carry(gb.Reg.ReadH(), gb.Reg.ReadB(), carry)
		gb.Reg.WriteH(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x0A: // LD A, [BC]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadBC())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(Z)

	case 0x0B: // DEC BC
		// M2
		gb.Reg.WriteBC(gb.Reg.ReadBC() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x0C: // INC C
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadC(), 1)
		gb.Reg.WriteC(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x0D:
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadC(), 1)
		gb.Reg.WriteC(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x0F: // RRCA
		// M2/M1
		result, carry := RotateRightCircular(gb.Reg.ReadA())

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(result)

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

	case 0x13: // INC DE
		// M2
		gb.Reg.WriteDE(gb.Reg.ReadDE() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x14: // INC D
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadD(), 1)
		gb.Reg.WriteD(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x15: // DEC D
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadD(), 1)
		gb.Reg.WriteD(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x17: // RLA
		// M2/M1
		result, carry := RotateLeft(gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(result)

	case 0x18:
		return
	case 0x19: // ADD HL, DE
		// M2
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadL(), gb.Reg.ReadE())
		gb.Reg.WriteL(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		// M3/M1
		result, halfCarry, carry = AddUnsigned8Carry(gb.Reg.ReadH(), gb.Reg.ReadD(), carry)
		gb.Reg.WriteH(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x1A: // LD A, [DE]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadDE())

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(Z)

	case 0x1B: // DEC DE
		// M2
		gb.Reg.WriteDE(gb.Reg.ReadDE() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x1C: // INC E
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadE(), 1)
		gb.Reg.WriteE(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x1D: // DEC E
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadE(), 1)
		gb.Reg.WriteE(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x1F: // RRA
		// M2/M1
		result, carry := RotateRight(gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(result)

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

	case 0x23: // INC HL
		// M2
		gb.Reg.WriteHL(gb.Reg.ReadHL() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x24: // INC H
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadH(), 1)
		gb.Reg.WriteH(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x25: // DEC H
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadH(), 1)
		gb.Reg.WriteH(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x27: // DAA
		// https://blog.ollien.com/posts/gb-daa
		// M2/M1
		var offset uint16
		var shouldCarry bool = false

		a_value := gb.Reg.ReadA()
		if (!gb.Reg.ReadSubstractionFlag() && ((a_value & 0xF) > 0x09)) || gb.Reg.ReadHalfCarryFlag() {
			offset = offset | 0x06
		}

		if (!gb.Reg.ReadSubstractionFlag() && (a_value > 0x99)) || gb.Reg.ReadCarryFlag() {
			offset = offset | 0x60
			shouldCarry = true
		}

		var result uint8
		if !gb.Reg.ReadSubstractionFlag() {
			result = uint8(uint16(a_value) + offset)
		} else {
			result = uint8(uint16(a_value) - offset)
		}

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(shouldCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteA(result)

	case 0x28:
		return
	case 0x29: // ADD HL, HL
		// M2
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadL(), gb.Reg.ReadL())
		gb.Reg.WriteL(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		// M3/M1
		result, halfCarry, carry = AddUnsigned8Carry(gb.Reg.ReadH(), gb.Reg.ReadH(), carry)
		gb.Reg.WriteH(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x2B: // DEC HL
		// M2
		gb.Reg.WriteHL(gb.Reg.ReadHL() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x2C: // INC L
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadL(), 1)
		gb.Reg.WriteL(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x2D: // DEC L
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadL(), 1)
		gb.Reg.WriteL(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x2F: // CPL
		// M2/M1
		gb.Reg.WriteA(^gb.Reg.ReadA())
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(true)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x33: // INC SP
		// M2
		gb.Reg.WriteSP(gb.Reg.ReadSP() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x34: // INC [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, halfCarry, _ := AddUnsigned8(Z, 1)
		gb.WriteOnAddress(result, gb.Reg.ReadHL())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x35:
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, halfCarry, _ := SubstractUnsigned8(Z, 1)
		gb.WriteOnAddress(result, gb.Reg.ReadHL())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x37: // SCF
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(true)
		opcode, interrupt := gb.FetchCycle(gb.Reg.ReadPC())
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x38:
		return
	case 0x39: // ADD HL, SP
		// M2
		sp_value := gb.Reg.ReadSP()
		//
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadL(), uint8(sp_value&0xFF))
		gb.Reg.WriteL(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		// M3/M1
		result, halfCarry, carry = AddUnsigned8Carry(gb.Reg.ReadH(), uint8(sp_value>>8), carry)
		gb.Reg.WriteH(result)

		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x3B: // DEC SP
		// M2
		gb.Reg.WriteSP(gb.Reg.ReadSP() + 1)

		// M3/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x3C: // INC A
		// M2/M1
		result, halfCarry, _ := AddUnsigned8(gb.Reg.ReadA(), 1)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x3D: // DEC A
		// M2/M1
		result, halfCarry, _ := SubstractUnsigned8(gb.Reg.ReadA(), 1)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x3F: // CCF
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(!gb.Reg.ReadCarryFlag())

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0x76: // HALT TODO
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
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadB())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x81: // ADD A, C
		// M2/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadC())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x82: // ADD A, D
		// M2/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadD())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x83: // ADD A, E
		// M2/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadE())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x84: // ADD A, H
		// M2/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadH())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x85: // ADD A, L
		// M2/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadL())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x86: // ADD A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x87: // ADD A, A
		// M2/M1
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadA())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

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
		result, halfCarry, carry := AddUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x90: // SUB A, B
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadB())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x91: // SUB A, C
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadC())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x92: // SUB A, D
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadD())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x93: // SUB A, E
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadE())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x94: // SUB A, H
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadH())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x95: // SUB A, L
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadL())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x96: // SUB A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x97: // SUB A, A
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadA())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x98: // SBC A, B
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadB(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x99: // SBC A, C
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadC(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x9A: // SBC A, D
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadD(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x9B: // SBC A, E
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadE(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x9C: // SBC A, H
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadH(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x9D: // SBC A, L
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadL(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x9E: // SBC A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())
		// M3/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), Z, gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0x9F: // SBC A, A
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA0: // AND A, B
		// M2/M1
		result := gb.Reg.ReadA() & gb.Reg.ReadB()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA1: // AND A, C
		// M2/M1
		result := gb.Reg.ReadA() & gb.Reg.ReadC()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA2: // AND A, D
		// M2/M1
		result := gb.Reg.ReadA() & gb.Reg.ReadD()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA3: // AND A, E
		// M2/M1
		result := gb.Reg.ReadA() & gb.Reg.ReadE()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA4: // AND A, H
		// M2/M1
		result := gb.Reg.ReadA() & gb.Reg.ReadH()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA5: // AND A, L
		// M2/M1
		result := gb.Reg.ReadA() & gb.Reg.ReadL()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA6: // AND A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result := gb.Reg.ReadA() & Z
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA7: // AND A, A
		// M2/M1
		// result is always equal to A
		result := gb.Reg.ReadA() & gb.Reg.ReadA()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA8: // XOR A, B
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadB()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xA9: // XOR A, C
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadC()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xAA: // XOR A, D
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadD()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xAB: // XOR A, E
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadE()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xAC: // XOR A, H
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadH()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xAD: // XOR A, L
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadL()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xAE: // XOR A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result := gb.Reg.ReadA() ^ Z
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xAF: // XOR A, A
		// M2/M1
		result := gb.Reg.ReadA() ^ gb.Reg.ReadA()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB0: // OR A, B
		// M2/M1
		result := gb.Reg.ReadA() | gb.Reg.ReadB()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB1: // OR A, C
		// M2/M1
		result := gb.Reg.ReadA() | gb.Reg.ReadC()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB2: // OR A, D
		// M2/M1
		result := gb.Reg.ReadA() | gb.Reg.ReadD()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB3: // OR A, E
		// M2/M1
		result := gb.Reg.ReadA() | gb.Reg.ReadE()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB4: // OR A, H
		// M2/M1
		result := gb.Reg.ReadA() | gb.Reg.ReadH()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB5: // OR A, L
		// M2/M1
		result := gb.Reg.ReadA() | gb.Reg.ReadL()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB6: // OR A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result := gb.Reg.ReadA() | Z
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB7: // OR A, A
		// M2/M1

		// does nothing
		result := gb.Reg.ReadA() | gb.Reg.ReadA()
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB8: // CP A, B
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadB())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xB9: // CP A, C
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadC())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xBA: // CP A, D
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadD())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xBB: // CP A, E
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadE())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xBC: // CP A, H
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadH())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xBD: // CP A, L
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadL())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xBE: // CP A, [HL]
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadHL())

		// M3/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xBF: // CP A, A
		// M2/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), gb.Reg.ReadA())
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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
		result, halfCarry, carry := AddUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xC7:
	case 0xC8:
	case 0xC9:
	case 0xCA:
	case 0xCB: // PREFIX
		// M2
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		switch opcode {
		case 0x00: // RLC B
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadB())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x01: // RLC C
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadC())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x02: // RLC D
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadD())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x03: // RLC E
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadE())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x04: // RLC H
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadH())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x05: // RLC L
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadL())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x06: // RLC [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := RotateLeftCircular(Z)

			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x07: // RLC A
			// M3/M1
			result, carry := RotateLeftCircular(gb.Reg.ReadA())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x08: // RRC B
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadB())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x09: // RRC C
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadC())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x0A: // RRC D
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadD())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x0B: // RRC E
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadE())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x0C: // RRC H
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadH())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x0D: // RRC L
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadL())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x0E: // RRC [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := RotateRightCircular(Z)

			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x0F: // RRC A
			// M3/M1
			result, carry := RotateRightCircular(gb.Reg.ReadA())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x10: // RL B
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadB(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x11: // RL C
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadC(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x12: // RL D
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadD(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x13: // RL E
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadE(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x14: // RL H
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadH(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x15: // RL L
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadL(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x16: // RL [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := RotateLeft(Z, gb.Reg.ReadCarryFlag())

			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x17: // RL A
			// M3/M1
			result, carry := RotateLeft(gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x18: // RR B
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadB(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x19: // RR C
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadC(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x1A: // RR D
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadD(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x1B: // RR E
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadE(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x1C: // RR H
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadH(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x1D: // RR L
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadL(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x1E: // RR [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := RotateRight(Z, gb.Reg.ReadCarryFlag())

			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x1F: // RR A
			// M3/M1
			result, carry := RotateRight(gb.Reg.ReadA(), gb.Reg.ReadCarryFlag())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x20: // SLA B
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadB())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x21: // SLA C
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadC())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x22: // SLA D
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadD())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x23: // SLA E
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadE())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x24: // SLA H
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadH())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x25: // SLA L
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadL())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x26: // SLA [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := ShiftLeftArithmetic(Z)
			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x27: // SLA A
			// M3/M1
			result, carry := ShiftLeftArithmetic(gb.Reg.ReadA())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x28: // SRA B
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadB())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x29: // SRA C
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadC())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x2A: // SRA D
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadD())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x2B: // SRA E
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadE())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x2C: // SRA H
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadH())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x2D: // SRA L
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadL())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x2E: // SRA [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := ShiftRightArithmetic(Z)
			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x2F: // SRA A
			// M3/M1
			result, carry := ShiftRightArithmetic(gb.Reg.ReadA())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x30: // SWAP B
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadB())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x31: // SWAP C
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadC())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x32: // SWAP D
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadD())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x33: // SWAP E
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadE())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x34: // SWAP H
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadH())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x35: // SWAP L
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadL())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x36: // SWAP [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result := SwapNibbles(Z)
			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x37: // SWAP A
			// M3/M1
			result := SwapNibbles(gb.Reg.ReadA())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(false)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x38: // SRL B
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadB())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteB(result)

		case 0x39: // SRL C
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadC())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteC(result)

		case 0x3A: // SRL D
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadD())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteD(result)

		case 0x3B: // SRL E
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadE())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteE(result)

		case 0x3C: // SRL H
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadH())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteH(result)

		case 0x3D: // SRL L
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadL())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteL(result)

		case 0x3E: // SRL [HL]
			// M3
			Z := gb.ReadByAddress(gb.Reg.ReadHL())

			// M4
			result, carry := ShiftRightLogically(Z)
			// M5/M1
			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.WriteOnAddress(result, gb.Reg.ReadHL())

		case 0x3F: // SRL A
			// M3/M1
			result, carry := ShiftRightLogically(gb.Reg.ReadA())

			gb.Reg.SetZeroFlag(result == 0)
			gb.Reg.SetSubstractionFlag(false)
			gb.Reg.SetHalfCarryFlag(false)
			gb.Reg.SetCarryFlag(carry)

			opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
			gb.Reg.IR = opcode
			gb.Reg.IE = interrupt
			gb.Reg.IncreasePC()

			gb.Reg.WriteA(result)

		case 0x40: // BIT 0, B
		case 0x41: // BIT 0, C
		case 0x42: // BIT 0, D
		case 0x43: // BIT 0, E
		case 0x44: // BIT 0, H
		case 0x45: // BIT 0, L
		case 0x46: // BIT 0, [HL]
		case 0x47: // BIT 0, A
		case 0x48: // BIT 1, B
		case 0x49: // BIT 1, C
		case 0x4A: // BIT 1, D
		case 0x4B: // BIT 1, E
		case 0x4C: // BIT 1, H
		case 0x4D: // BIT 1, L
		case 0x4E: // BIT 1, [HL]
		case 0x4F: // BIT 1, A
		case 0x50: // BIT 2, B
		case 0x51: // BIT 2, C
		case 0x52: // BIT 2, D
		case 0x53: // BIT 2, E
		case 0x54: // BIT 2, H
		case 0x55: // BIT 2, L
		case 0x56: // BIT 2, [HL]
		case 0x57: // BIT 2, A
		case 0x58: // BIT 3, B
		case 0x59: // BIT 3, C
		case 0x5A: // BIT 3, D
		case 0x5B: // BIT 3, E
		case 0x5C: // BIT 3, H
		case 0x5D: // BIT 3, L
		case 0x5E: // BIT 3, [HL]
		case 0x5F: // BIT 3, A
		case 0x60: // BIT 4, B
		case 0x61: // BIT 4, C
		case 0x62: // BIT 4, D
		case 0x63: // BIT 4, E
		case 0x64: // BIT 4, H
		case 0x65: // BIT 4, L
		case 0x66: // BIT 4, [HL]
		case 0x67: // BIT 4, A
		case 0x68: // BIT 5, B
		case 0x69: // BIT 5, C
		case 0x6A: // BIT 5, D
		case 0x6B: // BIT 5, E
		case 0x6C: // BIT 5, H
		case 0x6D: // BIT 5, L
		case 0x6E: // BIT 5, [HL]
		case 0x6F: // BIT 5, A
		case 0x70: // BIT 6, B
		case 0x71: // BIT 6, C
		case 0x72: // BIT 6, D
		case 0x73: // BIT 6, E
		case 0x74: // BIT 6, H
		case 0x75: // BIT 6, L
		case 0x76: // BIT 6, [HL]
		case 0x77: // BIT 6, A
		case 0x78: // BIT 7, B
		case 0x79: // BIT 7, C
		case 0x7A: // BIT 7, D
		case 0x7B: // BIT 7, E
		case 0x7C: // BIT 7, H
		case 0x7D: // BIT 7, L
		case 0x7E: // BIT 7, [HL]
		case 0x7F: // BIT 7, A
		case 0x80: // RES 0, B
		case 0x81: // RES 0, C
		case 0x82: // RES 0, D
		case 0x83: // RES 0, E
		case 0x84: // RES 0, H
		case 0x85: // RES 0, L
		case 0x86: // RES 0, [HL]
		case 0x87: // RES 0, A
		case 0x88: // RES 1, B
		case 0x89: // RES 1, C
		case 0x8A: // RES 1, D
		case 0x8B: // RES 1, E
		case 0x8C: // RES 1, H
		case 0x8D: // RES 1, L
		case 0x8E: // RES 1, [HL]
		case 0x8F: // RES 1, A
		case 0x90: // RES 2, B
		case 0x91: // RES 2, C
		case 0x92: // RES 2, D
		case 0x93: // RES 2, E
		case 0x94: // RES 2, H
		case 0x95: // RES 2, L
		case 0x96: // RES 2, [HL]
		case 0x97: // RES 2, A
		case 0x98: // RES 3, B
		case 0x99: // RES 3, C
		case 0x9A: // RES 3, D
		case 0x9B: // RES 3, E
		case 0x9C: // RES 3, H
		case 0x9D: // RES 3, L
		case 0x9E: // RES 3, [HL]
		case 0x9F: // RES 3, A
		case 0xA0: // RES 4, B
		case 0xA1: // RES 4, C
		case 0xA2: // RES 4, D
		case 0xA3: // RES 4, E
		case 0xA4: // RES 4, H
		case 0xA5: // RES 4, L
		case 0xA6: // RES 4, [HL]
		case 0xA7: // RES 4, A
		case 0xA8: // RES 5, B
		case 0xA9: // RES 5, C
		case 0xAA: // RES 5, D
		case 0xAB: // RES 5, E
		case 0xAC: // RES 5, H
		case 0xAD: // RES 5, L
		case 0xAE: // RES 5, [HL]
		case 0xAF: // RES 5, A
		case 0xB0: // RES 6, B
		case 0xB1: // RES 6, C
		case 0xB2: // RES 6, D
		case 0xB3: // RES 6, E
		case 0xB4: // RES 6, H
		case 0xB5: // RES 6, L
		case 0xB6: // RES 6, [HL]
		case 0xB7: // RES 6, A
		case 0xB8: // RES 7, B
		case 0xB9: // RES 7, C
		case 0xBA: // RES 7, D
		case 0xBB: // RES 7, E
		case 0xBC: // RES 7, H
		case 0xBD: // RES 7, L
		case 0xBE: // RES 7, [HL]
		case 0xBF: // RES 7, A
		case 0xC0: // SET 0, B
		case 0xC1: // SET 0, C
		case 0xC2: // SET 0, D
		case 0xC3: // SET 0, E
		case 0xC4: // SET 0, H
		case 0xC5: // SET 0, L
		case 0xC6: // SET 0, [HL]
		case 0xC7: // SET 0, A
		case 0xC8: // SET 1, B
		case 0xC9: // SET 1, C
		case 0xCA: // SET 1, D
		case 0xCB: // SET 1, E
		case 0xCC: // SET 1, H
		case 0xCD: // SET 1, L
		case 0xCE: // SET 1, [HL]
		case 0xCF: // SET 1, A
		case 0xD0: // SET 2, B
		case 0xD1: // SET 2, C
		case 0xD2: // SET 2, D
		case 0xD3: // SET 2, E
		case 0xD4: // SET 2, H
		case 0xD5: // SET 2, L
		case 0xD6: // SET 2, [HL]
		case 0xD7: // SET 2, A
		case 0xD8: // SET 3, B
		case 0xD9: // SET 3, C
		case 0xDA: // SET 3, D
		case 0xDB: // SET 3, E
		case 0xDC: // SET 3, H
		case 0xDD: // SET 3, L
		case 0xDE: // SET 3, [HL]
		case 0xDF: // SET 3, A
		case 0xE0: // SET 4, B
		case 0xE1: // SET 4, C
		case 0xE2: // SET 4, D
		case 0xE3: // SET 4, E
		case 0xE4: // SET 4, H
		case 0xE5: // SET 4, L
		case 0xE6: // SET 4, [HL]
		case 0xE7: // SET 4, A
		case 0xE8: // SET 5, B
		case 0xE9: // SET 5, C
		case 0xEA: // SET 5, D
		case 0xEB: // SET 5, E
		case 0xEC: // SET 5, H
		case 0xED: // SET 5, L
		case 0xEE: // SET 5, [HL]
		case 0xEF: // SET 5, A
		case 0xF0: // SET 6, B
		case 0xF1: // SET 6, C
		case 0xF2: // SET 6, D
		case 0xF3: // SET 6, E
		case 0xF4: // SET 6, H
		case 0xF5: // SET 6, L
		case 0xF6: // SET 6, [HL]
		case 0xF7: // SET 6, A
		case 0xF8: // SET 7, B
		case 0xF9: // SET 7, C
		case 0xFA: // SET 7, D
		case 0xFB: // SET 7, E
		case 0xFC: // SET 7, H
		case 0xFD: // SET 7, L
		case 0xFE: // SET 7, [HL]
		case 0xFF: // SET 7, A

		}

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

	case 0xD6: // SUB A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xD7:
	case 0xD8:
	case 0xD9:
	case 0xDA:
	case 0xDB:
	case 0xDC:
	case 0xDD:
	case 0xDE: // SBC A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result, halfCarry, carry := SubstractUnsigned8Carry(gb.Reg.ReadA(), Z, gb.Reg.ReadCarryFlag())
		gb.Reg.WriteA(result)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

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

	case 0xE6: // AND A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result := gb.Reg.ReadA() & Z
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(true)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xE7:
	case 0xE8: // ADD SP, e
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		sp_value := gb.Reg.ReadSP()
		result, halfCarry, carry := AddUnsigned8(uint8(sp_value&0xFF), Z)
		Z = result

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		// M4
		result = uint8(sp_value >> 8) // MSB
		if (Z >> 7) == 0b1 {
			result = result + 0xFF
		}
		if gb.Reg.ReadCarryFlag() {
			result = result + 0b1
		}

		// M5/M1
		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

		gb.Reg.WriteSP((uint16(result) << 8) | uint16(Z))

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
	case 0xEE: // XOR A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result := gb.Reg.ReadA() ^ Z
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()
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

	case 0xF6: // OR A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3/M1
		result := gb.Reg.ReadA() | Z
		gb.Reg.WriteA(result)

		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(false)
		gb.Reg.SetCarryFlag(false)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xF7:
	case 0xF8: // LD HL, SP+e
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()

		// M3
		result, halfCarry, carry := AddUnsigned8(uint8(gb.Reg.ReadSP()>>8), Z)

		gb.Reg.WriteL(result)

		gb.Reg.SetZeroFlag(false)
		gb.Reg.SetSubstractionFlag(false)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

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
	case 0xFE: // CP A, n8
		// M2
		Z := gb.ReadByAddress(gb.Reg.ReadPC())
		gb.Reg.IncreasePC()
		// M3/M1
		result, halfCarry, carry := SubstractUnsigned8(gb.Reg.ReadA(), Z)
		gb.Reg.SetZeroFlag(result == 0)
		gb.Reg.SetSubstractionFlag(true)
		gb.Reg.SetHalfCarryFlag(halfCarry)
		gb.Reg.SetCarryFlag(carry)

		opcode, interrupt := gb.FetchCycle(gb.Reg.PC)
		gb.Reg.IR = opcode
		gb.Reg.IE = interrupt
		gb.Reg.IncreasePC()

	case 0xFF:

	}

}

func NewGameBoy() GameBoy {
	return GameBoy{}
}
