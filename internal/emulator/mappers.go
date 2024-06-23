package emulator

/*
Mappers are not actually writing data. Instead they work similar to that of Virtual Memory
They basically take in an address and return the physical address. Pretty cool!
Additionally There are various types of Mapper that's supported by the NES, each with its own
quirks on how they translate the addresses.
*/

type Mapper interface {
    cpuMapRead(addr uint16, mapped_addr *uint32) bool;
    cpuMapWrite(addr uint16, mapped_addr *uint32) bool;
    ppuMapRead(addr uint16, mapped_addr *uint32) bool;
    ppuMapWrite(addr uint16, mapped_addr *uint32) bool;
}

// Subclass
type Mapper0 struct {
    prgBanks uint8
    chrBanks uint8
}

func (m *Mapper0) cpuMapRead(addr uint16, mapped_addr *uint32) bool {
    if (addr >= 0x8000 && addr <= 0xFFFF) {
        if m.prgBanks > 1 {
            *mapped_addr = uint32(addr) & 0x7FFF
        } else {
            *mapped_addr = uint32(addr) & 0x3FFF
        }
        return true
    }

    return false
}

func (m *Mapper0) cpuMapWrite(addr uint16, mapped_addr *uint32) bool {
    if (addr >= 0x8000 && addr <= 0xFFFF) {
        if m.prgBanks > 1 {
            *mapped_addr = uint32(addr) & 0x7FFF
        } else {
            *mapped_addr = uint32(addr) & 0x3FFF
        }
        return true
    }

    return false
}

func (m *Mapper0) ppuMapRead(addr uint16, mapped_addr *uint32) bool {
    if (addr >= 0x0000 && addr <= 0x1FFF) {
        *mapped_addr = uint32(addr)
         return true
    }
    return false
}

func (m *Mapper0) ppuMapWrite(addr uint16, mapped_addr *uint32) bool {
    return false
}
