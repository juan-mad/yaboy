package gameboy

type CPURegisters struct {
	AF uint16
	BC uint16
	DE uint16
	HL uint16
	SP uint16
	PC uint16
	IR uint8 // Instruction Register
	IE uint8 // Interrupt Enable
}

func (r CPURegisters) ReadA() byte {
	return byte(r.AF >> 8)
}

func (r *CPURegisters) WriteA(b byte) {
	r.AF = (uint16(b) << 8) | (r.AF & 0x00FF)
}

func (r CPURegisters) ReadF() byte {
	return byte(r.AF & 0x00FF)
}

func (r *CPURegisters) WriteF(b byte) {
	r.AF = (r.AF & 0xFF00) | uint16(b)
}

func (r CPURegisters) ReadB() byte {
	return byte(r.BC >> 8)
}

func (r *CPURegisters) WriteB(b byte) {
	r.BC = (uint16(b) << 8) | (r.BC & 0xFF)
}

func (r CPURegisters) ReadC() byte {
	return byte(r.BC & 0xFF)
}

func (r *CPURegisters) WriteC(b byte) {
	r.BC = (r.BC & 0xFF00) | uint16(b)
}

func (r CPURegisters) ReadD() byte {
	return byte(r.DE >> 8)
}

func (r *CPURegisters) WriteD(b byte) {
	r.DE = (uint16(b) << 8) | (r.DE & 0xFF)
}

func (r CPURegisters) ReadE() byte {
	return byte(r.DE & 0xFF)
}

func (r *CPURegisters) WriteE(b byte) {
	r.DE = (r.DE & 0xFF00) | uint16(b)
}

func (r CPURegisters) ReadH() byte {
	return byte(r.HL >> 8)
}

func (r *CPURegisters) WriteH(b byte) {
	r.HL = (uint16(b) << 8) | (r.HL & 0xFF)
}

func (r CPURegisters) ReadL() byte {
	return byte(r.HL & 0xFF)
}

func (r *CPURegisters) WriteL(b byte) {
	r.HL = (r.HL & 0xFF00) | uint16(b)
}

func (r CPURegisters) ReadAF() uint16 {
	return r.AF
}

func (r *CPURegisters) WriteAF(bb uint16) {
	r.AF = bb
}

func (r CPURegisters) ReadBC() uint16 {
	return r.BC
}

func (r *CPURegisters) WriteBC(bb uint16) {
	r.BC = bb
}

func (r CPURegisters) ReadDE() uint16 {
	return r.DE
}

func (r *CPURegisters) WriteDE(bb uint16) {
	r.DE = bb
}

func (r CPURegisters) ReadHL() uint16 {
	return r.HL
}

func (r *CPURegisters) WriteHL(bb uint16) {
	r.HL = bb
}

func (r CPURegisters) ReadPC() uint16 {
	return r.PC
}

func (r *CPURegisters) IncreasePC() {
	r.PC = r.PC + 1
}

func (r *CPURegisters) WritePC(bb uint16) {
	r.PC = bb
}

func (r CPURegisters) ReadSP() uint16 {
	return r.SP
}

func (r *CPURegisters) WriteSP(bb uint16) {
	r.SP = bb
}

// ReadZeroFlag returns whether the 7th bit of register F is set.
// This bit contains the Zero flag, and is set iff the result of an operation is zero.
// Used by conditional jumps.
func (r CPURegisters) ReadZeroFlag() bool {
	return (0b10000000 & r.AF) == 0b10000000
}

// SetZeroFlag sets the zero flag to the specified state. It should be set to 1 iff
// the result of an operation is zero.
// Used by conditional jumps.
func (r *CPURegisters) SetZeroFlag(state bool) {
	if state {
		r.AF = r.AF & ((0xFF << 8) | uint16(0b10000000))
	} else {
		r.AF = r.AF & ((0xFF << 8) & uint16(0b01111111))
	}
}

// ReadSubstractionFlag returns the substraction (N) flag.
// It is set iff the previous operation was a substraction.
func (r CPURegisters) ReadSubstractionFlag() bool {
	return (0b01000000 & r.AF) == 0b01000000
}

// SetSubstractionFlag returns the substraction (N) flag.
// It is set iff the previous operation was a substraction.
func (r *CPURegisters) SetSubstractionFlag(state bool) {
	if state {
		r.AF = r.AF & ((0xFF << 8) | uint16(0b01000000))
	} else {
		r.AF = r.AF & ((0xFF << 8) & uint16(0b10111111))
	}
}

// ReadHalfCarryFlag returns the half carry flag.
// It is set iff the previous operation raised a half carry.
func (r CPURegisters) ReadHalfCarryFlag() bool {
	return (0b00100000 & r.AF) == 0b00100000
}

// SetHalfCarryFlag returns the half carry flag.
// It is set iff the previous operation raised a half carry.
func (r CPURegisters) SetHalfCarryFlag(state bool) {
	if state {
		r.AF = r.AF & ((0xFF << 8) | uint16(0b00100000))
	} else {
		r.AF = r.AF & ((0xFF << 8) & uint16(0b11011111))
	}
}

// ReadCarryFlag returns whether the 4th bit of register F is set.
// This bit contains the Carry flag, and is set in these cases:
// - When the result of an 8-bit addition is higher than 0xFF.
// - When the result of a 16-bit addition is higher than 0xFFFF.
// - When the result of a substraction or comparison is lower than zero
// - When a rotate/shift operation shifts out a "1" bit.
// Used by conditional jumps and instructions such as ADC, SBC, RL, RLA, etc.
func (r CPURegisters) ReadCarryFlag() bool {
	return (0b00010000 & r.AF) == 0b00010000
}

// SetCarryFlag sets the 4th bit of register F to the passed state.
// This bit contains the Carry flag, and is set in these cases:
// - When the result of an 8-bit addition is higher than 0xFF.
// - When the result of a 16-bit addition is higher than 0xFFFF.
// - When the result of a substraction or comparison is lower than zero
// - When a rotate/shift operation shifts out a "1" bit.
// Used by conditional jumps and instructions such as ADC, SBC, RL, RLA, etc.
func (r *CPURegisters) SetCarryFlag(state bool) {
	if state {
		r.AF = r.AF & ((0xFF << 8) | uint16(0b00010000))
	} else {
		r.AF = r.AF & ((0xFF << 8) & uint16(0b11101111))
	}
}
