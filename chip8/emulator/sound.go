package emulator

import (
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

func BeepRoutine(soundTimer *bool) {
	speaker.Init(240, 0)
	sound, err := generators.SinTone(240, 119)
	if err != nil {
		return
	}
	for {
		if *soundTimer {
			speaker.Play(sound)
		}
	}
}
