package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sarchlab/onira/riscv"
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Ensure a file is specified
	if flag.NArg() < 1 {
		fmt.Println("Usage: disasm <elf_file>")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	// Open as an ELF file
	elfFile, err := elf.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening ELF file %s: %v\n", filename, err)
		os.Exit(1)
	}
	defer elfFile.Close()

	// Print the file format header information
	printFileHeader(filename, elfFile)

	// Process ELF file
	processElfFile(elfFile)
}

// printFileHeader prints the file format information in the expected format
func printFileHeader(filename string, elfFile *elf.File) {
	// Get just the base name of the file
	baseFileName := filepath.Base(filename)

	// Determine the machine architecture format string
	formatStr := "elf32-littleriscv"
	if elfFile.Class == elf.ELFCLASS64 {
		formatStr = "elf64-littleriscv"
	}

	// Print the header in the expected format
	fmt.Printf("\n%s:     file format %s\n\n\n", baseFileName, formatStr)
}

// processElfFile handles disassembly of an ELF file
func processElfFile(elfFile *elf.File) {
	// Load symbols from ELF file
	symbolTable := riscv.NewSymbolTable()
	symbolTable.LoadFromElf(elfFile)

	// Create decoder and printer
	decoder := riscv.NewInstDecoder()
	printer := riscv.NewInstPrinter(symbolTable)

	// Process each executable section
	for _, section := range elfFile.Sections {
		// Only process executable sections (code)
		if section.Flags&elf.SHF_EXECINSTR == 0 {
			continue
		}

		sectionData, err := section.Data()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading section %s: %v\n",
				section.Name, err)
			continue
		}

		fmt.Printf("Disassembly of section %s:\n", section.Name)

		// Disassemble the section
		disassembleSection(sectionData, uint32(section.Addr), decoder, printer)
	}
}

// disassembleSection disassembles a section of binary data
func disassembleSection(
	data []byte,
	startAddr uint32,
	decoder *riscv.InstDecoder,
	printer *riscv.InstPrinter,
) {
	for i := 0; i+3 < len(data); i += 4 {
		addr := startAddr + uint32(i)
		inst := decoder.Decode(addr, data[i:i+4])

		// Print the instruction with formatting
		fmt.Println(printer.FormatInstWithAddr(&inst))
	}
}
