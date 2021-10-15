package device

import (
	"log"
	"math/rand"
)

type LogicalInstructionType uint8

const (
	Load LogicalInstructionType = iota
	Or
	And
	Xor
	Add
	Sub
	Shr
	Subn
	Shl = 0xE
)

// Execute system instructions RET and CLR
func (p *Processor) executeSystemInstruction(instruction string) {
	switch instruction {
	case "0000":
		break
	case "00E0":
		// CLR: Clear Screen
		log.Print("INSTRUCTION: Clear the display.\n")
		p.display.ClearDisplay()
	case "00EE":
		// RET: Return from subroutine.
		log.Print("INSTRUCTION: Return from subroutine.\n")
		p.registers.SetProgramCounter(p.stack.Pop())
	default:
		log.Panicf("ERROR: Unknown Instruction %s.", instruction)
	}
}

// Execute skip instructions that skip a number of intsructions
// Depending on a condition. SE, SNE, SE, SNE.
func (p *Processor) executeSkipInstructions(instructionChar string, instruction uint16) {
	byteValue := byte(instruction & 0xFF) // For 3XNN and 4XNN
	register := getRegisterIndex(instructionChar[1])
	register2 := getRegisterIndex(instructionChar[2]) // For 5XY0
	msn := instructionChar[0]                         // Most significant nibble
	lsn := instructionChar[3]                         // Least significant nibble
	switch {
	case msn == '3' && p.registers.ReadRegister(register) == byteValue:
		// Incrementing the program counter now will effectively
		// Skip this instruction.
		log.Print("INSTRUCTION: Skip Variant 1.\n")
		p.registers.IncrementProgramCounter()
	case msn == '4' && p.registers.ReadRegister(register) != byteValue:
		log.Print("INSTRUCTION: Skip Variant 2.\n")
		p.registers.IncrementProgramCounter()
	case msn == '5' && lsn == '0' && p.registers.CompareRegisters(register, register2):
		log.Print("INSTRUCTION: Skip Variant 3.\n")
		p.registers.IncrementProgramCounter()
	case msn == '9' && lsn == '0' && !p.registers.CompareRegisters(register, register2):
		log.Print("INSTRUCTION: Skip Variant 4.\n")
		p.registers.IncrementProgramCounter()
	case msn == 'E' && (instruction<<8 == 0x9E00) && p.keyboards.IsKeyPressed(p.registers.ReadRegister(register)):
		log.Print("INSTRUCTION: Skip Variant 5.\n")
		p.registers.IncrementProgramCounter()
	case msn == 'E' && (instruction<<8 == 0xA100) && !p.keyboards.IsKeyPressed(p.registers.ReadRegister(register)):
		log.Print("INSTRUCTION: Skip Variant 6.\n")
		p.registers.IncrementProgramCounter()
	}
}

// Given the indexes for two registers and the operation type execute
// the operation.
func (p *Processor) executeLogicalInstructions(x uint8, y uint8, operationType LogicalInstructionType, instruction uint16) {
	switch operationType {
	case Load:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Load V%d, V%d.\n", x, y)
		p.registers.RegisterOperation(x, y, func(b1 byte, b2 byte) byte { return b2 })
	case Or:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Or V%d, V%d.\n", x, y)
		p.registers.RegisterOperation(x, y, func(b1 byte, b2 byte) byte { return b1 | b2 })
	case And:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit And V%d, V%d.\n", x, y)
		p.registers.RegisterOperation(x, y, func(b1, b2 byte) byte { return b1 & b2 })
	case Xor:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Xor V%d, V%d.\n", x, y)
		p.registers.RegisterOperation(x, y, func(b1, b2 byte) byte { return b1 ^ b2 })
	case Add:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Add V%d, V%d.\n", x, y)
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 + b2 },
			func(b1, b2 byte) byte {
				if (uint16(b1) + uint16(b2)) > 255 {
					return 1
				} else {
					return 0
				}
			},
		)
	case Sub:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Sub V%d, V%d.\n", x, y)
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 - b2 },
			func(b1, b2 byte) byte {
				if b1 > b2 {
					return 1
				} else {
					return 0
				}
			},
		)
	case Shr:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Shr V%d, V%d.\n", x, y)
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 >> 1 },
			func(b1, b2 byte) byte {
				return b1 & 0x1
			},
		)
	case Subn:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Subn V%d, V%d.\n", x, y)
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b2 - b1 },
			func(b1, b2 byte) byte {
				if b2 > b1 {
					return 1
				} else {
					return 0
				}
			},
		)
	case Shl:
		log.Printf("INSTRUCTION: Arithmatic Logic Unit Shl V%d, V%d.\n", x, y)
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 << 1 },
			func(b1, b2 byte) byte {
				return b1 >> 7
			},
		)
	default:
		log.Panicf("ERROR: Unknown Instruction %04X.", instruction)
	}
}

// Set the value of the register to the immediate ANDed with a
// Randomly generated number.
func (p *Processor) executeRandomAnd(register uint8, immediate byte) {
	log.Print("INSTRUCTION: Execute random number generation.\n")
	randomByte := byte(rand.Intn(255))
	p.registers.WriteRegister(register, randomByte&immediate)
}

