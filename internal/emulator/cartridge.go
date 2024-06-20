package emulator

import (
	"hash/crc32"
	"io"
	"os"
)

type MIRROR uint8;

const (
    MirrorHorizontal MIRROR = 0
    MirrorVertical MIRROR = 1
    MirrorSingle0 MIRROR = 2
    MirrorSingle1 MIRROR = 3
)

type Cartridge struct {
    MapperID uint8
    CHRBank uint8
    PRGBank uint8
    PRGMemory []uint8
    CHRMemory []uint8
    CRC32 uint32
    MirrorMode MIRROR 
    CHRRam bool
    mapper Mapper
}


func LoadCartridge(path string) *Cartridge {
    file, err := os.Open(path)
    if err != nil {
        // Uh oh.
        return nil
    }
    defer file.Close()
    header := make([]byte, 16)
    x, err := file.Read(header)

    if err != nil || x > 16{
        // Uh oh
        return nil
    }
    // Check header signature.
	if header[0] != 'N' || header[1] != 'E' || header[2] != 'S' || header[3] != 0x1A {
		return nil
	}

    var (
		mapperID   = (header[6] >> 4) | (header[7] & (1 << 4))
		prgBanks   = header[4]
		chrBanks   = header[5]
		hasTrainer = header[6]&(0x04) != 0
		// hasBattery = header[6]&(1<<1) != 0
		mirrorMode = header[6] & (1 << 0)
	)
    // If there's training info. Skip it (512 bytes)
    if hasTrainer {
        file.Seek(512, io.SeekCurrent)
    }
    
    var mapper Mapper
    switch mapperID {
    case 0:
        mapper = &Mapper0{prgBanks: prgBanks, chrBanks: chrBanks}
    }

    var ines_file_type uint8 = 1
    if ines_file_type == 0 {
        // TODO 
    } else if ines_file_type == 1 {
        h := crc32.NewIEEE()
        romReader := io.TeeReader(file, h)
	    prgData := make([]byte, int(prgBanks) * 1024 * 16)
        chrData := make([]byte, int(chrBanks) * 1024 * 8)

        romReader.Read(prgData)
        romReader.Read(chrData)
        // If no CHR Data then set 8KB for CHR RAM
        var chrRAM bool
        if len(chrData) == 0 {
            chrData = make([]byte, 1024 * 8)
            chrRAM = true
        }
        return &Cartridge{
                MapperID: mapperID,
                CHRBank: chrBanks, 
                PRGBank: prgBanks, 
                PRGMemory: prgData, 
                CHRMemory: chrData, 
                CRC32: h.Sum32(), 
                MirrorMode: MIRROR(mirrorMode),
                CHRRam: chrRAM,
                mapper: mapper}
    } else if ines_file_type == 2 {
        // TODO
    }

    return nil
}

func (this *Cartridge) cpuWrite(addr uint16, data uint8) bool {
    var mapped_addr uint32 = 0
    if this.mapper.cpuMapRead(addr, &mapped_addr) {
        this.PRGMemory[mapped_addr] = data
        return true
    }
    return false
}

func (this *Cartridge) cpuRead(addr uint16, buf *uint8) bool {
    var mapped_addr uint32 = 0
    if this.mapper.cpuMapRead(addr, &mapped_addr) {
        data := this.PRGMemory[mapped_addr]
        *buf = data
        return true
    }
    return false
}

func (this *Cartridge) ppuWrite(addr uint16, data uint8) bool {
    var mapped_addr uint32 = 0
    if this.mapper.cpuMapRead(addr, &mapped_addr) {
        this.CHRMemory[mapped_addr] = data
        return true
    }
   return false
}

func (this *Cartridge) ppuRead(addr uint16, buf *uint8) bool {
    var mapped_addr uint32 = 0
    if this.mapper.cpuMapRead(addr, &mapped_addr) {
        data := this.CHRMemory[mapped_addr]
        *buf = data
        return true
    }
    return false
}
