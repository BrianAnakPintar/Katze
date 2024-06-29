package emulator

import (
    "sync"
)

var lock = &sync.Mutex{}

type BUS struct {
    cpu *CPU
    cpuRam []uint8

    ppu *PPU

    cartridge *Cartridge
    systemClockCounter uint32
}

var BusInstance *BUS;

func (this *BUS) GetCPU() *CPU {
    return this.cpu
}

func (this *BUS) GetPPU() *PPU {
    return this.ppu
}

func GetBus() *BUS {
    if BusInstance == nil {
        lock.Lock()
        defer lock.Unlock()
        if BusInstance == nil {
            BusInstance = &BUS{cpu: nil, cpuRam: make([]uint8, 1024 * 2), ppu: MakePPU()}
        }
    }

    return BusInstance
}

func (bus *BUS) BusSetCPU(cpu *CPU) {
    bus.cpu = cpu;
}

func (bus *BUS) CpuWrite(addr uint16, val uint8) {
    if bus.cartridge.cpuWrite(addr, val) {

    } else if (addr >= 0 && addr <= 0x1FFF) {
        bus.cpuRam[addr & 0x07FF] = val;
    } else if (addr >= 0x2000 && addr <= 0x3FFF) {
        bus.CpuWrite(addr, val)
    }
}

func (bus *BUS) CpuRead(addr uint16) uint8 {
    var data uint8 = 0
    if bus.cartridge.cpuRead(addr, &data) {
        // Cartridge Addr Range
    } else if (addr >= 0 && addr <= 0x1FFF) {
        return bus.cpuRam[addr & 0x07FF];
    } else if (addr >= 0x2000 && addr <= 0x3FFF) {
        return bus.ppu.cpuRead(addr & 0x0007)
    } 
    return data;
}

// The Actual System Interface | Since the Bus acts like the system itself. This will be our main Emulator "class" 

// Inserts a Cartridge into the NES
func (bus *BUS) InsertCartridge(cartridge *Cartridge) {
    bus.cartridge = cartridge
    bus.ppu.connectCartridge(cartridge)
}

// Reset button
func (bus *BUS) Reset() {
    bus.cpu.Reset()
    bus.systemClockCounter = 0
}

// Does a full tick
func (bus *BUS) Clock() {
    bus.ppu.clock();

    if bus.systemClockCounter % 3 == 0 {
        bus.cpu.Tick();
    }

    if bus.ppu.nmi {
        bus.ppu.nmi = false
        bus.cpu.nmi()
    }
    bus.systemClockCounter++
}
