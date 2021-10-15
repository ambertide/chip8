package device

import (
	"log"
	"time"
)

// This type of function is used to operate on register values.
type OperationFunction func(byte, byte) byte

// This type of function is used to determine VF.
type CarryFunction func(byte, byte) byte

type chip8Registers struct {
	generalPurpose [16]byte
	iRegister      uint16
	programCounter uint16
	delayTimer     byte
	soundTimer     byte
	soundBuffer    *bool
}

// Write to a general purpose register.
func (r *chip8Registers) WriteRegister(registerIndex uint8, value byte) {
	r.generalPurpose[registerIndex] = value
}

// Read from a general purpose register.
func (r *chip8Registers) ReadRegister(registerIndex uint8) byte {
	return r.generalPurpose[registerIndex]
}

// Read register as BCD array, ie: 100 first then 10s then 1s.
func (r *chip8Registers) ReadRegisterBCD(registerIndex uint8) [3]byte {
	value := r.ReadRegister(registerIndex)
	return [3]byte{value / 100, (value / 10) % 10, value % 10}
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

// Write to the I register.
func (r *chip8Registers) WriteIRegister(value uint16) {
	r.iRegister = value
}

// Read the I register.
func (r *chip8Registers) ReadIRegister() uint16 {
	return r.iRegister
}

// Add the value of the source register to the I register.
func (r *chip8Registers) AccumulateIRegister(sourceRegister uint8) {
	r.iRegister += uint16(r.ReadRegister(sourceRegister))
}

// Set the value of I to the address of the sprite representing
// the digit stored in the source register.
func (r *chip8Registers) SetIDigitSprite(sourceRegister uint8) {
	characterIndex := r.ReadRegister(sourceRegister)
	log.Printf("Getting character index for %X\n", characterIndex)
	// Since characters consist of 5 bytes in their sprites,
	// Just times 5 should work.
	log.Printf("Writing address %03X\n to I register.\n", uint16(characterIndex)*5)
	r.WriteIRegister(uint16(characterIndex) * 5)
}

// Set the carry register VF to 1 if value is true.
func (r *chip8Registers) SetCarry(value bool) {
	var byteValue byte = 0x0
	if value {
		byteValue = 0x1
	}
	r.WriteRegister(15, byteValue)
}

// Set the sound timer to the value of the source register.
func (r *chip8Registers) SetSoundTimer(sourceRegister uint8) {
	r.soundTimer = r.ReadRegister(sourceRegister)
}

// Set the delay timer to the value of the source register.
func (r *chip8Registers) SetDelayTimer(sourceRegister uint8) {
	r.delayTimer = r.ReadRegister(sourceRegister)
}

// Set the value of the destination register to the value of the
// Sound timer.
func (r *chip8Registers) LoadSoundTimer(destinationRegister uint8) {
	r.WriteRegister(destinationRegister, r.soundTimer)
}

// Set the value of the destination register to the value of the
// Delay timer.
func (r *chip8Registers) LoadDelayTimer(destinationRegister uint8) {
	r.WriteRegister(destinationRegister, r.delayTimer)
}

// Write an array of bytes to the registers
func (r *chip8Registers) BlockWriteRegisters(registerData [16]byte, length uint8) {
	copy(r.generalPurpose[:length], registerData[:length])
}

// Return all the registers.
func (r *chip8Registers) BlockReadRegisters() [16]byte {
	var buffer [16]byte
	copy(buffer[:], r.generalPurpose[:])
	return buffer
}

// Decrement ST and DT.
func (r *chip8Registers) UpdateClockRegisters() {
	if r.soundTimer > 0 {
		r.soundTimer--
	}
	if r.delayTimer > 0 {
		r.delayTimer--
	}
	*r.soundBuffer = r.soundTimer > 0
}

func (r *chip8Registers) RegisterClockLoop() {
	for {
		r.UpdateClockRegisters()
		time.Sleep(time.Second / 60)
	}
}

// Initialise a new register with a sound buffer.
func NewRegisters(soundBuffer *bool) *chip8Registers {
	register := new(chip8Registers)
	register.soundBuffer = soundBuffer
	return register
}
