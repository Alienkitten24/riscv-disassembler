# Assignment 4

This assignment is given to ensure you to understand the organization of the RISC-V ISA. You are expected to implement a simple RISC-V disassembler. 

The main function is given in `cmd/main.go` file. Currently, it only prints empty strings. So you will need to implement the `decoder.Decode` and `printer.FormatInstWithAddr` methods. 

The `Inst` struct is given in `riscv/inst.go` file. The goal is to fill this structs with the `decoder.Decode` method. Then, you will print the content of this struct with the `printer.FormatInstWithAddr` method.

Two examples are provided in `samples` folder. The `program.elf` file will be your input. Your final output should be the same as in the `program.disasm` file. 

I will not divide this assignment problem into problems. But here is a grading schemes that I recommend you follow the order to implement.

1. (25pt) Print all the mnemonics of the instructions. 
2. (25pt) Print all the mnemonics with all the register operands.
3. (25pt) Print all the mnemonics with all the operands, with an example of label operands.
4. (25pt) Print all the mnemonics with all the operands, including labels. Print all the labels in the file.
5. (25pt bonus) Pass both the tests in `.github/workflows/go-build.yml` file. That requires your output to be exactly the same as the `program.disasm` file. You may need to modify the `FormatInstWithAddr` of the `InstPrinter` struct.