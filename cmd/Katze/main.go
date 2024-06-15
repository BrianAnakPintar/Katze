package main

import (
    "github.com/BrianAnakPintar/Katze/internal/cpu"
)

func main() {
    var main_cpu *cpu.CPU = cpu.MakeCPU();
    main_cpu.Scream();
    var instr []uint8 = []uint8{0xA9, 0x05, 0x00, 0x00};
    main_cpu.Interpret(instr);
}
