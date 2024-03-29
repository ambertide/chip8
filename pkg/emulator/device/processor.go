package device

import (
	"fmt"
	"math/rand"
	"time"
)

// Convert an instruction to its character version.
func getInstructionChar(instruction uint16) string {
	return fmt.Sprintf("%04X", instruction)
}

// Get register index from its hexadecimal representation.
func getRegisterIndex(hex byte) byte {
	switch {
	case '0' <= hex && hex <= '9':
		return byte(hex) - 48
	case 'a' <= hex && hex <= 'f':
		return byte(hex) - 87
	case 'A' <= hex && hex <= 'F':
		return byte(hex) - 55
	default:
		return 128
	}

}

// Main processor of the Chip-8
type Processor struct {
	memory    *chip8Memory
	registers *chip8Registers
	display   *chip8Display
	stack     *chip8Stack
	keyboards *chip8Keyboard
}

func NewProcessor(screenBuffer *[32]uint64, keyboardBuffer *uint16, soundBuffer *bool) *Processor {
	rand.Seed(time.Now().Unix())
	processor := new(Processor)
	processor.display = newDisplay(screenBuffer)
	//log.Println("Display initialised.")
	processor.memory = newMemory()
	//log.Println("Memory initialised.")
	processor.registers = NewRegisters(soundBuffer)
	//log.Println("Registers initialised.")
	processor.keyboards = NewKeyboard(keyboardBuffer)
	//log.Println("Keyboard initialised.")
	processor.stack = new(chip8Stack)
	//log.Println("Stack initialised.")
	go processor.registers.RegisterClockLoop()
	//log.Println("Register clock loop started.")
	return processor
}

// Load a standard Chip-8 Program to the memory
// And set the program counter accordingly.
func (p *Processor) LoadProgram(program []byte, programSize uint16) {
	// Load the program.
	p.memory.LoadProgram(program, programSize)
	// Set the PC to standard start location.
	p.registers.SetProgramCounter(0x198)
}

// Load an ETI program to the memory and set
// the program counter accordingly.
func (p *Processor) LoadETIProgram(program []byte, programSize uint16) {
	p.memory.LoadETIProgram(program, programSize)
	p.registers.SetProgramCounter(0x600)
}

// Fetch the current instruction.
func (p *Processor) fetchInstruction() uint16 {
	mostSignificantByte := p.memory.ReadMemory(p.registers.GetProgramCounter())
	leastSignificantByte := p.memory.ReadMemory(p.registers.GetProgramCounter() + 1)
	return uint16(mostSignificantByte)<<8 + uint16(leastSignificantByte)
}

// Returns true if the processor should halt.
func (p *Processor) ShouldHalt() bool {
	return p.registers.GetProgramCounter() >= 0xFFD
}

// Run a CPU Fetch/Execute cycle.
func (p *Processor) Cycle() {
	//Increment the PC.
	p.registers.IncrementProgramCounter()
	// Fetch the instruction.
	instruction := p.fetchInstruction()
	// Execute the instruction
	p.executeInstruction(instruction)
}
