package device

type pixel bool

type chip8Display struct {
	screen [64][32]pixel
}

func (d *chip8Display) ClearDisplay() {

}
