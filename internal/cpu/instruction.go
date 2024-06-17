package cpu

// Sets the flags for CPU based on the output of an operation (result)
func (this *CPU) SetFlags(result uint8) {
    // Sets the Z flag (ZERO FLAG)
    if result == 0 {
        this.Status = this.Status | FLAG_ZERO;
    } else {
        this.Status = this.Status & ^FLAG_ZERO; // Go does not have a bitwise NOT. So use XOR trick.
    }

    // Sets the N flag (NEGATIVE FLAG)
    if result & 0b1000_0000 != 0 {
        this.Status = this.Status | FLAG_NEGATIVE;
    } else {
        this.Status = this.Status & ^FLAG_NEGATIVE;
    }
}

func (this *CPU) setCFlag(active bool) {
    if active {
        this.Status = this.Status | FLAG_CARRY;
    } else {
        this.Status = this.Status & ^FLAG_CARRY;
    }
}

func (this *CPU) adc(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    value := this.Mem.read(addr);
    sum := uint16(this.Reg_a) + uint16(value) + uint16(this.Status&FLAG_CARRY)
    carry := sum > 0xFF
    overflow := (this.Reg_a^value)&0x80 == 0 && (this.Reg_a^byte(sum))&0x80 != 0
    this.Reg_a = byte(sum)
    
    this.SetFlags(this.Reg_a)
    if carry {
        this.Status |= FLAG_CARRY
    } else {
        this.Status &^= FLAG_CARRY
    }
    if overflow {
        this.Status |= FLAG_OVERFLOW
    } else {
        this.Status &^= FLAG_OVERFLOW
    }
}

func (this *CPU) lda(mode AddressingMode) {
    var addr uint16 = this.GetOperandAddress(mode);
    var val = this.Mem.read(addr);
    this.Reg_a = val;
    this.SetFlags(this.Reg_a);
}

func (this *CPU) tax() {
    this.Reg_x = this.Reg_a;
    this.SetFlags(this.Reg_x);
}

func (this *CPU) inx() {
    this.Reg_x += 1;
    this.SetFlags(this.Reg_x); 
}

func (this *CPU) and(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    tmp_data := this.Mem.read(addr);
    this.Reg_a &= tmp_data;
    this.SetFlags(this.Reg_a)
}

func (this *CPU) asl(mode AddressingMode) {
    var data uint8;
    if mode == Accumulator {
        data = this.Reg_a
        if data >> 7 == 1 {
            this.setCFlag(true)
        } else {
            this.setCFlag(false)
        }
        data = data << 1
        this.Reg_a = data
        this.SetFlags(this.Reg_a);
    } else {
        addr := this.GetOperandAddress(mode);
        data = this.Mem.read(addr);
        if data >> 7 == 1 {
            this.setCFlag(true)
        } else {
            this.setCFlag(false)
        }
        data = data << 1
        this.Mem.write(addr, data);
        this.SetFlags(data);
    }
}

func (this *CPU) bcc() {
    if this.Status & FLAG_CARRY == 0 {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr
    }
}

func (this *CPU) bcs() {
    if this.Status & FLAG_CARRY == FLAG_CARRY {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr
    }
}

func (this *CPU) beq() {
    if this.Status & FLAG_ZERO == FLAG_ZERO {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr
    }
}

func (this *CPU) bit(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    res := this.Reg_a & data;
    if res == 0 {
        this.Status |= FLAG_ZERO;
    } else {
        this.Status &^= FLAG_ZERO;
    }
    this.Status = (this.Status &^ FLAG_NEGATIVE) | (data & FLAG_NEGATIVE)
    this.Status = (this.Status &^ FLAG_OVERFLOW) | ((data >> 6) & 1 << 6)
}

func (this *CPU) bmi() {
    if this.Status & FLAG_NEGATIVE == FLAG_NEGATIVE {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr       
    }
}

func (this *CPU) bne() {
    if this.Status & FLAG_ZERO == 0 {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr       
    }
}

// TODO: Have to check whether if 0 FLAG means Positive or not
func (this *CPU) bpl() {
    if this.Status & FLAG_NEGATIVE == 0 {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr       
    }
}

func (this *CPU) bvc() {
    if this.Status & FLAG_OVERFLOW == 0 {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr       
    }
}

func (this *CPU) bvs() {
    if this.Status & FLAG_OVERFLOW == FLAG_OVERFLOW {
        jump := this.Mem.read(this.PC)
        jump_addr := this.PC + 1 + uint16(jump)
        this.PC = jump_addr       
    }
}

func (this *CPU) clc() {
    this.Status = this.Status & ^FLAG_CARRY;
}

