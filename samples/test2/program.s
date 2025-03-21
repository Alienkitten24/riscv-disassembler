# RISC-V Assembly example with conditions and loops:
# Calculates the sum of integers from 1 to n (n=10 in this example)

    .text
    .globl _start
_start:
    jal main

main:
    # Initialize variables
    addi t0, x0, 1        # t0 = counter = 1
    addi t1, x0, 10       # t1 = n = 10 (upper limit)
    addi t2, x0, 0        # t2 = sum = 0

loop:
    # Check if counter > n
    bgt t0, t1, exit_loop # If counter > n, exit loop
    
    # Add counter to sum
    add t2, t2, t0        # sum = sum + counter
    
    # Increment counter
    addi t0, t0, 1        # counter++
    
    # Loop back
    j loop                # Jump back to loop start

exit_loop:
    # Print result
    addi a7, x0, 1        # syscall 1 = print integer
    mv a0, t2             # move sum to a0
    ecall

    # Exit program
    addi a7, x0, 10       # syscall 10 = exit
    ecall
