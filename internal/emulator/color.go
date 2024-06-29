package emulator

type Color struct {
    R uint8
    G uint8 
    B uint8 
    A uint8 
}

func InitPaletteScreen() [0x40]Color {
    res := [0x40]Color{
        {84, 84, 84, 255}, // 0x00
        {0, 30, 116, 255}, // 0x01
        {8, 16, 144, 255}, // 0x02
        {48, 0, 136, 255}, // 0x03
        {68, 0, 100, 255}, // 0x04
        {92, 0, 48, 255}, // 0x05
        {84, 4, 0, 255}, // 0x06
        {60, 24, 0, 255}, // 0x07
        {32, 42, 0, 255}, // 0x08
        {8, 58, 0, 255}, // 0x09
        {0, 64, 0, 255}, // 0x0A
        {0, 60, 0, 255}, // 0x0B
        {0, 50, 60, 255}, // 0x0C
        {0, 0, 0, 255}, // 0x0D
        {0, 0, 0, 255}, // 0x0E
        {0, 0, 0, 255}, // 0x0F
        {152, 150, 152, 255}, // 0x10
        {8, 76, 196, 255}, // 0x11
        {48, 50, 236,255}, // 0x12
        {92, 30, 228,255}, // 0x13
        {136, 20, 176,255}, // 0x14
        {160, 20, 100,255}, // 0x15
        {152, 34, 32,255}, // 0x16
        {120, 60, 0, 255}, // 0x17
        {84, 90, 0, 255}, // 0x18
        {40, 114, 0, 255}, // 0x19
        {8, 124, 0, 255}, // 0x1A
        {0, 118, 40, 255}, // 0x1B
        {0, 102, 120, 255}, // 0x1C
        {0, 0, 0, 255}, // 0x1D
        {0, 0, 0, 255}, // 0x1E
        {0, 0, 0, 255}, // 0x1F
        {236, 238, 236, 255}, // 0x20
        {76, 154, 236, 255}, // 0x21
        {120, 124, 236, 255}, // 0x22
        {176, 98, 236, 255}, // 0x23
        {228, 84, 236, 255}, // 0x24
        {236, 88, 180, 255}, // 0x25
        {236, 106, 100, 255}, // 0x26
        {212, 136, 32, 255}, // 0x27
        {160, 170, 0, 255}, // 0x28
        {116, 196, 0, 255}, // 0x29
        {76, 208, 32, 255}, // 0x2A
        {56, 204, 108, 255}, // 0x2B
        {56, 180, 204, 255}, // 0x2C
        {60, 60, 60, 255}, // 0x2D
        {0, 0, 0, 255}, // 0x2E
        {0, 0, 0, 255}, // 0x2F
        {236, 238, 236, 255}, // 0x30
        {168, 204, 236, 255}, // 0x31
        {188, 188, 236, 255}, // 0x32
        {212, 178, 236, 255}, // 0x33
        {236, 174, 236, 255}, // 0x34
        {236, 174, 212, 255}, // 0x35
        {236, 180, 176, 255}, // 0x36
        {228, 196, 144, 255}, // 0x37
        {204, 210, 120, 255}, // 0x38
        {180, 222, 120, 255}, // 0x39
        {168, 226, 144, 255}, // 0x3A
        {152, 226, 180, 255}, // 0x3B
        {160, 214, 228, 255}, // 0x3C
        {160, 162, 160, 255}, // 0x3D
        {0, 0, 0, 255}, // 0x3E
        {0, 0, 0, 255}, // 0x3F
    }
    return res
}