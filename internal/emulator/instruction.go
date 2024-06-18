package emulator 

// Sets the flags for CPU based on the output of an operation (result)
func (this *CPU) SetZNFlag(result uint8) {
    // Sets the Z flag (ZERO FLAG)
    if result == 0 {
        this.STATUS = this.STATUS | FLAG_ZERO;
    } else {
        this.STATUS = this.STATUS & ^FLAG_ZERO; // Go does not have a bitwise NOT. So use XOR trick.
    }

    // Sets the N flag (NEGATIVE FLAG)
    if result & 0b1000_0000 != 0 {
        this.STATUS = this.STATUS | FLAG_NEGATIVE;
    } else {
        this.STATUS = this.STATUS & ^FLAG_NEGATIVE;
    }
}

func (this *CPU) SetCFlag(active bool) {
    if active {
        this.STATUS = this.STATUS | FLAG_CARRY;
    } else {
        this.STATUS = this.STATUS & ^FLAG_CARRY;
    }
}

func (this *CPU) carried() bool {
    return this.STATUS & FLAG_CARRY != 0
}

// ADC - Add with Carry
func adc(cpu *CPU, op Operand) {
    a := uint16(cpu.A);
    b := uint16(cpu.Read(op.address))
    var carry uint16
    if cpu.carried() {
        carry = 1
    } else {
        carry = 0
    }
    r := a + b + carry
	overflow := (a^b) & 0x80 == 0 && (a^r) & 0x80 != 0

	cpu.SetFlag(FLAG_CARRY, r > 0xFF)
	cpu.SetFlag(FLAG_OVERFLOW, overflow)
	cpu.A = uint8(r)
	cpu.SetZNFlag(cpu.A)

	if op.extra_cycle {
		cpu.cycles_left++
	}
}

// AND - Logical AND
func and(cpu *CPU, op Operand) {
    val := cpu.Read(op.address);
    cpu.A &= val;
    cpu.SetZNFlag(cpu.A)
    if op.extra_cycle {
        cpu.cycles_left++;
    }
}

// ASL - Arithmetic Shift Left.
func asl(cpu *CPU, op Operand) {
    if op.mode == Accumulator {
        val := cpu.A
        cpu.SetCFlag(val & 0x80 != 0)
        val = val << 1
        cpu.SetZNFlag(cpu.A)
        cpu.A = val
    } else {
        val := cpu.Read(op.address)
        cpu.SetCFlag(val & 0x80 != 0)
        val = val << 1
        cpu.SetZNFlag(cpu.A)
        cpu.Write(op.address, val)
    }
}

