package emulator 

/*
Flags are based on https://www.nesdev.org/wiki/Status_flags
Notice that bit-5 is not a flag.
*/

const (
    FLAG_CARRY      byte = 1 << 0  // 0000 0001
    FLAG_ZERO       byte = 1 << 1  // 0000 0010
    FLAG_INTERRUPT  byte = 1 << 2  // 0000 0100
    FLAG_DECIMAL    byte = 1 << 3  // 0000 1000
    FLAG_BREAK      byte = 1 << 4  // 0001 0000
    FLAG_BREAK2     byte = 1 << 5  // 0010 0000
    FLAG_OVERFLOW   byte = 1 << 6  // 0100 0000 
    FLAG_NEGATIVE   byte = 1 << 7  // 1000 0000
)

type AddressingMode byte 

const (
    Immediate AddressingMode = iota
    ZeroPage
    ZeroPageX
    ZeroPageY
    Absolute
    AbsoluteX
    AbsoluteY
    Indirect
    IndirectX
    IndirectY
    Accumulator
    Implied
    Relative
)

const STACK uint16 = 0x0100;
const STACK_RESET uint8 = 0xFD;

// Emulates the 6502 CPU.
type CPU struct {
    A uint8         // Accumulator
    X uint8         // Index Register x
    Y uint8         // Index Register y
    SP uint8        // Stack pointer.
    STATUS byte     // Flags for program.
    PC uint16       // Program Counter
    Bus *BUS        // The Bus connecting CPU to everything else.

    cycles_left uint8 // How many cycles left.
}

func MakeCPU() *CPU {
    cpu := CPU{A: 0, X: 0, Y: 0, SP: STACK_RESET, STATUS: 0, PC: 0, Bus: GetBus()};
    cpu.Bus.BusSetCPU(&cpu);
    return &cpu;
}

func (this *CPU) Write(addr uint16, data uint8) {
    this.Bus.Write(addr, data);
}

func (this *CPU) Read(addr uint16) uint8 {
    return this.Bus.Read(addr);
}

func (this *CPU) Read_u16(addr uint16) uint16 {
    lo := uint16(this.Bus.Read(addr))
    hi := uint16(this.Bus.Read(addr + 1))
    return hi << 8 | lo
}

func (this *CPU) ContainsFlag(flag uint8) bool {
    return this.STATUS & flag != 0
}

func (this *CPU) SetFlag(flag uint8, cond bool) {
    if cond {
        this.STATUS |= flag
    } else {
        this.STATUS = this.STATUS & ^flag
    }
}

func (this *CPU) push(data uint8) {
    this.Write(STACK + uint16(this.SP), data)
    this.SP--       // Recall that decrementing SP is allocating more space
}

func (this *CPU) push_u16(data uint16) {
    this.push(uint8(data >> 8))
    this.push(uint8(data))
}

func (this *CPU) pop() uint8 {
    this.SP++
    return this.Read(STACK + uint16(this.SP))
}

func (this *CPU) pop_u16() uint16 {
    lo := this.pop()
    hi := this.pop()
    return uint16(hi) << 8 | uint16(lo)
}

// Represents what happens in a single clock cycle.
func (this *CPU) Tick() {
    if this.cycles_left == 0 {
        opcode := this.Read(this.PC);
        this.PC++;

        var instr Instruction = Instructions[opcode];
        cycles := instr.cycles;

        operand := this.fetchOperand(instr.mode)
        instr.handler(this, operand)
        this.cycles_left += cycles - 1
    }
    this.cycles_left--;
}

func (cpu *CPU) fetchOperand(mode AddressingMode) Operand {
	switch mode {
	case Implied:
		return Operand{mode: mode}
	case Accumulator:
		return Operand{mode: mode}
	case Immediate:
		addr := cpu.PC
		cpu.PC++
		return Operand{mode: mode,address: addr}
	case ZeroPage:
		addr := uint16(cpu.Read(cpu.PC))
		cpu.PC++
		return Operand{mode: mode, address: addr}
	case ZeroPageX:
		addr := (uint16(cpu.Read(cpu.PC)) + uint16(cpu.X)) & 0x00FF
		cpu.PC++
		return Operand{mode: mode, address: addr}
	case ZeroPageY:
		addr := (uint16(cpu.Read(cpu.PC)) + uint16(cpu.Y)) & 0x00FF
		cpu.PC++
		return Operand{mode: mode, address: addr}
	case Absolute:
		addr := cpu.Read_u16(cpu.PC)
		cpu.PC += 2
		return Operand{mode: mode, address: addr,}
	case AbsoluteX:
		addr := cpu.Read_u16(cpu.PC)
		cpu.PC += 2
		addrX := addr + uint16(cpu.X)
        do_cycle := addr&0xFF00 != addrX&0xFF00
		return Operand{mode: mode, address: addrX, extra_cycle: do_cycle}
	case AbsoluteY:
		addr := cpu.Read_u16(cpu.PC)
		cpu.PC += 2
		addrY := addr + uint16(cpu.Y)
        do_cycle := addr&0xFF00 != addrY&0xFF00
        return Operand{mode: mode,address: addrY, extra_cycle: do_cycle}
	case Indirect:
		ptrAddr := cpu.Read_u16(cpu.PC)
		addr := cpu.Read_u16(ptrAddr)
		cpu.PC += 2
		return Operand{mode: mode, address: addr}
	case IndirectX:
		ptrAddr := (uint16(cpu.Read(cpu.PC)) + uint16(cpu.X)) & 0x00FF
		addr := cpu.Read_u16(ptrAddr)
		cpu.PC++
		return Operand{mode: mode,address: addr,}
	case IndirectY:
		ptrAddr := uint16(cpu.Read(cpu.PC))
		cpu.PC++
		addr := cpu.Read_u16(ptrAddr)
		addrY := addr + uint16(cpu.Y)
        do_cycle := addrY&0xFF00 != addr&0xFF00
		return Operand{mode: mode,address: addrY, extra_cycle: do_cycle}
	case Relative:
		rel := uint16(cpu.Read(cpu.PC))
		cpu.PC++
		if rel&(1<<7) != 0 {
			rel |= 0xFF00
		}
		addr := cpu.PC + rel
        do_cycle := addr & 0xFF00 != cpu.PC & 0xFF00
		return Operand{mode: mode,address: addr, extra_cycle: do_cycle}
	default:
		panic("invalid addressing mode")
	}
}
