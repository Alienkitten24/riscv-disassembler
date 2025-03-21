# Makefile for RISC-V assembly programs

# Tools
RISCV_PREFIX ?= riscv64-unknown-elf-
GCC = $(RISCV_PREFIX)gcc
LD = $(RISCV_PREFIX)ld
OBJDUMP = $(RISCV_PREFIX)objdump

# Flags
GCCFLAGS = -march=rv32i -mabi=ilp32 -nostdlib
LDFLAGS = -static

# Files
SRC = program.s
ELF = program.elf
DISASM = program.disasm

.PHONY: all clean

all: $(ELF) $(DISASM)

$(ELF): $(SRC)
	$(GCC) $(GCCFLAGS) -o $@ $<

$(DISASM): $(ELF)
	$(OBJDUMP) -d -Mno-aliases $< > $@

clean:
	rm -f $(ELF) $(DISASM)
