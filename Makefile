build/chip8:
	cd chip8 &&\
	go build -o ../build/chip8 &&\
	cd ..

all:
	cd chip8 &&\
	go build -o ../build/chip8 &&\
	cd ..

.PHONY: all
	