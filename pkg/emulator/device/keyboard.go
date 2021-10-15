package device

type chip8Keyboard struct {
	keyboardMask *uint16
}

// Returns true if a key is pressed.
func (k *chip8Keyboard) IsKeyPressed(keyValue byte) bool {
	keymask := uint16(1) << keyValue // Calculate mask from value.
	mask := *k.keyboardMask
	return mask&keymask != 0
}

// Decode which key is pressed, emulator
// Does not support multiple key presses
// and the rightmost will be selected.
func DecodeKey(keyboardMask uint16) byte {
	for i := 0; i < 16; i++ {
		keyboardMask := keyboardMask >> i
		if keyboardMask&1 == 1 {
			////log.Printf("User pressed %X\n", i)
			return byte(i)
		}
	}
	return 0
}

// Halt program execution until key is pressed.
func (k *chip8Keyboard) WaitForKeyPress() byte {
	////log.Println("Waiting for key press.")
	for {
		keyboradMask := *k.keyboardMask
		if keyboradMask != 0 {
			return DecodeKey(keyboradMask)
		}
	}
}

// Initialise a new keyboard whose buffer is shared
// Between emulator logic and device logic.
func NewKeyboard(keyboardBuffer *uint16) *chip8Keyboard {
	keyboard := new(chip8Keyboard)
	keyboard.keyboardMask = keyboardBuffer
	return keyboard
}
