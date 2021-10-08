package device

// This type of function is used to operate on register values.
type OperationFunction func(byte, byte) byte

// This type of function is used to determine VF.
type CarryFunction func(byte, byte) byte

type chip8Registers struct {
	generalPurpose [16]byte
	programCounter uint16
	delayTimer     byte
	soundTimer     byte
}

// Write to a general purpose register.
func (r *chip8Registers) WriteRegister(registerIndex uint8, value byte) {
	r.generalPurpose[registerIndex] = value
}

// Read from a general purpose register.
func (r *chip8Registers) ReadRegister(registerIndex uint8) byte {
	return r.generalPurpose[registerIndex]
}

// Add to the value of the register with the given
// index an immediate value and store it in that
// register.
func (r *chip8Registers) AddRegisterImmediate(registerIndex uint8, value byte) {
	r.generalPurpose[registerIndex] += value
}

// Given indexes of two registers, use their values in the operation function,
// Store the return value in the register x. Carry function
// is calculated BEFORE the operation and its result is stored in VF register.
func (r *chip8Registers) RegisterOperationWithCarry(x uint8, y uint8, operation OperationFunction, carry CarryFunction) {
	r.generalPurpose[15] = carry(r.generalPurpose[x], r.generalPurpose[y])
	r.generalPurpose[x] = operation(r.generalPurpose[x], r.generalPurpose[y])
}

// Given indexes of two registers, use their values in the operation function,
// Store the return value in the register x.
func (r *chip8Registers) RegisterOperation(x uint8, y uint8, operation OperationFunction) {

	r.generalPurpose[x] = operation(r.generalPurpose[x], r.generalPurpose[y])
}

// Set the program counter to a value.
func (r *chip8Registers) SetProgramCounter(value uint16) {
	r.programCounter = value
}

// Increment the program counter to the location of
// the next instruction.
func (r *chip8Registers) IncrementProgramCounter() {
	r.programCounter += 2
}

// Return the value of the program counter.
func (r *chip8Registers) GetProgramCounter() uint16 {
	return r.programCounter
}

// Compare two registers and return true if their values are equal.
func (r *chip8Registers) CompareRegisters(x uint8, y uint8) bool {
	return r.generalPurpose[x] == r.generalPurpose[y]
}
