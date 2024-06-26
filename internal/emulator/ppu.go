package emulator

type (
    ControlFlag = uint8
    MaskFlag = uint8
    StatusFlag = uint8
)

// CTRL REGISTER VALUES
const (
    CTRLIncrementMode       ControlFlag = 1 << 2
    CTRLSpritePatternAddr   ControlFlag = 1 << 3
    CTRLBGPatternAddr       ControlFlag = 1 << 4
    CTRLSpriteSize          ControlFlag = 1 << 5
    CTRLSlaveMode           ControlFlag = 1 << 6
    CTRLNMI                 ControlFlag = 1 << 7
)

// MASK REGISTER VALUES
const (
    MaskGreyscale           MaskFlag = 1
    MaskRenderBGLeft        MaskFlag = 1 << 1
    MaskRenderSpritesLeft   MaskFlag = 1 << 2
    MaskRenderBG            MaskFlag = 1 << 3
    MaskRenderSprites       MaskFlag = 1 << 4
    MaskEmphRed             MaskFlag = 1 << 5
    MaskEmphGreen           MaskFlag = 1 << 6
    MaskEmphBlue            MaskFlag = 1 << 7
)

// STATUS REGISTER VALUES
const (
    StatusSpriteOverflow StatusFlag = 1 << 5
    StatusSpriteZeroHit  StatusFlag = 1 << 6
    StatusVerticalBlank  StatusFlag = 1 << 7
)

type PPU struct {
    cart *Cartridge
    CTRL ControlFlag     // 0x2000
    MASK MaskFlag        // 0x2001
    STATUS StatusFlag    // 0x2002
    OAMADDR uint8        // 0x2003
    OAMDATA uint8        // 0x2004
    SCROLL uint8         // 0x2005
    ADDR uint8           // 0x2006
    DATA uint8           // 0x2007
    nameTable [2][1024]byte
    patternTable [2][4096]byte
    paletteTable [32]byte


    scanline int16
    cycle int16
    nmi bool
    
    addr_latch bool
    ppu_data_buf uint8
    ppu_addr uint16

    //DEBUG PURPOSES
    FrameComplete bool
}


// CONTROL REGISTER

func (this *PPU) SetControlFlag(flag uint8, cond bool) {
    if cond {
        this.CTRL |= flag
    } else {
        this.CTRL &= ^flag
    }
}

func (this *PPU) ControlContainsFlag(flag uint8) bool {
    return this.CTRL & flag != 0
}

// END CONTROL REGISTER

// MASK REGISTER

func (this *PPU) SetMaskFlag(flag uint8, cond bool) {
    if cond {
        this.MASK |= flag
    } else {
        this.MASK &= ^flag
    }
}

func (this *PPU) MaskContainsFlag(flag uint8) bool {
    return this.MASK & flag != 0
}

// END MASK REGISTER


// STATUS REGISTER

func (this *PPU) SetStatusFlag(flag uint8, cond bool) {
    if cond {
        this.STATUS |= flag
    } else {
        this.STATUS &= ^flag
    }
}

func (this *PPU) StatusContainsFlag(flag uint8) bool {
    return this.STATUS & flag != 0
}

// END OF STATUS REGISTER FUNCTIONS


func (this *PPU) cpuWrite(addr uint16, data uint8) {
    switch addr {
    case 0x0000:    // Control
        this.CTRL = data
    case 0x0001:    // Mask
        this.MASK = data
    case 0x0002:    // Status
        
    case 0x0003:    // OAM Address
        
    case 0x0004:    // OAM Data
        
    case 0x0005:    // Scroll
        
    case 0x0006:    // PPU Address
        if !this.addr_latch {

            this.ppu_addr = (this.ppu_addr & 0x00FF) | (uint16(data) << 8)
            this.addr_latch = true
        } else {
            this.ppu_addr = (this.ppu_addr & 0xFF00) | uint16(data)
            this.addr_latch = false
        }
    case 0x0007:    // PPU Data
        this.ppuWrite(this.ppu_addr, data)
        this.ppu_addr++
    }
}

func (this *PPU) cpuRead(addr uint16) uint8 {
    var data uint8 = 0
    switch addr {
    case 0x0000:    // Control
        
    case 0x0001:    // Mask
        
    case 0x0002:    // Status
        data = (this.STATUS & 0xE0) | (this.ppu_data_buf & 0x1F);
        this.SetStatusFlag(StatusVerticalBlank, false);
        this.addr_latch = false
    case 0x0003:    // OAM Address
        
    case 0x0004:    // OAM Data
        
    case 0x0005:    // Scroll
        
    case 0x0006:    // PPU Address
        
    case 0x0007:    // PPU Data
        data = this.ppu_data_buf;
        this.ppu_data_buf = this.ppuRead(this.ppu_addr);

        if (this.ppu_addr > 0x3F00) {
            data = this.ppu_data_buf
        }
        this.ppu_addr++
    }
    return data
}

