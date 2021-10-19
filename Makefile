BIN = chip8
ifeq ($(OS),Windows_NT)
BIN := $(BIN).exe
endif

build/$(BIN):
	cd cmd/chip8 &&\
	go build -o ../../build/$(BIN)
	cd ../..

all:
	cd cmd/chip8 &&\
	go build -o ../../build/$(BIN)
	cd ../..

.PHONY: all
	
