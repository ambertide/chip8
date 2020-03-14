//
// Created by egeem on 13/03/2020.
//

#ifndef CHIP8_CHIP8_MACHINE_H
#define CHIP8_CHIP8_MACHINE_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct {
    uint8_t* stack_bottom;
    uint8_t* stack_pointer;
} chip8_stack;

typedef struct {
    int screen[64][32];
} chip8_frame_buffer;

typedef struct {
    uint8_t registers[16];
    uint8_t memory[4096];
    chip8_stack return_stack;
    chip8_frame_buffer frame_buffer;
    uint8_t sound_timer;
    uint8_t delay_timer;
    uint16_t* program_counter;
    uint16_t index_register;

} chip8_machine;

uint16_t get_char();
#endif //CHIP8_CHIP8_MACHINE_H
