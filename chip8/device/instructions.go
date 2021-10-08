package device

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
	case "00E0":
		// CLR: Clear Screen
		p.display.ClearDisplay()
	case "00EE":
		// RET: Return from subroutine.
		p.registers.SetProgramCounter(p.stack.Pop())
	}
}

// Execute skip instructions that skip a number of intsructions
// Depending on a condition. SE, SNE, SE.
func (p *Processor) executeSkipInstructions(instructionChar string, instruction uint16) {
	byteValue := byte(instruction & 0xFF) // For 3XNN and 4XNN
	register := getRegisterIndex(instructionChar[1])
	register2 := getRegisterIndex(instructionChar[2]) // For 5XY0
	msb := instructionChar[0]                         // Most significant byte
	switch {
	case msb == '3' && p.registers.ReadRegister(register) == byteValue:
		// Incrementing the program counter now will effectively
		// Skip this instruction.
		p.registers.IncrementProgramCounter()
	case msb == '4' && p.registers.ReadRegister(register) != byteValue:
		p.registers.IncrementProgramCounter()
	case msb == '5' && instructionChar[3] == '0' && p.registers.CompareRegisters(register, register2):
		p.registers.IncrementProgramCounter()
	}
}

// Given the indexes for two registers and the operation type execute
// the operation.
func (p *Processor) executeLogicalInstructions(x uint8, y uint8, operationType LogicalInstructionType) {
	switch operationType {
	case Load:
		p.registers.RegisterOperation(x, y, func(b1 byte, b2 byte) byte { return b2 })
	case Or:
		p.registers.RegisterOperation(x, y, func(b1 byte, b2 byte) byte { return b1 | b2 })
	case And:
		p.registers.RegisterOperation(x, y, func(b1, b2 byte) byte { return b1 & b2 })
	case Xor:
		p.registers.RegisterOperation(x, y, func(b1, b2 byte) byte { return b1 ^ b2 })
	case Add:
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 + b2 },
			func(b1, b2 byte) byte {
				return byte(uint16(b1)+uint16(b2)>>16) & 0x1 // Calculate the carry bit.
			},
		)
	case Sub:
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
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 >> 1 },
			func(b1, b2 byte) byte {
				return b1 & 0x1
			},
		)
	case Subn:
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
		p.registers.RegisterOperationWithCarry(x, y, func(b1, b2 byte) byte { return b1 << 1 },
			func(b1, b2 byte) byte {
				return b1 >> 7
			},
		)
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
		p.registers.SetProgramCounter(instruction & 0x0FFF)
	case '2':
		// CALL: Call a subroutine.
		p.stack.Push(p.registers.GetProgramCounter())
		p.registers.SetProgramCounter(instruction & 0x0FFF)
	case '3', '4', '5':
		p.executeSkipInstructions(instructionCharacter, instruction)
	case '6':
		// LD, load immediate value to register.
		p.registers.WriteRegister(register, immediate)
	case '7':
		// ADD, add immediate value to register.
		p.registers.AddRegisterImmediate(register, immediate)
	case '8':
		p.executeLogicalInstructions(register, register2, LogicalInstructionType(instruction&0xF))
	}
}
