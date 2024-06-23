package emulator 

type Operand struct {
    mode AddressingMode
    address uint16
    extra_cycle bool
}

type Instruction struct {
    Code byte
    Name string
    Size uint8
    cycles uint8
    Mode AddressingMode
    handler func (*CPU, Operand)
}

// It would be ideal if we could store this in such a way that the index of the code IS the code.
var Instructions = map[byte]Instruction {
    0x69: {0x69, "ADC", 2, 2, Immediate, adc},
    0x65: {0x65, "ADC", 2, 3, ZeroPage,  adc},
    0x75: {0x75, "ADC", 2, 4, ZeroPageX, adc},
    0x6D: {0x6D, "ADC", 3, 4, Absolute,  adc},
    0x7D: {0x7D, "ADC", 3, 4, AbsoluteX, adc},
    0x79: {0x79, "ADC", 3, 4, AbsoluteY, adc},
    0x61: {0x61, "ADC", 2, 6, IndirectX, adc},
    0x71: {0x71, "ADC", 2, 5, IndirectY, adc},

    0x29: {0x29, "AND", 2, 2, Immediate, and},
    0x25: {0x25, "AND", 2, 3, ZeroPage,  and},
    0x35: {0x35, "AND", 2, 4, ZeroPageX, and},
    0x2D: {0x2D, "AND", 3, 4, Absolute,  and},
    0x3D: {0x3D, "AND", 3, 4, AbsoluteX, and},
    0x39: {0x39, "AND", 3, 4, AbsoluteY, and},
    0x21: {0x21, "AND", 2, 6, IndirectX, and},
    0x31: {0x31, "AND", 2, 5, IndirectY, and},

    0x0A: {0x0A, "ASL", 1, 2, Accumulator, asl},
    0x06: {0x06, "ASL", 2, 5, ZeroPage,    asl},
    0x16: {0x16, "ASL", 2, 6, ZeroPageX,   asl},
    0x0E: {0x0E, "ASL", 3, 6, Absolute,    asl},
    0x1E: {0x1E, "ASL", 3, 7, AbsoluteX,   asl},

    0x90: {0x90, "BCC", 2, 2, Relative, bcc},
    0xB0: {0xB0, "BCS", 2, 2, Relative, bcs},
    0xF0: {0xF0, "BEQ", 2, 2, Relative, beq},
    
    0x24: {0x24, "BIT", 2, 3, ZeroPage, bit},
    0x2C: {0x2C, "BIT", 3, 4, Absolute, bit},
    
    0x30: {0x30, "BMI", 2, 2, Relative, bmi},
    0xD0: {0xD0, "BNE", 2, 2, Relative, bne},
    0x10: {0x10, "BPL", 2, 2, Relative, bpl},
    
    0x00: {0x00, "BRK", 1, 7, Implied, brk},
    
    0x50: {0x50, "BVC", 2, 2, Relative, bvc},
    0x70: {0x70, "BVS", 2, 2, Relative, bvs},
    
    0x18: {0x18, "CLC", 1, 2, Implied, clc},
    0xD8: {0xD8, "CLD", 1, 2, Implied, cld},
    0x58: {0x58, "CLI", 1, 2, Implied, cli},
    0xB8: {0xB8, "CLV", 1, 2, Implied, clv},
    
    0xC9: {0xC9, "CMP", 2, 2, Immediate, cmp},
    0xC5: {0xC5, "CMP", 2, 3, ZeroPage, cmp},
    0xD5: {0xD5, "CMP", 2, 4, ZeroPageX, cmp},
    0xCD: {0xCD, "CMP", 3, 4, Absolute, cmp},
    0xDD: {0xDD, "CMP", 3, 4, AbsoluteX, cmp},
    0xD9: {0xD9, "CMP", 3, 4, AbsoluteY, cmp},
    0xC1: {0xC1, "CMP", 2, 6, IndirectX, cmp},
    0xD1: {0xD1, "CMP", 2, 5, IndirectY, cmp},
    
    0xE0: {0xE0, "CPX", 2, 2, Immediate, cpx},
    0xE4: {0xE4, "CPX", 2, 3, ZeroPage, cpx},
    0xEC: {0xEC, "CPX", 3, 4, Absolute, cpx},
    
    0xC0: {0xC0, "CPY", 2, 2, Immediate, cpy},
    0xC4: {0xC4, "CPY", 2, 3, ZeroPage, cpy},
    0xCC: {0xCC, "CPY", 3, 4, Absolute, cpy},
    
    0xC6: {0xC6, "DEC", 2, 5, ZeroPage, dec},
    0xD6: {0xD6, "DEC", 2, 6, ZeroPageX, dec},
    0xCE: {0xCE, "DEC", 3, 6, Absolute, dec},
    0xDE: {0xDE, "DEC", 3, 7, AbsoluteX, dec},
    
    0xCA: {0xCA, "DEX", 1, 2, Implied, dex},
    0x88: {0x88, "DEY", 1, 2, Implied, dey},
    
    0x49: {0x49, "EOR", 2, 2, Immediate, eor},
    0x45: {0x45, "EOR", 2, 3, ZeroPage, eor},
    0x55: {0x55, "EOR", 2, 4, ZeroPageX, eor},
    0x4D: {0x4D, "EOR", 3, 4, Absolute, eor},
    0x5D: {0x5D, "EOR", 3, 4, AbsoluteX, eor},
    0x59: {0x59, "EOR", 3, 4, AbsoluteY, eor},
    0x41: {0x41, "EOR", 2, 6, IndirectX, eor},
    0x51: {0x51, "EOR", 2, 5, IndirectY, eor},
    
    0xE6: {0xE6, "INC", 2, 5, ZeroPage, inc},
    0xF6: {0xF6, "INC", 2, 6, ZeroPageX, inc},
    0xEE: {0xEE, "INC", 3, 6, Absolute, inc},
    0xFE: {0xFE, "INC", 3, 7, AbsoluteX, inc},
    
    0xE8: {0xE8, "INX", 1, 2, Implied, inx},
    0xC8: {0xC8, "INY", 1, 2, Implied, iny},
    
    0x4C: {0x4C, "JMP", 3, 3, Absolute, jmp},
    0x6C: {0x6C, "JMP", 3, 5, Indirect, jmp},
    
    0x20: {0x20, "JSR", 3, 6, Absolute, jsr},
    
    0xA9: {0xA9, "LDA", 2, 2, Immediate, lda},
    0xA5: {0xA5, "LDA", 2, 3, ZeroPage, lda},
    0xB5: {0xB5, "LDA", 2, 4, ZeroPageX, lda},
    0xAD: {0xAD, "LDA", 3, 4, Absolute, lda},
    0xBD: {0xBD, "LDA", 3, 4, AbsoluteX, lda},
    0xB9: {0xB9, "LDA", 3, 4, AbsoluteY, lda},
    0xA1: {0xA1, "LDA", 2, 6, IndirectX, lda},
    0xB1: {0xB1, "LDA", 2, 5, IndirectY, lda},
    
    0xA2: {0xA2, "LDX", 2, 2, Immediate, ldx},
    0xA6: {0xA6, "LDX", 2, 3, ZeroPage, ldx},
    0xB6: {0xB6, "LDX", 2, 4, ZeroPageY, ldx},
    0xAE: {0xAE, "LDX", 3, 4, Absolute, ldx},
    0xBE: {0xBE, "LDX", 3, 4, AbsoluteY, ldx},
    
    0xA0: {0xA0, "LDY", 2, 2, Immediate, ldy},
    0xA4: {0xA4, "LDY", 2, 3, ZeroPage, ldy},
    0xB4: {0xB4, "LDY", 2, 4, ZeroPageX, ldy},
    0xAC: {0xAC, "LDY", 3, 4, Absolute, ldy},
    0xBC: {0xBC, "LDY", 3, 4, AbsoluteX, ldy},
    
    0x4A: {0x4A, "LSR", 1, 2, Accumulator, lsr},
    0x46: {0x46, "LSR", 2, 5, ZeroPage, lsr},
    0x56: {0x56, "LSR", 2, 6, ZeroPageX, lsr},
    0x4E: {0x4E, "LSR", 3, 6, Absolute, lsr},
    0x5E: {0x5E, "LSR", 3, 7, AbsoluteX, lsr},
    
    0xEA: {0xEA, "NOP", 1, 2, Implied, nop},
    
    0x09: {0x09, "ORA", 2, 2, Immediate, ora},
    0x05: {0x05, "ORA", 2, 3, ZeroPage, ora},
    0x15: {0x15, "ORA", 2, 4, ZeroPageX, ora},
    0x0D: {0x0D, "ORA", 3, 4, Absolute, ora},
    0x1D: {0x1D, "ORA", 3, 4, AbsoluteX, ora},
    0x19: {0x19, "ORA", 3, 4, AbsoluteY, ora},
    0x01: {0x01, "ORA", 2, 6, IndirectX, ora},
    0x11: {0x11, "ORA", 2, 5, IndirectY, ora},
    
    0x48: {0x48, "PHA", 1, 3, Implied, pha},
    0x08: {0x08, "PHP", 1, 3, Implied, php},
    0x68: {0x68, "PLA", 1, 4, Implied, pla},
    0x28: {0x28, "PLP", 1, 4, Implied, plp},
    
    0x2A: {0x2A, "ROL", 1, 2, Accumulator, rol},
    0x26: {0x26, "ROL", 2, 5, ZeroPage, rol},
    0x36: {0x36, "ROL", 2, 6, ZeroPageX, rol},
    0x2E: {0x2E, "ROL", 3, 6, Absolute, rol},
    0x3E: {0x3E, "ROL", 3, 7, AbsoluteX, rol},
    
    0x6A: {0x6A, "ROR", 1, 2, Accumulator, ror},
    0x66: {0x66, "ROR", 2, 5, ZeroPage, ror},
    0x76: {0x76, "ROR", 2, 6, ZeroPageX, ror},
    0x6E: {0x6E, "ROR", 3, 6, Absolute, ror},
    0x7E: {0x7E, "ROR", 3, 7, AbsoluteX, ror},
    
    0x40: {0x40, "RTI", 1, 6, Implied, rti},
    0x60: {0x60, "RTS", 1, 6, Implied, rts},
    
    0xE9: {0xE9, "SBC", 2, 2, Immediate, sbc},
    0xE5: {0xE5, "SBC", 2, 3, ZeroPage, sbc},
    0xF5: {0xF5, "SBC", 2, 4, ZeroPageX, sbc},
    0xED: {0xED, "SBC", 3, 4, Absolute, sbc},
    0xFD: {0xFD, "SBC", 3, 4, AbsoluteX, sbc},
    0xF9: {0xF9, "SBC", 3, 4, AbsoluteY, sbc},
    0xE1: {0xE1, "SBC", 2, 6, IndirectX, sbc},
    0xF1: {0xF1, "SBC", 2, 5, IndirectY, sbc},
    
    0x38: {0x38, "SEC", 1, 2, Implied, sec},
    0xF8: {0xF8, "SED", 1, 2, Implied, sed},
    0x78: {0x78, "SEI", 1, 2, Implied, sei},
    
    0x85: {0x85, "STA", 2, 3, ZeroPage, sta},
    0x95: {0x95, "STA", 2, 4, ZeroPageX, sta},
    0x8D: {0x8D, "STA", 3, 4, Absolute, sta},
    0x9D: {0x9D, "STA", 3, 5, AbsoluteX, sta},
    0x99: {0x99, "STA", 3, 5, AbsoluteY, sta},
    0x81: {0x81, "STA", 2, 6, IndirectX, sta},
    0x91: {0x91, "STA", 2, 6, IndirectY, sta},
    
    0x86: {0x86, "STX", 2, 3, ZeroPage, stx},
    0x96: {0x96, "STX", 2, 4, ZeroPageY, stx},
    0x8E: {0x8E, "STX", 3, 4, Absolute, stx},
    
    0x84: {0x84, "STY", 2, 3, ZeroPage, sty},
    0x94: {0x94, "STY", 2, 4, ZeroPageX, sty},
    0x8C: {0x8C, "STY", 3, 4, Absolute, sty},
    
    0xAA: {0xAA, "TAX", 1, 2, Implied, tax},
    0xA8: {0xA8, "TAY", 1, 2, Implied, tay},
    0xBA: {0xBA, "TSX", 1, 2, Implied, tsx},
    0x8A: {0x8A, "TXA", 1, 2, Implied, txa},
    0x9A: {0x9A, "TXS", 1, 2, Implied, txs},
    0x98: {0x98, "TYA", 1, 2, Implied, tya},
}
