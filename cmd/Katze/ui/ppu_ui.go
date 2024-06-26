package ui

import (
	NESpkg "github.com/BrianAnakPintar/Katze/internal/emulator"
	rl "github.com/gen2brain/raylib-go/raylib"
    "fmt"
)

func ShowPaletteAndPatternTable(x int32, y int32) {
    var size int32 = 6
    for p := int32(0); p < 8; p++ {
        for s := int32(0); s < 4; s++ {
            rl.DrawRectangle(x + p * (size + 5) + s * size, 340, size, size, rl.Black)
        }
    }    
}
