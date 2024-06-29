package ui

import (
	NESpkg "github.com/BrianAnakPintar/Katze/internal/emulator"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func colorToRaylibColor(c NESpkg.Color) rl.Color {
    return rl.NewColor(c.R, c.G, c.B, 255)
}

func ShowPaletteAndPatternTable(x int32, y int32, ppu *NESpkg.PPU) {
    var size int32 = 6
    for p := int32(0); p < 8; p++ {
        for s := int32(0); s < 4; s++ {
            c := ppu.GetColorFromPalette(uint8(p), uint8(s))
            rl.DrawRectangle(x + p * (size + 5) + s * size, 340, size, size, colorToRaylibColor(c))
        }
    }    

    // var texture rl.Texture2D = rl.LoadTextureFromImage(ppu.GetPatternTable(0, 0))
    // var texture2 rl.Texture2D = rl.LoadTextureFromImage(ppu.GetPatternTable(1, 0))
    // rl.DrawTexture(texture, x, y, rl.White)
    // rl.DrawTexture(texture2, x+300, y, rl.White)
}
