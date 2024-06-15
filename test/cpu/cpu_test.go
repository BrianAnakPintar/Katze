package cpu

import (
	"testing"

	"github.com/BrianAnakPintar/Katze/internal/cpu"
)

func TestLDA(t *testing.T) {
    var tests = []struct { input []uint8; expected uint8; z_flag uint8; n_flag uint8} {
        {[]uint8{0xA9, 0x05, 0x00}, 0x05, 0, 0},
        {[]uint8{0xa9, 0x05, 0xA9, 0xFF, 0x00}, 0xFF, 0, cpu.FLAG_NEGATIVE},
    }

    var t_cpu *cpu.CPU = cpu.MakeCPU();
    var instr []uint8 = []uint8{0xA9, 0x05, 0x00};
    for i, test := range tests {
        t_cpu.Load_and_run(test.input)
        if t_cpu.Reg_a != test.expected {
            t.Error("Expected: ", t_cpu.Reg_a, "Received: ", t_cpu.Reg_a);
        }
        if t_cpu.Status & cpu.FLAG_ZERO != test.z_flag {
            // Pass
            t.Error("Incorrect Z Flag, ", i)
        }
        if t_cpu.Status & cpu.FLAG_NEGATIVE != test.n_flag {
            t.Error("Incorrect N Flag", i)
        }
    }
    t_cpu.Load_and_run(instr);
}
