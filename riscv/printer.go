package riscv

import (
	"fmt"
)

// InstPrinter formats and prints RISC-V instructions
type InstPrinter struct {
	SymbolTable *SymbolTable
}

// NewInstPrinter creates a new instruction printer
func NewInstPrinter(symbolTable *SymbolTable) *InstPrinter {
	return &InstPrinter{
		SymbolTable: symbolTable,
	}
}

// FormatInst formats an instruction into a string representation
func (p *InstPrinter) FormatInst(i *Inst) string {
	return ""
}

// FormatLabel returns a formatted symbol label string if there's a symbol at
// the instruction address
func (p *InstPrinter) FormatLabel(i *Inst) string {
	if p.SymbolTable == nil {
		return ""
	}

	symbol := p.SymbolTable.GetSymbolAtAddress(i.Addr)
	if symbol != nil {
		return fmt.Sprintf("\n%08x <%s>:", symbol.Addr, symbol.Name)
	}

	return ""
}

// FormatInstWithAddr formats an instruction with its address and binary
// representation
func (p *InstPrinter) FormatInstWithAddr(i *Inst) string {
	return p.FormatInst(i)
}
