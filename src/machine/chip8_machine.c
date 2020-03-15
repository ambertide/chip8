//
// Created by egeem on 13/03/2020.
//

#include "chip8_machine.h"
#include "../grammar/chip8_instructions.h"
#include <math.h>

void chip8_stack_init(chip8_stack* stack) {
    stack->stack_pointer = malloc(64 * sizeof(uint16_t));
    stack->stack_bottom = stack->stack_pointer;
}

uint16_t chip8_stack_pop(chip8_stack* stack) {
    uint16_t return_code = *stack->stack_pointer;
    stack->stack_pointer--;
    return return_code;
}

void chip8_stack_push(chip8_stack* stack, uint16_t return_code) {
    stack->stack_pointer++;
    *stack->stack_pointer = return_code;
}

void chip8_frame_buffer_init(chip8_frame_buffer* buf) {
    for (int i = 0; i < 64; i++) {
        for (int j = 0; j < 32; j++) {
            buf->screen[i][j] = 0;
        }
    }
}

void chip8_machine_init(chip8_machine* machine) {
    chip8_stack_init(&machine->return_stack);
    chip8_frame_buffer_init(&machine->frame_buffer);
    machine->index_register = 0;
    machine->program_counter = (uint16_t *) &machine->memory;

}


uint16_t chip8_get_instruction(chip8_machine* machine) {
    uint16_t instruction = *machine->program_counter;
    return instruction;
}

void chip8_set_pc(chip8_machine* machine, uint16_t address) {
    machine->program_counter = machine->program_counter = (uint16_t *) (&machine->memory + address);
}

void chip8_skip_instruction(chip8_machine* machine) {
    machine->program_counter+=2;
}

void parse_cls(chip8_machine* machine, chip8_instruction* instruction) {
    chip8_frame_buffer_init(&machine->frame_buffer);
}

void parse_return(chip8_machine* machine) {
    uint8_t return_value = chip8_stack_pop(&machine->return_stack);
    chip8_set_pc(machine, return_value);
}

void parse_jump(chip8_machine* machine, chip8_instruction* instruction) {
    chip8_set_pc(machine, parse_hex(1, 3, instruction));
}

void parse_jump_offset(chip8_machine* machine, chip8_instruction* instruction) {
    chip8_set_pc(machine, parse_hex(1, 3, instruction) + machine->registers[0]);
}

void parse_call(chip8_machine* machine, chip8_instruction* instruction) {
    chip8_stack_push(&machine->return_stack, (machine->program_counter - (uint16_t*) &machine->memory));
    chip8_set_pc(machine, parse_hex(1, 3, instruction));
}

void parse_skip_if_val_equal(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] == parse_hex(2, 3, instruction)) {
        chip8_skip_instruction(machine);
    }
}

void parse_skip_if_val_not_equal(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] != parse_hex(2, 3, instruction)) {
        chip8_skip_instruction(machine);
    }
}

void parse_skip_if_reg_equal(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] == machine->registers[instruction->characters[2]]) {
        chip8_skip_instruction(machine);
    }
}

void parse_skip_if_reg_not_equal(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] != machine->registers[instruction->characters[2]]) {
        chip8_skip_instruction(machine);
    }
}

void parse_skip_if_button_pressed(chip8_machine* machine, chip8_instruction* instruction) {
    if (check_char(machine->registers[instruction->characters[1]], (char) getchar())) {
        chip8_skip_instruction(machine);
    }
}

void parse_skip_if_button_not_pressed(chip8_machine* machine, chip8_instruction* instruction) {
    if (!check_char(machine->registers[instruction->characters[1]], (char) getchar())) {
        chip8_skip_instruction(machine);
    }
}

void parse_load_reg_val(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] = parse_hex(2, 3, instruction);
}

void parse_load_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] = machine->registers[instruction->characters[2]];
}

void parse_load_index_register(chip8_machine* machine, chip8_instruction* instruction) {
    machine->index_register = parse_hex(1,3,instruction);
}

void parse_load_from_delay_timer(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] = machine->delay_timer;
}

void parse_load_keypress(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] = get_char();
}

void parse_load_to_delay_timer(chip8_machine* machine, chip8_instruction* instruction){
    machine->delay_timer = machine->registers[instruction->characters[1]];
}

void parse_load_to_sound_timer(chip8_machine* machine, chip8_instruction* instruction){
    machine->sound_timer = machine->registers[instruction->characters[1]];
}

void parse_add_register_byte(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] += parse_hex(2, 3, instruction);
}

void parse_or_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] |= machine->registers[instruction->characters[2]];
}

void parse_and_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] &= machine->registers[instruction->characters[2]];
}

void parse_xor_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[instruction->characters[1]] ^= machine->registers[instruction->characters[2]];
}

void parse_add_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    if ((uint32_t) machine->registers[instruction->characters[1]] + (uint32_t) machine->registers[instruction->characters[2]] > 255) {
        machine->registers[15] = 1; // Set VF to 1 if there is overflow
    } else {
        machine->registers[15] = 0; // Otherwise to 0
    }
    machine->registers[instruction->characters[1]] =
            ((uint32_t) machine->registers[instruction->characters[1]]
            + (uint32_t) machine->registers[instruction->characters[2]]) % 255;
}

