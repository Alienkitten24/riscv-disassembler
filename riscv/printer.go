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
	
	offset := uint32(int32(i.Addr) + i.Imm)
	// fmt.Println("hola")
	// fmt.Printf("%x\n", i.Addr)
	// fmt.Printf("%x\n", int32(i.Addr))
	// fmt.Printf("%x\n", i.Imm)
	// fmt.Printf("%x\n", offset)
	symbol := p.SymbolTable.GetSymbolAtAddress(offset)

	switch i.Opcode {
	case "0110111": // lui
		return fmt.Sprintf("%s\t%s,%d", i.Op, i.Rd, i.Imm)
	case "0010111":  // auipc
		return fmt.Sprintf("%s\t%s,%x", i.Op, i.Rd, offset)
	case "1101111":  // jal
		return fmt.Sprintf("%s\t%s,%x <%s>", i.Op, i.Rd, offset, symbol.Name)
	case "1100111":  // jalr
	case "1100011":  // branches
		return fmt.Sprintf("%s\t%s,%s,%x <%s>", i.Op, i.Rs1, i.Rs2, offset, symbol.Name)
	case "0000011":  // loads
		return fmt.Sprintf("%s\t%s, %x(%s)", i.Op, i.Rd, offset, i.Rs1)
	case "0100011":  // stores
		return fmt.Sprintf("%s\t%s, %x(%s)", i.Op, i.Rs2, offset, i.Rs1)
	case "0010011":  // imm arithmatic
		return fmt.Sprintf("%s\t%s,%s,%d", i.Op, i.Rd, i.Rs1, i.Imm)
	case "0110011":  // reg arithmatic
		return fmt.Sprintf("%s\t%s,%s,%s", i.Op, i.Rd, i.Rs1, i.Rs2)
	case "0001111":  
	case "1110011":  // sys calls
		return fmt.Sprintf("%s", i.Op)
	}

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
	labelStr := p.FormatLabel(i)
	if labelStr != "" {
		fmt.Println(labelStr)
	}
	return fmt.Sprintf("   %x:\t%08x          \t%s", i.Addr, i.Bin, p.FormatInst(i))
}