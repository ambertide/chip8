package emulator

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const maxBitmask = uint64(1) << 63

type Graphics struct {
	screen      *[32]uint64
	pixelSprite *pixel.Sprite
	window      *pixelgl.Window
	batch       *pixel.Batch
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
		for bitmask := uint64(1); bitmask <= maxBitmask; bitmask = bitmask << 1 {
			x -= 1
			if row&bitmask != 0 { // This means pixel is lit
				// Append the location of the pixel to the matrices as a matrix.
				matrices = append(matrices, pixel.IM.Moved(pixel.V(float64(x), float64(y)).Scaled(10)))
			}
			if bitmask == maxBitmask {
				break
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

func NewGraphics(screenBuffer *[32]uint64) *Graphics {
	graphics := new(Graphics)
	graphics.screen = screenBuffer
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

func (g *Graphics) Mainloop() {
	for !g.window.Closed() {
		g.batch.Clear()
		g.drawPixels()
		g.batch.Draw(g.window)
		g.window.Update()
	}
}

func RunGraphics(screenBuffer *[32]uint64) {
	graphics := NewGraphics(screenBuffer)
	graphics.Mainloop()
}
