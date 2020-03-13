cmake_minimum_required(VERSION 3.14)
project(chip8 C)

set(CMAKE_C_STANDARD 99)

add_executable(chip8 src/main.c src/grammar/chip8_instructions.c src/grammar/chip8_instructions.h src/machine/chip8_machine.c src/machine/chip8_machine.h)