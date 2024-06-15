package cpu

// Sets the flags for CPU based on the output of an operation (result)
func (this *CPU) SetFlags(result uint8) {
    // Sets the Z flag (ZERO FLAG)
    if result == 0 {
        this.Status = this.Status | FLAG_ZERO;
    } else {
        this.Status = this.Status & ^FLAG_ZERO; // Go does not have a bitwise NOT. So use XOR trick.
    }

    // Sets the N flag (NEGATIVE FLAG)
    if result & 0b1000_0000 != 0 {
        this.Status = this.Status | FLAG_NEGATIVE;
    } else {
        this.Status = this.Status & ^FLAG_NEGATIVE;
    }
}

func (this *CPU) lda(val uint8) {
    this.Reg_a = val;
    this.SetFlags(this.Reg_a);   
}

func (this *CPU) tax() {
    this.Reg_x = this.Reg_a;
    this.SetFlags(this.Reg_x);
}

func (this *CPU) inx() {
    this.Reg_x += 1;
    this.SetFlags(this.Reg_x); 
}