// BCC - Branch if Carry Clear
func bcc(cpu *CPU, op Operand) {
    if !cpu.ContainsFlag(FLAG_CARRY) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// BCS - Branch if Carry Set
func bcs(cpu *CPU, op Operand) {
    if cpu.ContainsFlag(FLAG_CARRY) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// BEQ - Branch if Equal
func beq(cpu *CPU, op Operand) {
    if cpu.ContainsFlag(FLAG_ZERO) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// BIT - Bit Test
func bit(cpu *CPU, op Operand) {
    val := cpu.Read(op.address);
    cpu.SetFlag(FLAG_ZERO, val & cpu.A == 0)
    cpu.SetFlag(FLAG_OVERFLOW, val & (1 << 6) != 0)
    cpu.SetFlag(FLAG_NEGATIVE, val & (1 << 7) != 0)
}

// BMI - Branch if Minus
func bmi(cpu *CPU, op Operand) {
    if cpu.ContainsFlag(FLAG_NEGATIVE) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}
// BNE - Branch if Not Equal
func bne(cpu *CPU, op Operand) {
    if !cpu.ContainsFlag(FLAG_ZERO) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// BPL - Branch if Positive
func bpl(cpu *CPU, op Operand) {
    if !cpu.ContainsFlag(FLAG_NEGATIVE) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// BRK - BREAK, HALT, STOP STOP STOP. There are actually more to this
func brk(cpu *CPU, op Operand) {
    cpu.push_u16(cpu.PC + 1)
    cpu.push(cpu.STATUS) // TODO: Recheck this
    cpu.SetFlag(FLAG_INTERRUPT, true)
    cpu.PC = cpu.Read_u16(0xFFFE)
    return
}

// BVC - Branch if Overflow Clear
func bvc(cpu *CPU, op Operand) {
    if !cpu.ContainsFlag(FLAG_OVERFLOW) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// BVS - Branch if Overflow Set
func bvs(cpu *CPU, op Operand) {
    if cpu.ContainsFlag(FLAG_OVERFLOW) {
        cpu.PC = op.address;
        cpu.cycles_left++;
        
        if op.extra_cycle {
            cpu.cycles_left += 1; // TODO: Recheck if +1 or +2
        }
    }
}

// CLC - Clear Carry Flag
func clc(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_CARRY, false)
}

// CLD - Clear Decimal Mode
func cld(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_DECIMAL, false)
}

// CLI - Clear Interrupt Disable
func cli(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_INTERRUPT, false)
}

// CLV - Clear Overflow FLag
func clv(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_OVERFLOW, false)
}

// CMP - Compare
func cmp(cpu *CPU, op Operand) {
    acc := cpu.A
    data := cpu.Read(op.address)
    
    cpu.SetFlag(FLAG_CARRY, acc >= data)
    res := acc - data
    cpu.SetZNFlag(res)

    if op.extra_cycle {
        cpu.cycles_left++
    }
}

// CPX - Compare X Register
func cpx(cpu *CPU, op Operand) {
    acc := cpu.X
    data := cpu.Read(op.address)
    
    cpu.SetFlag(FLAG_CARRY, acc >= data)
    res := acc - data
    cpu.SetZNFlag(res)

    if op.extra_cycle {
        cpu.cycles_left++
    }
}

// CPY - Compare Y Register
func cpy(cpu *CPU, op Operand) {
    acc := cpu.Y
    data := cpu.Read(op.address)
    
    cpu.SetFlag(FLAG_CARRY, acc >= data)
    res := acc - data
    cpu.SetZNFlag(res)

    if op.extra_cycle {
        cpu.cycles_left++
    }
}

// DEC - Decrement val in Memory
func dec(cpu *CPU, op Operand) {
    val := cpu.Read(op.address)
    val--
    cpu.Write(op.address, val)
    cpu.SetZNFlag(val)
}

// DEX - Decrement val in X Register
func dex(cpu *CPU, op Operand) {
    cpu.X--
    cpu.SetZNFlag(cpu.X)
}

// DEY - Decrement val in Y Register
func dey(cpu *CPU, op Operand) {
    cpu.Y--
    cpu.SetZNFlag(cpu.Y)
}

// EOR - Exclusive OR
func eor(cpu *CPU, op Operand) {
    acc := cpu.A
    val := cpu.Read(op.address)
    res := acc ^ val
    cpu.SetZNFlag(res)
    cpu.A = res

    if op.extra_cycle {
        cpu.cycles_left++
    }
}

// INC - Increment a val in memory
func inc(cpu *CPU, op Operand) {
    val := cpu.Read(op.address) + 1
    cpu.Write(op.address, val)
    cpu.SetZNFlag(val)
}

// INX - Increment X register
func inx(cpu *CPU, op Operand) {
    cpu.X++
    cpu.SetZNFlag(cpu.X)
}

func iny(cpu *CPU, op Operand) {
    cpu.Y++
    cpu.SetZNFlag(cpu.Y)
}

// JMP - Change PC to addr
func jmp(cpu *CPU, op Operand) {
    cpu.PC = op.address
}

// JSR - Jump to Subroutine. Like CALL for y86. Push addr to stack and then jumps
func jsr(cpu *CPU, op Operand) {
    cpu.push_u16(cpu.PC-1)
    cpu.PC = op.address
}

// LDA - Load value of memory to Register A
func lda(cpu *CPU, op Operand) {
    val := cpu.Read(op.address)
    cpu.A = val
    cpu.SetZNFlag(cpu.A)
    if op.extra_cycle {
        cpu.cycles_left++
    }
}

func ldx(cpu *CPU, op Operand) {
    val := cpu.Read(op.address)
    cpu.X = val
    cpu.SetZNFlag(cpu.X)
    if op.extra_cycle {
        cpu.cycles_left++
    }
}

func ldy(cpu *CPU, op Operand) {
    val := cpu.Read(op.address)
    cpu.Y = val
    cpu.SetZNFlag(cpu.Y)
    if op.extra_cycle {
        cpu.cycles_left++
    }
}

// LSR - Logical Shift Right
func lsr(cpu *CPU, op Operand) {
    if op.mode == Accumulator {
        val := cpu.A
        cpu.SetCFlag(val & 0x01 != 0)
        val = val >> 1
        cpu.SetZNFlag(cpu.A)
        cpu.A = val
    } else {
        val := cpu.Read(op.address)
        cpu.SetCFlag(val & 0x01 != 0)
        val = val >> 1
        cpu.SetZNFlag(cpu.A)
        cpu.Write(op.address, val)
    }
}

func nop(cpu *CPU, op Operand) {
    return
}

// ORA - Logical Inclusive OR
func ora(cpu *CPU, op Operand) {
    cpu.A = cpu.A | cpu.Read(op.address)
    cpu.SetZNFlag(cpu.A)
    if op.extra_cycle {
        cpu.cycles_left++
    }
}

func pha(cpu *CPU, op Operand) {
    cpu.push(cpu.A)
}

func php(cpu *CPU, op Operand) {
    cpu.push(cpu.STATUS) // recheck.
}

func pla(cpu *CPU, op Operand) {
    cpu.A = cpu.pop()
    cpu.SetZNFlag(cpu.A)
}

func plp(cpu *CPU, op Operand) {
    cpu.STATUS = cpu.pop()
}

func rol(cpu *CPU, op Operand) {
    var carry uint8
        if cpu.carried() {
            carry = 1
        } else {
            carry = 0
        }
    if op.mode == Accumulator {
        data := cpu.A
        cpu.SetFlag(FLAG_CARRY, data & 0x80 != 0)
        data = data << 1 | carry
        cpu.SetZNFlag(data)
        cpu.A = data
    } else {
        data := cpu.Read(op.address)
        cpu.SetFlag(FLAG_CARRY, data & 0x80 != 0)
        data = data << 1 | carry
        cpu.SetZNFlag(data)
        cpu.Write(op.address, data)
    }
}

func ror(cpu *CPU, op Operand) {
    var carry uint8
        if cpu.carried() {
            carry = 1
        } else {
            carry = 0
        }
    if op.mode == Accumulator {
        data := cpu.A
        cpu.SetFlag(FLAG_CARRY, data & 0x01 != 0)
        data = data >> 1 | carry << 7
        cpu.SetZNFlag(data)
        cpu.A = data
    } else {
        data := cpu.Read(op.address)
        cpu.SetFlag(FLAG_CARRY, data & 0x01 != 0)
        data = data >> 1 | carry << 7
        cpu.SetZNFlag(data)
        cpu.Write(op.address, data)
    }
}

func rti(cpu *CPU, op Operand) {
    cpu.STATUS = cpu.pop()          // Sketched out. But just try.
    cpu.SetFlag(FLAG_BREAK, false)
    cpu.PC = cpu.pop_u16()
}

func rts(cpu *CPU, op Operand) {
    addr := cpu.pop_u16()
    cpu.PC = addr + 1
}

// SBC - Subtract with Carry
func sbc(cpu *CPU, op Operand) {
    var carry uint16
    if cpu.carried() {
        carry = 1
    } else {
        carry = 0
    }

    a := uint16(cpu.A)
    b := uint16(cpu.Read(op.address))
    r := a - b - (1-carry)
    overflow := (a^b)&0x80 != 0 && (a^r)&0x80 != 0

    cpu.SetFlag(FLAG_CARRY, r < 0x100)
	cpu.SetFlag(FLAG_OVERFLOW, overflow)
	cpu.A = uint8(r)
	cpu.SetZNFlag(cpu.A)

    if op.extra_cycle {
        cpu.cycles_left++
    }
}


func sec(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_CARRY, true) 
}

func sed(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_DECIMAL, true) 
}

func sei(cpu *CPU, op Operand) {
    cpu.SetFlag(FLAG_INTERRUPT, true) 
}

func sta(cpu *CPU, op Operand) {
    cpu.Write(op.address, cpu.A)
}

func stx(cpu *CPU, op Operand) {
    cpu.Write(op.address, cpu.X)
}

func sty(cpu *CPU, op Operand) {
    cpu.Write(op.address, cpu.Y) 
}

func tax(cpu *CPU, op Operand) {
    cpu.X = cpu.A
    cpu.SetZNFlag(cpu.X)
}

func tay(cpu *CPU, op Operand) {
    cpu.Y = cpu.A
    cpu.SetZNFlag(cpu.Y)
}

func tsx(cpu *CPU, op Operand) {
    cpu.X = cpu.SP
    cpu.SetZNFlag(cpu.X)
}

func txa(cpu *CPU, op Operand) {
    cpu.A = cpu.X
    cpu.SetZNFlag(cpu.A)
}

func txs(cpu *CPU, op Operand) {
    cpu.SP = cpu.X
}

func tya(cpu *CPU, op Operand) {
    cpu.A = cpu.Y
    cpu.SetZNFlag(cpu.A)
}
