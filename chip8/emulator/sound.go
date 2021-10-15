package emulator

import (
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

func BeepRoutine(soundTimer *bool) {
	speaker.Init(44100, 735)
	sound, err := generators.SinTone(44100, 1190)
	if err != nil {
		return
	}
	for {
		if *soundTimer {
			speaker.Play(beep.Take(44100/60, sound))
		}
		time.Sleep(time.Second / 60)
	}
}
