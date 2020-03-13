//
// Created by egeem on 13/03/2020.
//

#ifndef CHIP8_CHIP8_INSTRUCTIONS_H
#define CHIP8_CHIP8_INSTRUCTIONS_H
typedef struct {unsigned char characters[4];} chip8_instruction;
typedef enum {
    SYS,
    CLS,
    RET,
    JP,
    CALL,
    SE,
    SNE,
    SE_A,
    LD,
    ADD,
    LD_A,
    OR,
    AND,
    XOR,
    ADD_A,
    SUB,
    SHR,
    SUBN,
    SHL,
    SNE_A,
    LD_B,
    JP_A,
    RND,
    DRW,
    SKP,
    SKNP,
    LD_C,
    LD_D,
    LD_E,
    LD_F,
    ADD_B,
    LD_G,
    LD_H,
    LD_I,
    LD_J
} chip8_op_code;
chip8_op_code scan_instruction(chip8_instruction* instruction);
void chip8_instruction_init(chip8_instruction* instruction, uint16_t register_);

#endif //CHIP8_CHIP8_INSTRUCTIONS_H
