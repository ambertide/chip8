package main

import (
	"flag"
	"os"

	"github.com/ambertide/chip8/pkg/emulator"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	clockSpeed := flag.Uint64("speed", 500, "Sets the speed of the main processor in Hz.")
	programPath := flag.String("rom", "", "Path to the rom file for chip8.")
	flag.Parse()
	if *programPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	pixelgl.Run(func() { emulator.RunEmulator(*clockSpeed, *programPath) })
}
