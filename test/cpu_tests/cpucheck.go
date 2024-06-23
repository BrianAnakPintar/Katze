package main

import (
	"fmt"
    "os"
    "strings"
	c "github.com/BrianAnakPintar/Katze/internal/emulator"
)

func main() {
    var nes *c.BUS = c.GetBus();
    var cpu *c.CPU = c.MakeCPU();
    nes.BusSetCPU(cpu)

    var game *c.Cartridge = c.LoadCartridge("nestest.nes");
    if game == nil {
        fmt.Println("ERROR");
        return
    }
    nes.InsertCartridge(game);
    nes.Reset();

    // Open a file for logging
    logFile, err := os.OpenFile("cpu_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error opening log file:", err)
        return
    }
    defer logFile.Close()

    nes.GetCPU().Reset()
    nes.GetCPU().PC = 0xC000    // For some reason this is required lol

    for {
        if nes.GetCPU().InstructionFinished() {
            stats := trace(nes.GetCPU())
            if _, err := logFile.WriteString(stats + "\n"); err != nil {
                fmt.Println("Error writing to log file:", err)
            }
        }

        nes.GetCPU().Tick()

        if nes.GetCPU().PC == 0xC66E {
            break
        }
    }
}

func trace(cpu *c.CPU) string {
	code := cpu.Read(cpu.PC)
	ops, ok := c.Instructions[code]
	if !ok {
		panic(fmt.Sprintf("Unknown opcode: %02X", code))
	}

	begin := cpu.PC
	hex_dump := []byte{code}

	mem_addr, stored_value := uint16(0), byte(0)
	switch ops.Mode {
	case c.Immediate, c.Implied, c.Relative, c.Accumulator:
		// Do nothing
	default:
		addr := cpu.DebugGetAddress(ops.Mode)
		mem_addr = addr
		stored_value = cpu.Read(addr)
	}

	var tmp string
	switch ops.Size{
	case 1:
		switch ops.Code{
		case 0x0A, 0x4A, 0x2A, 0x6A:
			tmp = "A "
		default:
			tmp = ""
		}
	case 2:
		address := cpu.Read(begin + 1)
		hex_dump = append(hex_dump, address)

		switch ops.Mode {
		case c.Immediate:
			tmp = fmt.Sprintf("#$%02X", address)
		case c.ZeroPage:
			tmp = fmt.Sprintf("$%02X = %02X", mem_addr, stored_value)
		case c.ZeroPageX:
			tmp = fmt.Sprintf("$%02X,X @ %02X = %02X", address, mem_addr, stored_value)
		case c.ZeroPageY:
			tmp = fmt.Sprintf("$%02X,Y @ %02X = %02X", address, mem_addr, stored_value)
		case c.IndirectX:
			tmp = fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", address, uint16(address+cpu.X), mem_addr, stored_value)
		case c.IndirectY:
			tmp = fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", address, uint16(mem_addr-uint16(cpu.Y)), mem_addr, stored_value)
		default:
			addr := (begin + 2) + uint16(int8(address))
			tmp = fmt.Sprintf("$%04X", addr)
		}
	case 3:
		address_lo := cpu.Read(begin + 1)
		address_hi := cpu.Read(begin + 2)
		hex_dump = append(hex_dump, address_lo, address_hi)

		address := cpu.Read_u16(begin + 1)

		switch ops.Mode {
		case c.Absolute:
			tmp = fmt.Sprintf("$%04X = %02X", mem_addr, stored_value)
		case c.AbsoluteX:
			tmp = fmt.Sprintf("$%04X,X @ %04X = %02X", address, mem_addr, stored_value)
		case c.AbsoluteY:
			tmp = fmt.Sprintf("$%04X,Y @ %04X = %02X", address, mem_addr, stored_value)
		default:
            if ops.Code == 0x6C {
				jmp_addr := address
				if address&0x00FF == 0x00FF {
					lo := cpu.Read(address)
					hi := cpu.Read(address & 0xFF00)
					jmp_addr = uint16(hi)<<8 | uint16(lo)
				} else {
					jmp_addr = cpu.Read_u16(address)
				}
				tmp = fmt.Sprintf("($%04X) = %04X", address, jmp_addr)
			} else {
				tmp = fmt.Sprintf("$%04X", address)
			}
		}
	default:
		tmp = ""
	}

	hex_str := strings.ToUpper(strings.Join(hexStrings(hex_dump), " "))
	asm_str := fmt.Sprintf("%04X  %-8s %-4s %s", begin, hex_str, ops.Name, tmp)

	return fmt.Sprintf("%-47s A:%02X X:%02X Y:%02X P:%02X SP:%02X",
		asm_str, cpu.A, cpu.X, cpu.Y, cpu.STATUS, cpu.SP)
}

func hexStrings(bytes []byte) []string {
	hexes := make([]string, len(bytes))
	for i, b := range bytes {
		hexes[i] = fmt.Sprintf("%02X", b)
	}
	return hexes
}
