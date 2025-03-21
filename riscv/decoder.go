package riscv

// InstDecoder handles decoding of RISC-V instructions
type InstDecoder struct {
	// Can store configuration options here if needed
	// For example, architecture variant (RV32I, RV64I, etc.)
}

// NewInstDecoder creates a new instruction decoder
func NewInstDecoder() *InstDecoder {
	return &InstDecoder{}
}

// Decode decodes a RISC-V instruction at the given address
func (d *InstDecoder) Decode(addr uint32, buf []byte) Inst {
	inst := Inst{Addr: addr}

	// TODO: Add your code here.

	return inst
}
