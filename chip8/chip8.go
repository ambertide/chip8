package main

import (
	"log"
	"os"
	"time"

	"github.com/ambertide/chip8/device"
	"github.com/ambertide/chip8/emulator"
	"github.com/faiface/pixel/pixelgl"
)

type Emulator struct {
	screenBuffer   [32]uint64
	processor      *device.Processor
	keyboardBuffer uint16
}

func NewEmulator() *Emulator {
	emulator := new(Emulator)
	emulator.processor = device.NewProcessor(&emulator.screenBuffer, &emulator.keyboardBuffer)
	return emulator
}

func (e *Emulator) RunEmulator(program []byte, programSize uint16) {
	e.processor.LoadProgram(program, programSize)
	for !e.processor.ShouldHalt() {
		e.processor.Cycle()
		time.Sleep(time.Second / 500)
	}
}

func emulatorCode(emulator *Emulator) {
	romPath := os.Args[1]
	file, err := os.Open(romPath)
	if err != nil {
		panic(err)
	}
	var program [0xFFF]byte // Maximum ROM size.
	programSize, err2 := file.Read(program[:])
	if err2 != nil {
		panic(err)
	}
	emulator.RunEmulator(program[:], uint16(programSize))
}

func runGraphics() {
	e := NewEmulator()
	log.Println("Emulator initialised.")
	go emulatorCode(e)
	log.Println("Emulator goroutine dispatched.")
	emulator.RunGraphics(&e.screenBuffer, &e.keyboardBuffer)
}

func main() {
	pixelgl.Run(runGraphics)
}
