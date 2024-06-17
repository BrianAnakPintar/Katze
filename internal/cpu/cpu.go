package cpu

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
    NoneAddressing
)

const STACK uint16 = 0x0100;
const STACK_RESET uint8 = 0xFD;

// Emulates the 6502 CPU.
type CPU struct {
    Reg_a byte      // Accumulator
    Reg_x byte      // Index Register x
    Reg_y byte      // Index Register y
    SP byte         // Stack pointer.
    Status byte     // Flags for program.
    PC uint16       // Program Counter
    Mem *Memory     // The RAM :O
}

func MakeCPU() *CPU {
    cpu := CPU{Reg_a: 0, Reg_x: 0, Reg_y: 0, SP: STACK_RESET, Status: 0, PC: 0, Mem: MakeMemory(65536)};
    return &cpu;
}

// Fetches an instruction and executes them in a single tick.
// Note that this is not a Pipelined implementation.
func (this *CPU) Run() {
    for {
        var code byte = this.Mem.read(this.PC);
        var instr instruction = Instructions[code]
        this.PC++;
        switch instr.name {
        case "ADC":
        case "AND":
            this.and(instr.mode);
        case "ASL":
            this.asl(instr.mode);
        case "BCC":
            this.bcc();
        case "BCS":
            this.bcs();
        case "BEQ":
            this.beq();
        case "BIT":
            this.bit(instr.mode);
        case "BMI":
            this.bmi();
        case "BNE":
            this.bne();
        case "BPL":
            this.bpl();
        case "BRK":
            return
        case "BVC":
            this.bvc();
        case "BVS":
            this.bvs();
        case "CLC":
            this.clc();
        case "CLD":
            this.cld();
        case "CLI":
            this.cli();
        case "CLV":
            this.clv();
        case "CMP":
            this.cmp(instr.mode);
        case "CPX":
            this.cpx(instr.mode);
        case "CPY":
            this.cpy(instr.mode);
        case "DEC":
            this.dec(instr.mode);
        case "DEX":
            this.dex();
        case "DEY":
            this.dey();
        case "EOR":
            this.eor(instr.mode);
        case "INC":
            this.inc(instr.mode);
        case "INX":
            this.inx();
        case "INY":
            this.iny();
        case "JMP":
            this.jmp(instr.mode);
        case "JSR":
            this.jsr();
        case "LDA":
            this.lda(instr.mode);
        case "LDX":
            this.ldx(instr.mode);
        case "LDY":
            this.ldy(instr.mode);
        case "LSR":
            this.lsr(instr.mode);
        case "NOP":
            break // Do nothing
        case "ORA":
            this.ora(instr.mode);
        case "PHA":
            this.pha();
        case "PHP":
            this.php();
        case "PLA":
            this.pla();
        case "PLP":
            this.plp();
        case "ROL":
            this.ROL(instr.mode);
        case "ROR":
            this.ROR(instr.mode);
        case "RTI":
            this.RTI();
        case "RTS":
            this.RTS();
        case "SBC":
            this.SBC(instr.mode)
        case "SEC":
            this.sec();
        case "SED":
            this.SED();
        case "SEI":
            this.SEI();
        case "STA":
            this.sta(instr.mode)
        case "STX":
            this.stx(instr.mode)
        case "STY":
            this.sty(instr.mode)
        case "TAX":
            this.TAX()
        case "TAY":
            this.TAY()
        case "TSX":
            this.TSX()
        case "TXA":
            this.TXA()
        case "TXS":
            this.TXS()
        case "TYA":
            this.TYA()
        }
        this.PC += uint16(instr.size);
        break
    }
}

func (this *CPU) reset() {
    this.Reg_a = 0;
    this.Reg_x = 0;
    this.Reg_y = 0;
    this.Status = 0;
    this.SP = STACK_RESET;
    this.PC = this.Mem.readU16(0xFFFC);
}

func (this *CPU) Load_and_run(program []byte) {
    this.Mem.load(program);
    this.reset();
    this.Run();
}

/*
Note for self: Recall in CPSC 213, our assembly has different versions of `ld`.
For example: (a) `ld $5, r0` (b) `ld r1, r0` ... (n) `ld 4(r1), r0`

This is basically the same thing. In which there are different versions, notice that the only thing
it really changes is the "address from where we read/store in case of st"
*/
func (this *CPU) GetOperandAddress(mode AddressingMode) uint16 {
    switch mode {
    case Immediate:
        return this.PC;
    case ZeroPage:
        return uint16(this.Mem.read(this.PC));
    case Absolute:
        return this.Mem.readU16(this.PC);
    case ZeroPageX:
        pos := this.Mem.read(this.PC);
        addr := uint16(pos + this.Reg_x);
        return addr;
    case ZeroPageY:
        pos := this.Mem.read(this.PC);
        addr := uint16(pos + this.Reg_y);
        return addr;
    case AbsoluteX:
        pos := this.Mem.readU16(this.PC);
        addr := pos + uint16(this.Reg_x);
        return addr;
    case AbsoluteY:
        pos := this.Mem.readU16(this.PC);
        addr := pos + uint16(this.Reg_y);
        return addr;
    case IndirectX:
        base := this.Mem.read(this.PC);
        ptr := uint8(base) + this.Reg_x;
        lo := this.Mem.read(uint16(ptr));
        hi := this.Mem.read(uint16(ptr) + 1);
        return (uint16(hi) << 8) | (uint16(lo));
    case IndirectY:
        base := this.Mem.read(this.PC);
        lo := this.Mem.read(uint16(base));
        hi := this.Mem.read(uint16(base) + 1);
        deref_base := (uint16(hi) << 8) | (uint16(lo));       
        return deref_base + uint16(this.Reg_y);
    case NoneAddressing:
        this.Scream();
    }
    return 0;
}

func (this *CPU) stack_push(data byte) {
    this.Mem.write(STACK + uint16(this.SP), data);
    this.SP = this.SP - 1;
}

func (this *CPU) stack_pop() byte {
    this.SP = this.SP + 1;
    return this.Mem.read(STACK + uint16(this.SP));
}

func (this *CPU) stack_push_U16(data uint16) {
    hi := uint8(data >> 8)
    lo := uint8(data & 0xff)
    this.stack_push(hi);
    this.stack_push(lo);
}

func (this *CPU) stack_pop_U16() uint16 {
    lo := this.stack_pop();
    hi := this.stack_pop();
    return uint16(hi) << 8 | uint16(lo);
}

func (e *CPU) Scream() {
    panic("It's over... GGs")
}
