package emulator

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Graphics struct {
	screen         *[32]uint64
	pixelSprite    *pixel.Sprite
	window         *pixelgl.Window
	batch          *pixel.Batch
	keyboardBuffer *uint16
}

var keysToChip8 = map[pixelgl.Button]uint16{
	pixelgl.Key0: 1,
	pixelgl.Key1: 2,
	pixelgl.Key2: 4,
	pixelgl.Key3: 8,
	pixelgl.Key4: 16,
	pixelgl.Key5: 32,
	pixelgl.Key6: 64,
	pixelgl.Key7: 128,
	pixelgl.Key8: 256,
	pixelgl.Key9: 512,
	pixelgl.KeyA: 1024,
	pixelgl.KeyB: 2048,
	pixelgl.KeyC: 4096,
	pixelgl.KeyD: 8192,
	pixelgl.KeyE: 16384,
	pixelgl.KeyF: 32768,
}

// Almost verbatim from the Pixel tutorial in
// https://github.com/faiface/pixel/wiki/Drawing-a-Sprite
func (g *Graphics) loadPixelSprite() pixel.Picture {
	file, err := os.Open("assets/pixelSprite.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, _, _ := image.Decode(file)
	picture := pixel.PictureDataFromImage(img)
	g.pixelSprite = pixel.NewSprite(picture, picture.Bounds())
	return picture
}

// Calculate the matrices to locate sprites.
func (g *Graphics) calculateMatrices() []pixel.Matrix {
	matrices := []pixel.Matrix{}
	for y, row := range g.screen {
		x := 64 // We are coming from the reverse side.
		for bitmask := uint64(1); bitmask != 0; bitmask = bitmask << 1 {
			x -= 1
			if row&bitmask != 0 { // This means pixel is lit
				// Append the location of the pixel to the matrices as a matrix.
				matrices = append(matrices, pixel.IM.Moved(pixel.V(float64(x), float64(32-y)).Scaled(10)))
			}
		}
	}
	return matrices
}

// Draw the pixels to the screen.
func (g *Graphics) drawPixels() {
	pixelLocations := g.calculateMatrices()
	for _, pixelLocation := range pixelLocations {
		g.pixelSprite.Draw(g.batch, pixelLocation)
	}

}

func NewGraphics(screenBuffer *[32]uint64, keyboardBuffer *uint16) *Graphics {
	graphics := new(Graphics)
	graphics.screen = screenBuffer
	graphics.keyboardBuffer = keyboardBuffer
	var err error
	config := pixelgl.WindowConfig{
		Title:  "Chip8",
		Bounds: pixel.R(0, 0, 640, 320),
		VSync:  true,
	}
	graphics.window, err = pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}
	pixelPicture := graphics.loadPixelSprite()
	graphics.batch = pixel.NewBatch(&pixel.TrianglesData{}, pixelPicture)
	return graphics
}

// Reset and recalculate the keyboard buffer
// From a slice of pressed keys.
func (g *Graphics) updateKeyboardBuffer(keyValues []uint16) {
	newMask := uint16(0x0)
	for _, keyValue := range keyValues {
		newMask ^= keyValue
	}
	*g.keyboardBuffer = newMask
}

// Handle keyboard presses by the user.
func (g *Graphics) handleKeyboard() {
	pressedKeys := []uint16{}
	for key, value := range keysToChip8 {
		if g.window.Pressed(key) {
			pressedKeys = append(pressedKeys, value)
		}
	}
	g.updateKeyboardBuffer(pressedKeys)
}

// Loop through the graphics engine.
func (g *Graphics) Mainloop() {
	for !g.window.Closed() {
		g.handleKeyboard()
		g.batch.Clear()
		g.drawPixels()
		g.batch.Draw(g.window)
		g.window.Update()
	}
}

func RunGraphics(screenBuffer *[32]uint64, keyboardBuffer *uint16) {
	log.Println("Graphic initialisation starting...")
	graphics := NewGraphics(screenBuffer, keyboardBuffer)
	log.Println("Graphics initialised")
	graphics.Mainloop()

}