func (this *CPU) cld() {
    this.Status = this.Status & ^FLAG_DECIMAL;
}

func (this *CPU) cli() {
    this.Status = this.Status & ^FLAG_INTERRUPT;
}

func (this *CPU) clv() {
    this.Status = this.Status & ^FLAG_OVERFLOW;
}

func (this *CPU) cmp(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    if (this.Reg_a >= data) {
        this.Status = this.Status | FLAG_CARRY;
    } else {
        this.Status = this.Status & ^FLAG_CARRY;
    }
    
    this.SetFlags(this.Reg_a - data)
}

func (this *CPU) cpx(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    if (this.Reg_x >= data) {
        this.Status = this.Status | FLAG_CARRY;
    } else {
        this.Status = this.Status & ^FLAG_CARRY;
    }
    
    this.SetFlags(this.Reg_x - data)
}

func (this *CPU) cpy(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    if (this.Reg_y >= data) {
        this.Status = this.Status | FLAG_CARRY;
    } else {
        this.Status = this.Status & ^FLAG_CARRY;
    }
    
    this.SetFlags(this.Reg_y - data)
}

func (this *CPU) dec(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    data--;
    this.Mem.write(addr, data);
    this.SetFlags(data);
}

func (this *CPU) dex() {
    this.Reg_x--;
    this.SetFlags(this.Reg_x);
}

func (this *CPU) dey() {
    this.Reg_y--;
    this.SetFlags(this.Reg_y);
}

func (this *CPU) eor(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    this.Reg_a = this.Reg_a ^ data;
    this.SetFlags(this.Reg_a);
}

func (this *CPU) inc(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    data++;
    this.Mem.write(addr, data);
    this.SetFlags(data);
}

func (this *CPU) iny() {
    this.Reg_y += 1;
    this.SetFlags(this.Reg_y); 
}

func (this *CPU) jmp(mode AddressingMode) {
    if mode == Absolute {
        addr := this.Mem.readU16(this.PC);
        this.PC = addr;
    } else if mode == Indirect {
        addr := this.Mem.readU16(this.PC);
        var res uint16;
        if addr & 0x00FF == 0x00FF {
            lo := this.Mem.read(addr);
            hi := this.Mem.read(addr & 0x00FF)
            res = uint16(hi) << 8 | uint16(lo)
        } else {
            res = this.Mem.readU16(addr)
        }
        this.PC = res;
    }
}

func (this *CPU) jsr() {
    this.stack_push_U16(this.PC + 1) 
    addr := this.Mem.readU16(this.PC)
    this.PC = addr
}

func (this *CPU) ldx(mode AddressingMode) {
    var addr uint16 = this.GetOperandAddress(mode);
    var val = this.Mem.read(addr);
    this.Reg_x = val;
    this.SetFlags(this.Reg_x);
}

func (this *CPU) ldy(mode AddressingMode) {
    var addr uint16 = this.GetOperandAddress(mode);
    var val = this.Mem.read(addr);
    this.Reg_y = val;
    this.SetFlags(this.Reg_y);
}

func (this *CPU) lsr(mode AddressingMode) {
    var data uint8;
    if mode == Accumulator {
        data = this.Reg_a
        if data & 1 == 1 {
            this.setCFlag(true)
        } else {
            this.setCFlag(false)
        }
        data = data >> 1
        this.Reg_a = data
        this.SetFlags(this.Reg_a);
    } else {
        addr := this.GetOperandAddress(mode);
        data = this.Mem.read(addr);
        if data & 1 == 1 {
            this.setCFlag(true);
        } else {
            this.setCFlag(false);
        }
        data = data >> 1
        this.Mem.write(addr, data);
        this.SetFlags(data);
    }
}

func (this *CPU) ora(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    data := this.Mem.read(addr);
    this.Reg_a = this.Reg_a | data;
    this.SetFlags(this.Reg_a);
}

func (this *CPU) pha() {
    this.stack_push(this.Reg_a);
}

func (this *CPU) php() {
    tmp_status := this.Status;
    tmp_status = tmp_status | FLAG_BREAK | FLAG_BREAK2;
    this.stack_push(tmp_status);
}

func (this *CPU) pla() {
    this.Reg_a = this.stack_pop();
    this.SetFlags(this.Reg_a);
}

// PLP - Pull Processor Status from stack
func (this *CPU) plp() {
    this.Status = this.stack_pop()
    this.Status &^= FLAG_BREAK
    this.Status |= FLAG_BREAK2
}

