package emulator

import (
	"log"
	"os"
	"time"

	"github.com/ambertide/chip8/emulator/device"
)

type Emulator struct {
	screenBuffer   [32]uint64
	processor      *device.Processor
	keyboardBuffer uint16
	clockSpeed     uint64
	programPath    string
}

func NewEmulator(clockSpeed uint64, programPath string) *Emulator {
	emulator := new(Emulator)
	emulator.clockSpeed = clockSpeed
	emulator.programPath = programPath
	emulator.processor = device.NewProcessor(&emulator.screenBuffer, &emulator.keyboardBuffer)
	return emulator
}

func (e *Emulator) RunEmulator(program []byte, programSize uint16) {
	e.processor.LoadProgram(program, programSize)
	for !e.processor.ShouldHalt() {
		e.processor.Cycle()
		time.Sleep(time.Second / time.Duration(e.clockSpeed))
	}
}

func (emulator *Emulator) emulatorCode() {
	romPath := emulator.programPath
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

// Run the emulator subroutines.
func RunEmulator(clockSpeed uint64, programPath string) {
	e := NewEmulator(clockSpeed, programPath)
	log.Println("Emulator initialised.")
	go e.emulatorCode()
	log.Println("Emulator goroutine dispatched.")
	RunGraphics(&e.screenBuffer, &e.keyboardBuffer)
}