func (this *PPU) ppuWrite(addr uint16, data uint8) {
    addr &= 0x3FFF;
    if this.cart.ppuWrite(addr, data) {

    } else if (addr >= 0 && addr <= 0x1FFF) {
        this.patternTable[(addr & 0x1000) >> 12][addr & 0x0FFF] = data
    } else if (addr >= 0x2000 && addr <= 0x3EFF) {
        addr &= 0x0FFF
        if (this.cart.MirrorMode == MirrorVertical) {
            // Vertical
		    if (addr >= 0x0000 && addr <= 0x03FF) {
				this.nameTable[0][addr & 0x03FF] = data;
            }
			if (addr >= 0x0400 && addr <= 0x07FF) {
				this.nameTable[1][addr & 0x03FF] = data;
            }
			if (addr >= 0x0800 && addr <= 0x0BFF) {
				this.nameTable[0][addr & 0x03FF] = data;
            }
			if (addr >= 0x0C00 && addr <= 0x0FFF) {
				this.nameTable[1][addr & 0x03FF] = data;
            }
        } else if (this.cart.MirrorMode == MirrorHorizontal) {
            // Horizontal 
			if (addr >= 0x0000 && addr <= 0x03FF) {
				this.nameTable[0][addr & 0x03FF] = data;
            }
			if (addr >= 0x0400 && addr <= 0x07FF) {
				this.nameTable[0][addr & 0x03FF] = data;
            }
			if (addr >= 0x0800 && addr <= 0x0BFF) {
				this.nameTable[1][addr & 0x03FF] = data;
            }
			if (addr >= 0x0C00 && addr <= 0x0FFF) {
				this.nameTable[1][addr & 0x03FF] = data;
            }
        }
    } else if (addr >= 0x3F00 && addr <= 0x3FFF) {
        addr &= 0x001F
        if addr == 0x0010 {
            addr = 0x0000
        } else if addr == 0x0014 {
            addr = 0x0004
        } else if addr == 0x0018 {
            addr = 0x0008
        } else if addr == 0x001C {
            addr = 0x000C
        }
        this.paletteTable[addr] = data
    }
}

func (this *PPU) ppuRead(addr uint16) uint8 {
    var data uint8 = 0
    addr &= 0x3FFF;
    
    if this.cart.ppuRead(addr, &data) {

    } else if (addr >= 0 && addr <= 0x1FFF) {
        data = this.patternTable[(addr & 0x1000) >> 12][addr & 0x0FFF]
    } else if (addr >= 0x2000 && addr <= 0x3FFF) {
        addr &= 0x0FFF
        if (this.cart.MirrorMode == MirrorVertical) {
            // Vertical
			if (addr >= 0x0000 && addr <= 0x03FF) {
				data = this.nameTable[0][addr & 0x03FF];
            }
			if (addr >= 0x0400 && addr <= 0x07FF) {
				data = this.nameTable[1][addr & 0x03FF];
            }
			if (addr >= 0x0800 && addr <= 0x0BFF) {
				data = this.nameTable[0][addr & 0x03FF];
            }
			if (addr >= 0x0C00 && addr <= 0x0FFF) {
				data = this.nameTable[1][addr & 0x03FF];
            }
        } else if (this.cart.MirrorMode == MirrorHorizontal) {
            // Horizontal 
			if (addr >= 0x0000 && addr <= 0x03FF) {
				data = this.nameTable[0][addr & 0x03FF];
            }
			if (addr >= 0x0400 && addr <= 0x07FF) {
				data = this.nameTable[0][addr & 0x03FF];
            }
			if (addr >= 0x0800 && addr <= 0x0BFF) {
				data = this.nameTable[1][addr & 0x03FF];
            }
			if (addr >= 0x0C00 && addr <= 0x0FFF) {
				data = this.nameTable[1][addr & 0x03FF];
            }
        }
    } else if (addr >= 0x3F00 && addr <= 0x3FFF) {
        addr &= 0x001F
        if addr == 0x0010 {
            addr = 0x0000
        } else if addr == 0x0014 {
            addr = 0x0004
        } else if addr == 0x0018 {
            addr = 0x0008
        } else if addr == 0x001C {
            addr = 0x000C
        }

        if this.MaskContainsFlag(MaskGreyscale) {
            data = this.paletteTable[addr] & 0x30
        } else {
            data = this.paletteTable[addr] & 0x3F
        }
    }

    return data
}

func (this *PPU) connectCartridge(c *Cartridge) {
    this.cart = c
}

func (this *PPU) clock() {
    if (this.scanline == -1 && this.cycle == 1) {
        this.SetStatusFlag(StatusVerticalBlank, false)
    }

    if (this.scanline == 241 && this.cycle == 1) {
        this.SetStatusFlag(StatusVerticalBlank, true)
        if this.ControlContainsFlag(CTRLNMI) {
            this.nmi = true
        }

    }

    this.cycle++;
    if this.cycle >= 341 {
        this.cycle = 0;
        this.scanline++;

        if this.scanline >= 261 {
            this.scanline = -1
            this.FrameComplete = true
        }
    }
}
