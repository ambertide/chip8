# A Chip8 Emulator in Go

chip8.go is a simple Chip8 emulator, compliant with the technical standard laid out in the 
[Cowgod's Manual](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM). Graphics and sound are both
powered by Faiface's [Beep](https://github.com/faiface/beep/) and [Pixel](https://github.com/faiface/pixel) libraries.

## Installation

For Linux distrobutions, you will need `libasound2-dev` package. Your go version should be 17+

```
git clone https://github.com/ambertide/chip8
cd chip8
make all
```

The built file can be accessed in the `build/chip8` directory. If you have added the `chip8` to the path, simply

```
chip8 -rom myrom.ch8
```

to play a rom.

You can also specify the speed using `-speed` flag, by default, the speed is 500MHz
