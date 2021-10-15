package device

type chip8Stack struct {
	// Although technically a register
	// It fits here much better.
	stackPointer uint16
	// This is where the addresses are hold.
	addresses [16]uint16
}

// Pop an address from the stack.
func (s *chip8Stack) Pop() uint16 {
	s.stackPointer--
	stackValue := s.addresses[s.stackPointer]
	return stackValue
}

// Push an address to the stack.
func (s *chip8Stack) Push(address uint16) {
	s.addresses[s.stackPointer] = address
	s.stackPointer++
}