// Draw to the screen the sprite loaded from n bytes from the memory address at
// Register I starting from screen coordinates (x, y) set VF to true if
// there is collision.
func (p *Processor) executeDrawInstruction(x uint8, y uint8, n byte) {
	log.Printf("INSTRUCTION: Draw at (V%d, V%d)\n", x, y)
	memoryAddress := p.registers.ReadIRegister()
	spriteData := p.memory.BlockReadFromMemory(memoryAddress, memoryAddress+uint16(n))
	startX, startY := p.registers.ReadRegister(x), p.registers.ReadRegister(y)
	collision := p.display.DrawSprite(startX, startY, n, spriteData)
	p.registers.SetCarry(collision)
}

func (p *Processor) executeRegisterInstructions(register uint8, instruction uint16) {
	// There are 9 subtypes of F operations.
	subtype := instruction & 0xFF
	switch subtype {
	case 0x07:
		// LD: Load delay timer to VX
		log.Print("INSTRUCTION: Store DT.\n")
		p.registers.LoadDelayTimer(register)
	case 0x0A:
		// LD: Wait and load key to VX
		log.Print("INSTRUCTION: Wait for Key.\n")
		p.registers.WriteRegister(register, p.keyboards.WaitForKeyPress())
	case 0x15:
		// LD: Load VX to delay timer
		log.Print("INSTRUCTION: Write DT\n")
		p.registers.SetDelayTimer(register)
	case 0x18:
		// LD: Set Sound timer to VX
		log.Print("INSTRUCTION: Write ST\n")
		p.registers.SetSoundTimer(register)
	case 0x1E:
		// ADD: Accumulate VX to I
		log.Print("INSTRUCTION: Accumulate I\n")
		p.registers.AccumulateIRegister(register)
	case 0x29:
		// LD: Set I to the location for the sprite
		// of the digit in the register.
		log.Print("INSTRUCTION: Set I to sprite.\n")
		p.registers.SetIDigitSprite(register)
	case 0x33:
		// LD: Store BCD representation of VX in memory.
		log.Print("INSTRUCTION: Load BCD\n")
		bcd := p.registers.ReadRegisterBCD(register)
		addrrStart := p.registers.ReadIRegister()
		p.memory.BlockWriteToMemory(addrrStart, addrrStart+3, bcd[:])
	case 0x55:
		// LD: Store registers to memory.
		log.Print("INSTRUCTION: Dump registers to memory.\n")
		addrStart := p.registers.ReadIRegister()
		registers := p.registers.BlockReadRegisters()
		p.memory.BlockWriteToMemory(addrStart, addrStart+uint16(register+1), registers[:register+1])
	case 0x65:
		// LD: Load registers from memory.
		log.Print("INSTRUCTION: Load registers from memory.\n")
		addrStart := p.registers.ReadIRegister()
		registers := p.memory.BlockReadFromMemory(addrStart, addrStart+uint16(register+1))
		var registersCopy [16]byte
		copy(registersCopy[:register+1], registers[:register+1])
		p.registers.BlockWriteRegisters(registersCopy, register+1)
		//fmt.Printf("0x%03X => %#v\n%#v\n%#v", addrStart, registers, registersCopy, p.registers.generalPurpose)
	default:
		log.Panicf("ERROR: Unknown Instruction %04X.", instruction)
	}
}

// Execute the next instruction
func (p *Processor) executeInstruction(instruction uint16) {
	instructionCharacter := getInstructionChar(instruction)
	register := getRegisterIndex(instructionCharacter[1])  // For many instructions.
	register2 := getRegisterIndex(instructionCharacter[2]) // For some instructions.
	immediate := byte(instruction & 0xFF)                  // For some instructions.
	switch instructionCharacter[0] {
	case '0':
		p.executeSystemInstruction(instructionCharacter)
	case '1':
		// JP: Jump to the location.
		log.Print("INSTRUCTION: Jump to location.\n")
		p.registers.SetProgramCounter((instruction & 0x0FFF) - 2)
	case '2':
		// CALL: Call a subroutine.
		log.Print("INSTRUCTION: Call a subroutine.\n")
		p.stack.Push(p.registers.GetProgramCounter())
		p.registers.SetProgramCounter((instruction & 0x0FFF) - 2)
	case '3', '4', '5', '9', 'E':
		p.executeSkipInstructions(instructionCharacter, instruction)
	case '6':
		// LD, load immediate value to register.
		log.Print("INSTRUCTION: Write an immediate value to a register.\n")
		p.registers.WriteRegister(register, immediate)
	case '7':
		// ADD, add immediate value to register.
		log.Print("INSTRUCTION: Add an immediate value to a register.\n")
		p.registers.AddRegisterImmediate(register, immediate)
	case '8':
		p.executeLogicalInstructions(register, register2, LogicalInstructionType(instruction&0xF), instruction)
	case 'A':
		// LD: Load to I register.
		log.Print("INSTRUCTION: Write to an I register.\n")
		p.registers.WriteIRegister(instruction & 0xFFF)
	case 'B':
		// JP: Jump to V0 + NNN.
		log.Print("INSTRUCTION: Jump to a location V0 + immediate.\n")
		p.registers.SetProgramCounter((uint16(p.registers.ReadRegister(0)) + instruction&0xFFF) - 2)
	case 'C':
		// RND: Set VX tp Random byte AND immediate
		p.executeRandomAnd(register, immediate)
	case 'D':
		// DRW draw a sprite to the screen.
		p.executeDrawInstruction(register, register2, byte(instruction&0xF))
	case 'F':
		// Register instructions
		p.executeRegisterInstructions(register, instruction)
	}

}