void parse_sub_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] > machine->registers[instruction->characters[2]]) {
        machine->registers[15] = 1; // Set VF to 1 if x > y
    } else {
        machine->registers[15] = 0; // Otherwise 0
    }
    machine->registers[instruction->characters[1]] -= machine->registers[instruction->characters[2]];
}

void parse_sub_neg_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[2]] > machine->registers[instruction->characters[1]]) {
        machine->registers[15] = 1; // Set VF to 1 if y > x
    } else {
        machine->registers[15] = 0; // Otherwise to 0
    }
    machine->registers[instruction->characters[2]] -= machine->registers[instruction->characters[1]];
}

void parse_shr_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[15] = machine->registers[instruction->characters[1]] % 2; // 1 if LSB is 1 otherwise 0
    machine->registers[instruction->characters[1]] >>= 1;
}

void parse_shl_reg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->registers[15] = machine->registers[instruction->characters[1]] >= 128 ;
    machine->registers[instruction->characters[1]] <<= 1;
}

void parse_add_ireg_reg(chip8_machine* machine, chip8_instruction* instruction) {
    machine->index_register = ((uint32_t) machine->index_register + machine->registers[instruction->characters[1]]) % 65535;
}

void parse_rand(chip8_machine* machine, chip8_instruction* instruction) {
    uint8_t random = rand() % 256;
    machine->registers[instruction->characters[1]] = parse_hex(2, 3, instruction) & random;
}

void parse_load_from_reg_mem(chip8_machine* machine, chip8_instruction* instruction) {
    int mem_location = machine->index_register;
    int reg_location = instruction->characters[1];
    for (int i = 0; i <= reg_location; i++) {
        machine->memory[mem_location + i] = machine->registers[i];
    }
}

void parse_load_from_mem_reg(chip8_machine* machine, chip8_instruction* instruction) {
    int mem_location = machine->index_register;
    int reg_location = instruction->characters[1];
    for (int i = 0; i <= reg_location; i++) {
        machine->registers[i] = machine->memory[mem_location + i];
    }
}

void parse_load_bcd(chip8_machine* machine, chip8_instruction* instruction) {
    uint8_t number = machine->registers[1];
    for (int i = machine->index_register; i < machine->index_register + 3; i++) {
        machine->memory[i] = (number / (uint8_t) pow(10, 3 - i)) % 10;
    }
}

void parse_instruction(chip8_machine* machine) {
    uint16_t instruction_val = chip8_get_instruction(machine);
    chip8_instruction instruction;
    chip8_instruction_init(&instruction, instruction_val);
    switch (scan_instruction(&instruction)) {
        case SYS: break;
        case CLS: parse_cls(machine, &instruction); break;
        case RET: parse_return(machine); break;
        case JP: parse_jump(machine, &instruction); break;
        case CALL: parse_call(machine, &instruction); break;
        case SE: parse_skip_if_val_equal(machine, &instruction); break;
        case SNE: parse_skip_if_val_not_equal(machine, &instruction); break;
        case SE_A: parse_skip_if_reg_equal(machine, &instruction); break;
        case LD: parse_load_reg_val(machine, &instruction); break;
        case ADD: parse_add_register_byte(machine, &instruction); break;
        case LD_A: parse_load_reg_reg(machine, &instruction); break;
        case OR: parse_or_reg_reg(machine, &instruction); break;
        case AND: parse_and_reg_reg(machine, &instruction); break;
        case XOR: parse_xor_reg_reg(machine, &instruction); break;
        case ADD_A: parse_add_reg_reg(machine, &instruction); break;
        case SUB: parse_sub_reg_reg(machine, &instruction); break;
        case SHR: parse_shr_reg_reg(machine, &instruction); break;
        case SUBN: parse_sub_neg_reg_reg(machine, &instruction); break;
        case SHL: parse_shl_reg_reg(machine, &instruction); break;
        case SNE_A: parse_skip_if_reg_not_equal(machine, &instruction); break;
        case LD_B: parse_load_index_register(machine, &instruction); break;
        case JP_A: parse_jump_offset(machine, &instruction); break;
        case RND: parse_rand(machine, &instruction); break;
        case DRW: break;
        case SKP: parse_skip_if_button_pressed(machine, &instruction); break;
        case SKNP: parse_skip_if_button_not_pressed(machine, &instruction); break;
        case LD_C: parse_load_from_delay_timer(machine, &instruction); break;
        case LD_D: parse_load_keypress(machine, &instruction); break;
        case LD_E: parse_load_to_delay_timer(machine, &instruction); break;
        case LD_F: parse_load_to_sound_timer(machine, &instruction); break;
        case ADD_B: parse_add_ireg_reg(machine, &instruction); break;
        case LD_G: break;
        case LD_H: parse_load_bcd(machine, &instruction); break;
        case LD_I: parse_load_from_reg_mem(machine, &instruction); break;
        case LD_J: parse_load_from_mem_reg(machine, &instruction); break;
    }
}

void execute_timers(chip8_machine* machine) {
    if (machine->delay_timer>0) {
        machine->delay_timer--;
    }

    if (machine->sound_timer>0) {
        machine->sound_timer--;
    }
}

void run_next(chip8_machine* machine) {
    parse_instruction(machine);
    execute_timers(machine);
}