// ROL - Rotate Left
func (this *CPU) ROL(mode AddressingMode) {
    if mode == Accumulator {
        carry := (this.Reg_a >> 7) & 1
        this.Reg_a = (this.Reg_a << 1) | (this.Status & FLAG_CARRY)
        if carry != 0 {
            this.Status |= FLAG_CARRY
        } else {
            this.Status &^= FLAG_CARRY
        }
        this.SetFlags(this.Reg_a)
    } else {
        addr := this.GetOperandAddress(mode)
        value := this.Mem.read(addr)
        carry := (value >> 7) & 1
        value = (value << 1) | (this.Status & FLAG_CARRY)
        if carry != 0 {
            this.Status |= FLAG_CARRY
        } else {
            this.Status &^= FLAG_CARRY
        }
        this.Mem.write(addr, value)
        this.SetFlags(value)
    }
}

// ROR - Rotate Right
func (this *CPU) ROR(mode AddressingMode) {
    if mode == Accumulator {
        carry := this.Reg_a & 1
        this.Reg_a = (this.Reg_a >> 1) | ((this.Status & FLAG_CARRY) << 7)
        if carry != 0 {
            this.Status |= FLAG_CARRY
        } else {
            this.Status &^= FLAG_CARRY
        }
        this.SetFlags(this.Reg_a)
    } else {
        addr := this.GetOperandAddress(mode)
        value := this.Mem.read(addr)
        carry := value & 1
        value = (value >> 1) | ((this.Status & FLAG_CARRY) << 7)
        if carry != 0 {
            this.Status |= FLAG_CARRY
        } else {
            this.Status &^= FLAG_CARRY
        }
        this.Mem.write(addr, value)
        this.SetFlags(value)
    }
}

// RTI - Return from Interrupt
func (this *CPU) RTI() {
    this.Status = this.stack_pop()
    this.Status &^= FLAG_BREAK
    this.Status |= FLAG_BREAK2
    this.PC = this.stack_pop_U16()
}

// RTS - Return from Subroutine
func (this *CPU) RTS() {
    this.PC = this.stack_pop_U16() + 1
}

// SBC - Subtract with Carry
func (this *CPU) SBC(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    value := this.Mem.read(addr);
    value = ^value // Invert the bits of the value (2's complement for subtraction)
    sum := uint16(this.Reg_a) + uint16(value) + uint16(this.Status&FLAG_CARRY)
    carry := sum > 0xFF
    overflow := (this.Reg_a^value)&0x80 == 0 && (this.Reg_a^byte(sum))&0x80 != 0
    this.Reg_a = byte(sum)

    this.SetFlags(this.Reg_a)
    if carry {
        this.Status |= FLAG_CARRY
    } else {
        this.Status &^= FLAG_CARRY
    }
    if overflow {
        this.Status |= FLAG_OVERFLOW
    } else {
        this.Status &^= FLAG_OVERFLOW
    }
}

// SEC - Set Carry Flag
func (this *CPU) sec() {
    this.Status |= FLAG_CARRY
}

// SED - Set Decimal Flag
func (this *CPU) SED() {
    this.Status |= FLAG_DECIMAL
}

// SEI - Set Interrupt Disable
func (this *CPU) SEI() {
    this.Status |= FLAG_INTERRUPT
}

func (this *CPU) sta(mode AddressingMode) {
    var addr uint16 = this.GetOperandAddress(mode);
    this.Mem.write(addr, this.Reg_a);
}

// STX - Store X Register
func (this *CPU) stx(mode AddressingMode) {
    addr := this.GetOperandAddress(mode);
    this.Mem.write(addr, this.Reg_x)
}

// STY - Store Y Register
func (this *CPU) sty(mode AddressingMode) {
    addr := this.GetOperandAddress(mode)
    this.Mem.write(addr, this.Reg_y)
}

// TAX - Transfer Accumulator to X
func (this *CPU) TAX() {
    this.Reg_x = this.Reg_a
    this.SetFlags(this.Reg_x)
}

// TAY - Transfer Accumulator to Y
func (this *CPU) TAY() {
    this.Reg_y = this.Reg_a
    this.SetFlags(this.Reg_y)
}

// TSX - Transfer Stack Pointer to X
func (this *CPU) TSX() {
    this.Reg_x = this.SP
    this.SetFlags(this.Reg_x)
}

// TXA - Transfer X to Accumulator
func (this *CPU) TXA() {
    this.Reg_a = this.Reg_x
    this.SetFlags(this.Reg_a)
}

// TXS - Transfer X to Stack Pointer

func (this *CPU) TXS() {
    this.SP = this.Reg_x
}

// TYA - Transfer Y to Accumulator
func (this *CPU) TYA() {
    this.Reg_a = this.Reg_y
    this.SetFlags(this.Reg_a)
}

