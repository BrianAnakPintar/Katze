package main

import (
	"fmt"
	c "github.com/BrianAnakPintar/Katze/internal/emulator"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var prog []uint8 = []uint8{0xA2, 0x0A, 0x8E, 0x00, 0x00, 0xA2, 0x03, 0x8E, 
                           0x01, 0x00, 0xAC, 0x00, 0x00, 0xA9, 0x00, 0x18, 
                           0x6D, 0x01, 0x00, 0x88, 0xD0, 0xFA, 0x8D, 0x02, 
                           0x00, 0xEA, 0xEA, 0xEA}

func main() {
    var cpu c.CPU = *c.MakeCPU();
    cpu.Load(prog);
    cpu.Reset()
    screenWidth := int32(1500)
	screenHeight := int32(800)
	rl.InitWindow(screenWidth, screenHeight, "CPU")

	defer rl.CloseWindow()
    memOffset := 0x0000
	memRowCount := 16 // Number of memory rows to display on the screen

	for !rl.WindowShouldClose() {
        handleInput(&cpu)
		// Handle input for scrolling
		if rl.IsKeyPressed(rl.KeyUp) && memOffset > 0 {
			memOffset -= memRowCount
		}
		if rl.IsKeyPressed(rl.KeyDown) && memOffset < (65536-memRowCount) {
			memOffset += memRowCount
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Display CPU details
		rl.DrawText("6502 CPU Emulator", 10, 10, 20, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("Accumulator (A): 0x%02X", cpu.A), 10, 50, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Index Register X (X): 0x%02X", cpu.X), 10, 80, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Index Register Y (Y): 0x%02X", cpu.Y), 10, 110, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Stack Pointer (SP): 0x%02X", cpu.SP), 10, 140, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Status Register (STATUS): 0b%08b", cpu.STATUS), 10, 170, 20, rl.Black)
		rl.DrawText(fmt.Sprintf("Program Counter (PC): 0x%04X", cpu.PC), 10, 200, 20, rl.Black)
        rl.DrawText(fmt.Sprintf("Instruction: ", (c.Instructions[cpu.Read(cpu.PC)]).Name), 10, 230, 20, rl.Blue)

		// Display memory content
		rl.DrawText("Memory", 500, 10, 20, rl.DarkGray)
		startY := int32(50)
		for i := 0; i < memRowCount; i++ {
			addr := memOffset + i*16
			memLine := fmt.Sprintf("0x%04X: ", addr)
			for j := 0; j < 16; j++ {
				if addr+j < 65536 {
					memLine += fmt.Sprintf("%02X ", cpu.Bus.Mem[addr+j])
				}
			}
			rl.DrawText(memLine, 500, startY+int32(i*20), 20, rl.Black)
		}

		rl.EndDrawing()
	}
}

func handleInput(cpu *c.CPU) {
    if rl.IsKeyPressed(rl.KeySpace) {
        for {
            cpu.Tick()
            if cpu.InstructionFinished() {
                break
            }
        }
    }
}
