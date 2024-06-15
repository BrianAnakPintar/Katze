package cpu

import "fmt"

/*
Flags are based on https://www.nesdev.org/wiki/Status_flags
Notice that bit-5 is not a flag.
*/

const (
    FLAG_CARRY      uint8 = 1 << 0  // 0000 0001
    FLAG_ZERO       uint8 = 1 << 1  // 0000 0010
    FLAG_INTERRUPT  uint8 = 1 << 2  // 0000 0100
    FLAG_DECIMAL    uint8 = 1 << 3  // 0000 1000
    FLAG_BREAK      uint8 = 1 << 4  // 0001 0000
    FLAG_OVERFLOW   uint8 = 1 << 6  // 0100 0000 
    FLAG_NEGATIVE   uint8 = 1 << 7  // 1000 0000
)

// Emulates the 6502 CPU.
type CPU struct {
    Reg_a uint8     // Accumulator
    Reg_x uint8     // Index Register x
    Reg_y uint8     // Index Register y
    SP uint8        // Stack pointer.
    Status uint8    // Flags for program.
    PC uint16       // Program Counter
    Mem *Memory     // The RAM :O
}

func MakeCPU() *CPU {
    cpu := CPU{Reg_a: 0, Reg_x: 0, Reg_y: 0, SP: 0, Status: 0, PC: 0, Mem: MakeMemory(65536)};
    return &cpu;
}

func (this *CPU) Run() {
    for {
        var op_code uint8 = this.Mem.read(this.PC);
        this.PC++;
        
        switch op_code {
        case 0xA9:  // LDA #$ff
            var val uint8 = this.Mem.read(this.PC);
            this.PC++;
            this.lda(val);
        case 0xAA:  // TAX
            this.tax();    
        case 0xE8:  // INX
            this.inx();
        case 0x00:  // BRK
            return;
        }
    }
}

func (this *CPU) reset() {
    this.Reg_a = 0;
    this.Reg_x = 0;
    this.Reg_y = 0;
    this.Status = 0;
    this.PC = this.Mem.readU16(0xFFFC);
}

func (this *CPU) Load_and_run(program []byte) {
    this.Mem.load(program);
    this.reset();
    this.Run();
}

func (e *CPU) Scream() {
    fmt.Println("AAAAAAH, the CPU is burning!!!!");
}
