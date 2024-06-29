package main

import (
	"fmt"
	ui "github.com/BrianAnakPintar/Katze/cmd/Katze/ui"
	NESpkg "github.com/BrianAnakPintar/Katze/internal/emulator"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func HandleInput(nes *NESpkg.BUS) {
    if rl.IsKeyDown(rl.KeyC) {
        for {
            nes.Clock()
            if nes.GetCPU().InstructionFinished() {
                break
            }
        }
        for {
            nes.Clock()
            if !nes.GetCPU().InstructionFinished() {
                break
            }
        }
    }
    // Clock until frame complete
    if rl.IsKeyPressed(rl.KeyF) {
        for {
            nes.Clock()
            if nes.GetPPU().FrameComplete {
                break
            }
        }

        for {
            nes.Clock()
            if nes.GetCPU().InstructionFinished() {
                break
            }
        }

        nes.GetPPU().FrameComplete = false
    }
}

func main() {
    var nes *NESpkg.BUS = NESpkg.GetBus();
    var cpu *NESpkg.CPU = NESpkg.MakeCPU();
    nes.BusSetCPU(cpu)

    var game *NESpkg.Cartridge = NESpkg.LoadCartridge("nestest.nes");
    if game == nil {
        fmt.Println("ERROR");
        return
    }
    nes.InsertCartridge(game);
    
    nes.Reset();    
    nes.GetCPU().PC = 0xC000    // For some reason this is required lol
    var screenWidth int32 = 256 * 3
    var screenHeight int32 = 240 * 3
    rl.InitWindow(screenWidth, screenHeight, "Katze")
    defer rl.CloseWindow()

    for !rl.WindowShouldClose() {
        HandleInput(nes)
        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)
        ui.ShowCPU(screenWidth * 2/3, 10, cpu)
        ui.ShowPaletteAndPatternTable(100, 100, nes.GetPPU())
        rl.DrawText("KATZE", screenWidth/2, screenHeight/2, 20, rl.Black)
        rl.EndDrawing()
    }
}
