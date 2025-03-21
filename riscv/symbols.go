package riscv

import (
	"debug/elf"
	"fmt"
)

// Symbol represents a symbol in the program
type Symbol struct {
	Name string
	Addr uint32
	Size uint32
	Type string // e.g., "FUNC", "OBJECT", etc.
}

// SymbolTable holds all the symbols in the program
type SymbolTable struct {
	Symbols      []Symbol
	AddrToSymbol map[uint32]*Symbol
}

// NewSymbolTable creates a new empty symbol table
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Symbols:      make([]Symbol, 0),
		AddrToSymbol: make(map[uint32]*Symbol),
	}
}

// AddSymbol adds a symbol to the table
func (st *SymbolTable) AddSymbol(
	name string, addr uint32, size uint32, symbolType string,
) {
	sym := Symbol{
		Name: name,
		Addr: addr,
		Size: size,
		Type: symbolType,
	}
	st.Symbols = append(st.Symbols, sym)

	// Store pointer to the symbol for quick address lookup
	st.AddrToSymbol[addr] = &st.Symbols[len(st.Symbols)-1]
}

// GetSymbolAtAddress returns the symbol at the given address, or nil if none
// exists
func (st *SymbolTable) GetSymbolAtAddress(addr uint32) *Symbol {
	if sym, exists := st.AddrToSymbol[addr]; exists {
		return sym
	}

	return nil
}

// GetSymbolForAddress returns the symbol name for a given address in
// jump/branch instructions It will append offset if the address is within a
// symbol but not at its start
func (st *SymbolTable) GetSymbolForAddress(addr uint32) string {
	// First check for exact match
	if sym, exists := st.AddrToSymbol[addr]; exists {
		return sym.Name
	}

	// Now check for an address within a function
	for _, sym := range st.Symbols {
		if addr >= sym.Addr && addr < sym.Addr+sym.Size {
			offset := addr - sym.Addr
			if offset == 0 {
				return sym.Name
			}

			return fmt.Sprintf("%s+%d", sym.Name, offset)
		}
	}

	// No symbol found, return empty string
	return ""
}

// LoadFromElf extracts symbols from an ELF file and adds them to the table
func (st *SymbolTable) LoadFromElf(elfFile *elf.File) {
	st.loadRegularSymbols(elfFile)
	st.loadDynamicSymbols(elfFile)
}

// loadRegularSymbols loads regular symbols from the ELF file
func (st *SymbolTable) loadRegularSymbols(elfFile *elf.File) {
	symbols, err := elfFile.Symbols()
	if err != nil {
		return
	}

	for _, sym := range symbols {
		if isValidSymbol(sym) {
			symbolType := determineSymbolType(elfFile, sym)
			size := determineSymbolSize(sym.Size)
			st.AddSymbol(sym.Name, uint32(sym.Value), size, symbolType)
		}
	}
}

// loadDynamicSymbols loads dynamic symbols from the ELF file
func (st *SymbolTable) loadDynamicSymbols(elfFile *elf.File) {
	dynSymbols, err := elfFile.DynamicSymbols()
	if err != nil {
		return
	}

	for _, sym := range dynSymbols {
		if isValidDynamicSymbol(sym) {
			size := determineSymbolSize(sym.Size)
			st.AddSymbol(sym.Name, uint32(sym.Value), size, "DYNSYM")
		}
	}
}

// isValidSymbol checks if a symbol should be included in the table
func isValidSymbol(sym elf.Symbol) bool {
	return sym.Value != 0 && sym.Name != ""
}

// isValidDynamicSymbol checks if a dynamic symbol should be included
func isValidDynamicSymbol(sym elf.Symbol) bool {
	return sym.Value != 0 && sym.Name != "" && sym.Info&0xf != 0
}

// determineSymbolSize returns an appropriate size for the symbol
func determineSymbolSize(size uint64) uint32 {
	if size == 0 {
		return 4 // Default to at least one instruction
	}

	return uint32(size)
}

// determineSymbolType examines an ELF symbol to determine its type
func determineSymbolType(elfFile *elf.File, sym elf.Symbol) string {
	// Check for symbol info type first
	infoType := getSymbolInfoType(sym)
	if infoType != "" {
		return infoType
	}

	// If not determined by info, try section-based identification
	return getSectionBasedType(elfFile, sym)
}

// getSymbolInfoType determines symbol type from its info field
func getSymbolInfoType(sym elf.Symbol) string {
	switch sym.Info & 0xf {
	case 0: // STT_NOTYPE
		return "NOTYPE"
	case 1: // STT_OBJECT
		return "OBJECT"
	case 2: // STT_FUNC
		return "FUNC"
	case 3: // STT_SECTION
		return "SECTION"
	case 4: // STT_FILE
		return "FILE"
	case 5: // STT_COMMON
		return "COMMON"
	case 6: // STT_TLS
		return "TLS"
	default:
		return "" // Not determined by info
	}
}

// getSectionBasedType determines type based on section properties
func getSectionBasedType(elfFile *elf.File, sym elf.Symbol) string {
	// Check special sections
	if sym.Section >= elf.SHN_LORESERVE {
		return getSpecialSectionType(sym.Section)
	}

	// Check regular sections
	if int(sym.Section) < len(elfFile.Sections) {
		section := elfFile.Sections[sym.Section]
		return getSectionFlagsType(section)
	}

	return "UNKNOWN"
}

// getSpecialSectionType handles special section index values
func getSpecialSectionType(section elf.SectionIndex) string {
	switch section {
	case elf.SHN_ABS:
		return "ABS"
	case elf.SHN_COMMON:
		return "COMMON"
	case elf.SHN_UNDEF:
		return "UNDEF"
	default:
		return "UNKNOWN"
	}
}

// getSectionFlagsType determines type based on section flags
func getSectionFlagsType(section *elf.Section) string {
	if section.Flags&elf.SHF_EXECINSTR != 0 {
		return "CODE"
	}

	if section.Flags&elf.SHF_WRITE != 0 {
		return "DATA"
	}

	return "UNKNOWN"
}
