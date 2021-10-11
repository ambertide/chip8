// Contains structs and methods for the chip-8 memory.
package device

const RamStartLocation = 0x200

// Givena a chip-8 adress, calculate the location
// of the corresponding value in the ram array.
func calculateRAMOffset(address uint16) uint16 {
	return address - RamStartLocation
}

// Check if an address is in the reserved range.
func isAddressInReservedRange(address uint16) bool {
	return address < RamStartLocation
}

type chip8Memory struct {
	// Reserved for the Interpreter.
	reserved [512]byte
	// Program space
	ram [3584]byte
}

// Read a single cell from memory.
func (m *chip8Memory) ReadMemory(address uint16) byte {
	if isAddressInReservedRange(address) {
		return m.reserved[address]
	} else {
		return m.ram[calculateRAMOffset(address)]
	}
}

// Write a single cell of memory.
func (m *chip8Memory) WriteMemory(address uint16, value byte) bool {
	if isAddressInReservedRange(address) {
		return false
	} else {
		m.ram[calculateRAMOffset(address)] = value
		return true
	}
}

// Read from a block of memory between the start and stop addresses.
func (m *chip8Memory) BlockWriteToMemory(start uint16, stop uint16, data []byte) bool {
	if start < 0x200 {
		// Reserved memory cannot be manipulated.
		return false
	}
	// Calculate the destination slice.
	dst := m.ram[calculateRAMOffset(start):calculateRAMOffset(stop)]
	// And copy the data.
	copy(dst, data)
	return true
}

// Read a block of memory between the start and stop adresses.
func (m *chip8Memory) BlockReadFromMemory(start uint16, stop uint16) []byte {
	// Create a buffer to hold copied values
	// Preallocate it as well.
	buffer := make([]byte, stop-start+1)
	// Since the read request may cross the reserved boundry.
	// We must calculate the ranges we will copy from them.
	var reservedStart, reservedStop, ramStart, ramStop uint16
	if isAddressInReservedRange(start) {
		reservedStart = start
	} else {
		ramStart = calculateRAMOffset(start)
	}
	if isAddressInReservedRange(stop) {
		reservedStop = stop
	} else {
		ramStop = calculateRAMOffset(stop)
	}
	// Finally we must calculate the stop point of the reserved
	// Copy in the buffer.
	seperator := reservedStop - reservedStart
	// Now we can finally copy our values.
	copy(buffer[:seperator], m.reserved[reservedStart:reservedStop])
	copy(buffer[seperator:], m.ram[ramStart:ramStop])
	return buffer
}

// Load a program to the chip8 memory given the size and the data
// of the program. isETI660 is used to determine the loading start
// location.
func (m *chip8Memory) loadProgram(program []byte, programSize uint16, isETI660 bool) {
	var startLocation uint16 = RamStartLocation // The RAM start location
	if isETI660 {
		// ETI600 programs start in another location.
		startLocation += 0x400
	}
	// Write the program to memory.
	m.BlockWriteToMemory(startLocation, startLocation+programSize, program)
}

// Load a traditional Chip-8 program to the memory
// Given its data and program size.
func (m *chip8Memory) LoadProgram(program []byte, programSize uint16) {
	m.loadProgram(program, programSize, false)
}

// Load an ETI600 Chip-8 program to the memory
// Given its data and program size.
func (m *chip8Memory) LoadETIProgram(program []byte, programSize uint16) {
	m.loadProgram(program, programSize, true)
}

// Load the reserved sections of the memory
// With apporpirate data.
func (m *chip8Memory) LoadReserved() {
	characterSprites := []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	copy(m.reserved[0:81], characterSprites)
}

func newMemory() *chip8Memory {
	memory := new(chip8Memory)
	memory.LoadReserved()
	return memory
}
