package device

type chip8Display struct {
	// A 64x32 pixel screen
	// Can easilly be represented
	// This way.
	screen [32]uint64
	// Global buffer used to
	// Communicate between Goroutines.
	screenBuffer *[32]uint64
}

func newDisplay(screenBuffer *[32]uint64) *chip8Display {
	display := new(chip8Display)
	display.screenBuffer = screenBuffer
	return display
}

func (d *chip8Display) SyncBuffer() {
	copy((*d.screenBuffer)[:], d.screen[:])
}

// Clear the display.
func (d *chip8Display) ClearDisplay() {
	for i := 0; i < 32; i++ {
		d.screen[i] = 0
	}
	d.SyncBuffer()
}

// Draw a sprite of height height into the screen
// Starting from x and y, return true if any pixels
// are erased.
func (d *chip8Display) DrawSprite(x byte, y byte, height byte, sprite []byte) bool {
	collusion := false
	for i, spriteRow := range sprite[:height] {
		displayRow := d.screen[(i+int(y))%32]
		// Convert sprite row to have space.
		paddedSpriteRow := uint64(spriteRow)
		// Align the sprite row to its XOR location.
		alignedSpriteRow := paddedSpriteRow
		if x < 56 {
			alignedSpriteRow = paddedSpriteRow << (56 - x)
		} else {
			alignedSpriteRow = paddedSpriteRow >> (x - 56)
		}
		// Check for collusion
		collusion = collusion || (alignedSpriteRow&displayRow != 0)
		// And XOR the screen.
		d.screen[(i+int(y))%32] ^= alignedSpriteRow
	}
	d.SyncBuffer()
	return collusion
}
