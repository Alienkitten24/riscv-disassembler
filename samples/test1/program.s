# Simplest RISC-V Assembly example: Compute 1 + 1 and print

    .text
    .globl _start
_start:
    jal main

main:
    addi t0, x0, 1        # t0 = 1
    addi t1, x0, 1        # t1 = 1
    add t2, t0, t1        # t2 = t0 + t1 (1 + 1)

    # Use ecall to print the integer result (assume emulator-defined)
    addi a7, x0, 1        # syscall 1 = print integer
    mv a0, t2             # move result to a0
    ecall

    # Exit program using ecall (emulator-defined)
    addi a7, x0, 10       # syscall 10 = exit
    ecall
