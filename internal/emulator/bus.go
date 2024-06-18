package emulator

import (
    "sync"
)

var lock = &sync.Mutex{}

type BUS struct {
    cpu *CPU
    Mem []uint8
}

var BusInstance *BUS;

func GetBus() *BUS {
    if BusInstance == nil {
        lock.Lock()
        defer lock.Unlock()
        if BusInstance == nil {
            BusInstance = &BUS{cpu: nil, Mem: make([]uint8, 1024 * 64)}
        }
    }

    return BusInstance
}

func (bus *BUS) BusSetCPU(cpu *CPU) {
    bus.cpu = cpu;
}

func (bus *BUS) Write(addr uint16, val uint8) {
    if (addr >= 0 && addr <= 0xFFFF) {
        bus.Mem[addr] = val;
    }
}

func (bus *BUS) Read(addr uint16) uint8 {
    if (addr >= 0 && addr <= 0xFFFF) {
        return bus.Mem[addr];
    }   
    return 0x00;
}
