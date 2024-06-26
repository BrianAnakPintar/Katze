package main

import (
	// "fmt"
	NESpkg "github.com/BrianAnakPintar/Katze/internal/emulator"
	rl "github.com/gen2brain/raylib-go/raylib"
    ui "github.com/BrianAnakPintar/Katze/cmd/Katze/ui"
)

func main() {
    var nes *NESpkg.BUS = NESpkg.GetBus();
    var cpu *NESpkg.CPU = NESpkg.MakeCPU();
    nes.BusSetCPU(cpu)

    var game *NESpkg.Cartridge = NESpkg.LoadCartridge("nestest.nes");
    nes.InsertCartridge(game);
    nes.Reset();

    var screenWidth int32 = 256 * 3
    var screenHeight int32 = 240 * 3
    rl.InitWindow(screenWidth, screenHeight, "Katze")
    defer rl.CloseWindow()

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)
        ui.ShowCPU(screenWidth * 2/3, 10, cpu)
        rl.DrawText("KATZE", screenWidth/2, screenHeight/2, 20, rl.Black)
        rl.EndDrawing()
    }
}
