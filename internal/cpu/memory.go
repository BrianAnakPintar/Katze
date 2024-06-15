package cpu


type Memory struct {
    size uint32
    buf []byte
}

func MakeMemory(s uint32) *Memory {
    // For NES we need 65536 for s. Since we require 64KB
    mem := Memory{size: s, buf: make([]byte, s)};
    return &mem;
}

func (mem *Memory) read(addr uint16) byte {
    return mem.buf[addr];
}

func (mem *Memory) write(addr uint16, data byte) {
    mem.buf[addr] = data;
}

func (mem *Memory) readU16(addr uint16) uint16 {
    var lo byte = mem.read(addr);
    var hi byte = mem.read(addr+1);
    return (uint16(hi) << 8) | (uint16(lo))
}

func (mem *Memory) writeU16(addr uint16, data uint16) {
    mem.buf[addr] = byte(data & 0xFF)
    mem.buf[addr+1] = byte((data >> 8) & 0xFF)
}

func (mem *Memory) load(program []byte) {
    copy(mem.buf[0x8000:], program);
    mem.writeU16(0xFFFC, 0x8000)
}
