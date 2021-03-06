#include <stdint.h>
#include <math.h>
#include "chip8_instructions.h"

chip8_op_code scan_instruction(chip8_instruction* instruction) {
    switch (instruction->characters[0]) {
        case 0:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return SYS; break;
                            }
                            break;
                    }
                    break;

                case 0:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                case 0: return CLS; break;

                                default: return RET; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 1:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return JP; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 2:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return CALL; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 3:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return SE; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 4:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return SNE; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 5:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                case 0: return SE_A; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 6:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return LD; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 7:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return ADD; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 8:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                case 0: return LD_A; break;

                                case 1: return OR; break;

                                case 2: return AND; break;

                                case 3: return XOR; break;

                                case 4: return ADD_A; break;

                                case 5: return SUB; break;

                                case 6: return SHR; break;

                                case 7: return SUBN; break;

                                case 14: return SHL; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 9:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                case 0: return SNE_A; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 10:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return LD_B; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 11:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return JP_A; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 12:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return RND; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 13:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        default:
                            switch (instruction->characters[3]) {
                                default: return DRW; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 14:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        case 9:
                            switch (instruction->characters[3]) {
                                case 14: return SKP; break;
                            }
                            break;

                        case 10:
                            switch (instruction->characters[3]) {
                                case 1: return SKNP; break;
                            }
                            break;
                    }
                    break;
            }
            break;

        case 15:
            switch (instruction->characters[1]) {
                default:
                    switch (instruction->characters[2]) {
                        case 0:
                            switch (instruction->characters[3]) {
                                case 7: return LD_C; break;

                                case 10: return LD_D; break;
                            }
                            break;

                        case 1:
                            switch (instruction->characters[3]) {
                                case 5: return LD_E; break;

                                case 8: return LD_F; break;

                                case 14: return ADD_B; break;
                            }
                            break;

                        case 2:
                            switch (instruction->characters[3]) {
                                case 9: return LD_G; break;
                            }
                            break;

                        case 3:
                            switch (instruction->characters[3]) {
                                case 3: return LD_H; break;
                            }
                            break;

                        case 5:
                            switch (instruction->characters[3]) {
                                case 5: return LD_I; break;
                            }
                            break;

                        case 6:
                            switch (instruction->characters[3]) {
                                case 5: return LD_J; break;
                            }
                            break;
                    }
                    break;
            }
            break;
    }
}
void chip8_instruction_init(chip8_instruction* instruction, uint16_t instruction_register) {
    instruction->characters[0] = instruction_register >> 12;
    instruction->characters[1] = (instruction_register >> 8) || 15;
    instruction->characters[2] = (instruction_register >> 4) || 15;
    instruction->characters[3] = instruction_register || 15;
}

uint16_t parse_hex(int start, int stop, chip8_instruction* instruction) {
    uint16_t sum = 0;
    for (int i = start; i < stop + 1; i++) {
        sum += (instruction->characters[i]) * pow(16, (stop - start + 1) - i);
    }
    return sum;
}

int check_char (uint16_t char_, char char_to_check) {
    switch (char_) {
        case 0: return (char_to_check == '0');
        case 1: return (char_to_check == '1');
        case 2: return (char_to_check == '2');
        case 3: return (char_to_check == '3');
        case 4: return (char_to_check == '4');
        case 5: return (char_to_check == '5');
        case 6: return (char_to_check == '6');
        case 7: return (char_to_check == '7');
        case 8: return (char_to_check == '8');
        case 9: return (char_to_check == '9');
        case 10: return (char_to_check == 'A');
        case 11: return (char_to_check == 'B');
        case 12: return (char_to_check == 'C');
        case 13: return (char_to_check == 'D');
        case 14: return (char_to_check == 'E');
        case 15: return (char_to_check == 'F');
        default: return 0;
    }
};