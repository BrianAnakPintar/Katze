package emulator

type PPU struct {
    cart *Cartridge
    nameTable [2][1024]byte
    paletteTable [32]byte
}

func (this *PPU) cpuWrite(addr uint16, data uint8) {
    switch addr {
    case 0x0000:    // Control
        
    case 0x0001:    // Mask
        
    case 0x0002:    // Status
        
    case 0x0003:    // OAM Address
        
    case 0x0004:    // OAM Data
        
    case 0x0005:    // Scroll
        
    case 0x0006:    // PPU Address
        
    case 0x0007:    // PPU Data
        
    }
}

func (this *PPU) cpuRead(addr uint16) uint8 {
    var data uint8 = 0
    switch addr {
    case 0x0000:    // Control
        
    case 0x0001:    // Mask
        
    case 0x0002:    // Status
        
    case 0x0003:    // OAM Address
        
    case 0x0004:    // OAM Data
        
    case 0x0005:    // Scroll
        
    case 0x0006:    // PPU Address
        
    case 0x0007:    // PPU Data
        
    }
    return data
}

func (this *PPU) ppuWrite(addr uint16, data uint8) {
    addr &= 0x3FFF;
    if this.cart.ppuWrite(addr, data) {

    }
}

func (this *PPU) ppuRead(addr uint16) uint8 {
    var data uint8 = 0
    addr &= 0x3FFF;
    
    if this.cart.ppuRead(addr, &data) {

    }

    return data
}

func (this *PPU) connectCartridge(c *Cartridge) {
    this.cart = c
}
