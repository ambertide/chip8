build/chip8:
	cd cmd/chip8 &&\
	go build -o ../../build/chip8
	cd ../..

all:
	cd cmd/chip8 &&\
	go build -o ../../build/chip8
	cd ../..

.PHONY: all
	