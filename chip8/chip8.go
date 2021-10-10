package main

import (
	"os"
	"time"

	"github.com/ambertide/chip8/device"
	"github.com/ambertide/chip8/emulator"
	"github.com/faiface/pixel/pixelgl"
)

type Emulator struct {
	screenBuffer [32]uint64
	halt         bool
	processor    *device.Processor
}

func NewEmulator() *Emulator {
	emulator := new(Emulator)
	emulator.processor = device.NewProcessor(&emulator.screenBuffer)
	return emulator
}

func (e *Emulator) RunEmulator(program []byte) {
	e.processor.LoadProgram(program, 478)
	for !e.processor.ShouldHalt() {
		e.processor.Cycle()
		time.Sleep(1 / 60 * time.Second)
	}
}

func emulatorCode(emulator *Emulator) {
	file, err := os.Open("../test_opcode.ch8")
	if err != nil {
		panic(err)
	}
	var program [478]byte
	_, err2 := file.Read(program[:])
	if err2 != nil {
		panic(err)
	}
	emulator.RunEmulator(program[:])
}

func runGraphics() {
	e := NewEmulator()
	go emulatorCode(e)
	emulator.RunGraphics(&e.screenBuffer)
}

func main() {
	pixelgl.Run(runGraphics)
}
