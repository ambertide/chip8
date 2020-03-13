//
// Created by egeem on 13/03/2020.
//

#include "chip8_machine.h"
#include "../grammar/chip8_instructions.h"

void chip8_stack_init(chip8_stack* stack) {
    stack->stack_pointer = malloc(64 * sizeof(uint8_t));
    stack->stack_bottom = stack->stack_pointer;
}

uint8_t chip8_stack_pop(chip8_stack* stack) {
    uint8_t return_code = *stack->stack_pointer;
    stack->stack_pointer--;
    return return_code;
}

void chip8_stack_push(chip8_stack* stack, uint8_t return_code) {
    stack->stack_pointer++;
    *stack->stack_pointer = return_code;
}

void chip8_machine_init(chip8_machine* machine) {
    chip8_stack_init(&machine->return_stack);
    machine->index_register = 0;
    machine->program_counter = (uint16_t *) &machine->memory;
}


uint16_t chip8_get_instruction(chip8_machine* machine) {
    machine->index_register = *machine->program_counter;
    uint16_t instruction = machine->index_register;
    return instruction;
}

uint16_t parse_instruction(chip8_machine* machine) {
    uint16_t instruction = chip8_get_instruction(machine);
    chip8_instruction instruction_;
    chip8_instruction_init(&instruction_, instruction);
    switch (scan_instruction(&instruction_)) {
        case SYS: break;
        case CLS: break;
        case RET: break;
        case JP: break;
        case CALL: break;
        case SE: break;
        case SNE: break;
        case SE_A: break;
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
        case SNE_A: break;
        case LD_B: break;
        case JP_A: break;
        case RND: break;
        case DRW: break;
        case SKP: break;
        case SKNP: break;
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
