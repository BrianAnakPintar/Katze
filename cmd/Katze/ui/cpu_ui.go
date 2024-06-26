package ui

import (
	NESpkg "github.com/BrianAnakPintar/Katze/internal/emulator"
	rl "github.com/gen2brain/raylib-go/raylib"
    "fmt"
)

func toggleColor(cond bool) rl.Color {
    if cond {
        return rl.Lime
    }
    return rl.Maroon
}

func showFlags(x int32, y int32, cpu *NESpkg.CPU) {
    var padding int32 = 16

    rl.DrawText("Flags:", x, y, 16, rl.SkyBlue)
    rl.DrawText("N", x + padding*4, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_NEGATIVE)))
    rl.DrawText("V", x + padding*5, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_OVERFLOW)))
    rl.DrawText("-", x + padding*6, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_BREAK2)))
    rl.DrawText("B", x + padding*7, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_BREAK)))
    rl.DrawText("D", x + padding*8, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_DECIMAL)))
    rl.DrawText("I", x + padding*9, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_INTERRUPT)))
    rl.DrawText("Z", x + padding*10, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_ZERO)))
    rl.DrawText("C", x + padding*11, y, 16, toggleColor(cpu.ContainsFlag(NESpkg.FLAG_CARRY)))

}

func showRegisters(x int32, y int32, cpu *NESpkg.CPU) {
    xReg := fmt.Sprintf("X: 0x%x", cpu.X)
    yReg := fmt.Sprintf("Y: 0x%x", cpu.Y)
    acc := fmt.Sprintf("A: 0x%x", cpu.A)
    SP := fmt.Sprintf("SP: 0x%x", cpu.SP)
    PC := fmt.Sprintf("PC: 0x%x", cpu.PC)
    rl.DrawText(xReg, x, y, 16, rl.Black)
    rl.DrawText(yReg, x + 16 * 4, y, 16, rl.Black)

    rl.DrawText(acc, x, y + 16, 16, rl.Black)
    rl.DrawText(SP, x + 16 * 4, y + 16, 16, rl.Black)

    rl.DrawText(PC, x, y + 16 * 2, 16, rl.Black)
}

// Displays Debugging information for the CPU
func ShowCPU(x int32, y int32, cpu *NESpkg.CPU) {
    rl.DrawText("CPU INFORMATION", x, y, 18, rl.Blue)
    showFlags(x, y + 18, cpu)
    showRegisters(x, y + 18 * 2, cpu)
}
