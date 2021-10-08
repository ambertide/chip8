package device

type chip8Keyboard struct {
	keyboardMask uint16
}

// Returns true if a key is pressed.
func (k *chip8Keyboard) IsKeyPressed(keyValue byte) bool {
	return k.keyboardMask>>uint16(keyValue)&1 == 1
}

func (k *chip8Keyboard) WaitForKeyPress() byte {
	return 'A'
}
