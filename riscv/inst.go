package riscv

// Inst represents a single decoded RISC-V instruction.
type Inst struct {
	Addr uint32
	Bin  uint32
	Op   string
	Rd   string
	Rs1  string
	Rs2  string
	Imm  int32
}
