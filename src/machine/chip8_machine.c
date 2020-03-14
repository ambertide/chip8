//
// Created by egeem on 13/03/2020.
//

#include "chip8_machine.h"
#include "../grammar/chip8_instructions.h"

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

void chip8_machine_init(chip8_machine* machine) {
    chip8_stack_init(&machine->return_stack);
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

}

void parse_return(chip8_machine* machine) {
    uint8_t return_value = chip8_stack_pop(&machine->return_stack);
    chip8_set_pc(machine, return_value);
}

void parse_jump(chip8_machine* machine, chip8_instruction* instruction) {
    chip8_set_pc(machine, parse_hex(1, 3, instruction));
}

void parse_call(chip8_machine* machine, chip8_instruction* instruction) {
    chip8_stack_push(&machine->return_stack, (machine->program_counter - (uint16_t*) &machine->memory));
    chip8_set_pc(machine, parse_hex(1, 3, instruction));
}

void parse_skip_if_val_equal(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] == parse_hex(2, 3, instruction)) {
        chip8_skip_instruction()
    }
}

void parse_skip_if_val_not_equal(chip8_machine* machine, chip8_instruction* instruction) {
    if (machine->registers[instruction->characters[1]] != parse_hex(2, 3, instruction)) {
        chip8_skip_instruction()
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
        case LD: break;
        case ADD: break;
        case LD_A: break;
        case OR: break;
        case AND: break;
        case XOR: break;
        case ADD_A: break;
        case SUB: break;
        case SHR: break;
        case SUBN: break;
        case SHL: break;
        case SNE_A: parse_skip_if_reg_not_equal(machine, &instruction); break;
        case LD_B: break;
        case JP_A: break;
        case RND: break;
        case DRW: break;
        case SKP: parse_skip_if_button_pressed(machine, &instruction); break;
        case SKNP: parse_skip_if_button_not_pressed(machine, &instruction); break;
        case LD_C: break;
        case LD_D: break;
        case LD_E: break;
        case LD_F: break;
        case ADD_B: break;
        case LD_G: break;
        case LD_H: break;
        case LD_I: break;
        case LD_J: break;
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