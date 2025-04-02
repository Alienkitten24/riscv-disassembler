package riscv
import (
	// "fmt"
	"strconv"
	// "reflect"
)

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
	// inst := Inst{Addr: addr}

	// TODO: Add your code here.
	// fmt.Println(addr)
	// for _, element := range buf {
	// 	fmt.Printf("%x ", element)
	// }
	// fmt.Print()


	bits := bytesToBits(buf)

	// fmt.Println(readBitRange(0,31,bits))
	
	var destinationReg, funct3, funct7, sourceReg1, sourceReg2 string
	var immediate int32

	opcode := readBitRange(0,6,bits)
	// funct3 = readBitRange(12,14,bits)
	// funct7 = readBitRange(25,31,bits)
	funct3 = "" 
	funct7 = ""

	// instId := instructionIdentifier{opcode, funct3, funct7} 
	// instData := instructionMap[instId]	

	// fmt.Println("LOW TEST")
	// fmt.Println(instId)
	// fmt.Println(instData)

	instType := instructionTypeMap[opcode]

	switch instType{
	case "R":
		destinationReg = readBitRange(7,11,bits)
		funct3 = readBitRange(12,14,bits)
		sourceReg1 = readBitRange(15,19,bits)
		sourceReg2 = readBitRange(20,24,bits)
		funct7 = readBitRange(25,31,bits)
	case "I":
		destinationReg = readBitRange(7,11,bits)
		funct3 = readBitRange(12,14,bits)
		sourceReg1 = readBitRange(15,19,bits)
		immNum, err := strconv.ParseInt(readBitRange(20,31,bits), 2, 32)
		if err != nil {
			panic(err)
		}
		immediate = convert2sCompliment(immNum)
	case "S":
		funct3 = readBitRange(12,14,bits)
		sourceReg1 = readBitRange(15,19,bits)
		sourceReg2 = readBitRange(20,24,bits)
		immBitStr := readBitRange(5,11,bits) + readBitRange(0,4,bits) 
		immNum, err := strconv.ParseInt(immBitStr, 2, 32)
		if err != nil {
			panic(err)
		}
		immediate = convert2sCompliment(immNum)
	case "B":
		funct3 = readBitRange(12,14,bits)
		sourceReg1 = readBitRange(15,19,bits)
		sourceReg2 = readBitRange(20,24,bits)
		middleBits := bits[32-12:32-7] // grab the 'middle' 5 bits of the bit string
		endBits := bits[0:32-25] // grab the 'last' 7 bits of the bit string
		immBitStr := readBitRange(6,6,endBits) + readBitRange(0,0,middleBits) + readBitRange(0,5,endBits) + readBitRange(1,4,middleBits) + "0"
		// fmt.Println("hello")
		// fmt.Println(middleBits)
		// fmt.Println(endBits)
		// fmt.Println(immBitStr)
		immNum, err := strconv.ParseInt(immBitStr, 2, 32)
		if err != nil {
			panic(err)
		}
		immediate = convert2sCompliment(immNum)
		// fmt.Println(immediate)
	case "U":
		destinationReg = readBitRange(7,11,bits)
		immNum, err := strconv.ParseInt(readBitRange(12,31,bits), 2, 32)
		if err != nil {
			panic(err)
		}
		immediate = convert2sCompliment(immNum)
	case "J":
		destinationReg = readBitRange(7,11,bits)
		// immBitStr := readBitRange(12,31,bits)
		// immNum, err := strconv.ParseInt(immBitStr, 2, 32)
		endBits := bits[0:32-12] // grab the 'last' 20 bits of the bit string
		immBitStr := readBitRange(19,19,endBits) + readBitRange(0,7,endBits) + readBitRange(8,8,endBits) + readBitRange(9,18,endBits) + "0"
		// fmt.Print("J hello")
		// fmt.Println(endBits)
		immNum, err := strconv.ParseInt(immBitStr, 2, 32)
		// fmt.Println(reflect.TypeOf(immNum))
		// fmt.Println(immNum)
		if err != nil {
			panic(err)
		}
		immediate = convert2sCompliment(immNum)
		// fmt.Printf("ULTRA TEST %s = %d \n", immBitStr, immediate)
	}

	instId := instructionIdentifier{opcode, funct3, funct7} 
	instData := instructionMap[instId]	

	// fmt.Println("LOW TEST")
	// fmt.Println(instId)
	// fmt.Println(instData)
	// fmt.Println("lowtest" + funct3 + funct7 + destinationReg)

	inst := Inst{Addr: addr, Bin: bitsToUInt(bits), Op: instData.mnemonic, Rd: registerMap[bitStringToInt(destinationReg)], Rs1: registerMap[bitStringToInt(sourceReg1)], Rs2: registerMap[bitStringToInt(sourceReg2)], Imm: immediate, Opcode: opcode} 

	// fmt.Println(inst)
	// fmt.Println(instData)

	// fmt.Println(registerMap[5])

	return inst
}

func bytesToBits(data []byte) []int {
	bits := make([]int, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			bits[(len(data)-1-i)*8+7-j] = int((b >> j) & 1) // Preserve bit order but reverse byte order
		}
	}
	return bits
}

func readBitRange(start int, end int, bits []int) string {
	subBits := bits[len(bits)-end-1:len(bits)-start]	

	strBits := make([]byte, len(subBits))	
	for i, bit := range subBits {
		if bit == 1 {
			strBits[i] = '1'
		} else {
			strBits[i] = '0'
		}
	}

	return string(strBits)
}

// convert bit string to int
func bitStringToInt(bits string) int {
	if bits == "" {
		return 0 // TODO REMOVE THIs
	}
	num, err := strconv.ParseInt(bits, 2, 0)
	if err != nil {
		panic(err)
	} else {
		return int(num)
	}
}

// convert bit array to uint32
func bitsToUInt(bits []int) uint32 {
	var result uint32
	result = 0
	for _, bit := range bits {
		result = (result << 1) | uint32(bit) // Shift left and add the bit
	}
	return result
}

func convert2sCompliment(val int64) int32 {
	if val >= (1 << 20) { // If the 21st bit is set (negative number in two's complement)
		val -= (1 << 21) // Subtract 2^21 to get the correct signed value
	}
	return int32(val)
}

// An instruction can be identified by its opcode, funct3, and funct7 (if present)
// for exmple, an addi inst is uniquely identified by opcode=0010011 and funct3=000
type instructionIdentifier struct {
	opcode string
	funct3 string
	funct7 string
}

// An instruction has a name and type
// for example, the name addi and the type I
type instructionData struct {
	mnemonic string
	instructionType string
}

// TODO FINISH WITH FENCE AND OTHERS
var instructionTypeMap = map[string]string {
	"0110111":  "U", // lui
	"0010111":  "U", // auipc
	"1101111":  "J", // jal
	"1100111":  "I", // jalr
	"1100011":  "B", // branches
	"0000011":  "I", // loads
	"0100011":  "S", // stores
	"0010011":  "I", // imm arithmatic
	"0110011":  "R", // reg arithmatic
	"0001111":  "", 
	"1110011":  "I", // sys calls
}

var instructionMap = map[instructionIdentifier]instructionData {
	{"0110111", "", ""}: {"lui", "U"},
	{"0010111", "", ""}: {"auipc", "U"},
	{"1101111", "", ""}: {"jal", "J"},
	{"1100111", "000", ""}: {"jalr", "I"},
	{"1100011", "000", ""}: {"beq", "B"},
	{"1100011", "001", ""}: {"bne", "B"},
	{"1100011", "100", ""}: {"blt", "B"},
	{"1100011", "000", ""}: {"bge", "B"},
	{"1100011", "000", ""}: {"bltu", "B"},
	{"1100011", "000", ""}: {"bgeu", "B"},
	{"0000011", "000", ""}: {"lb", "I"}, 
	{"0000011", "001", ""}: {"lh", "I"}, 
	{"0000011", "010", ""}: {"lw", "I"}, 
	{"0000011", "100", ""}: {"lbu", "I"}, 
	{"0000011", "101", ""}: {"lhu", "I"}, 
	{"0100011", "000", ""}: {"sb", "S"}, 
	{"0100011", "001", ""}: {"sh", "S"}, 
	{"0100011", "010", ""}: {"sw", "S"}, 
	{"0010011", "000", ""}: {"addi", "I"}, 
	{"0010011", "010", ""}: {"slti", "I"}, 
	{"0010011", "011", ""}: {"sltiu", "I"}, 
	{"0010011", "100", ""}: {"xori", "I"}, 
	{"0010011", "110", ""}: {"ori", "I"}, 
	{"0010011", "111", ""}: {"andi", "I"}, 
	{"0010011", "001", "0000000"}: {"slli", "I"}, 
	{"0010011", "101", "0000000"}: {"srli", "I"}, 
	{"0010011", "101", "0100000"}: {"srai", "I"}, 
	{"0110011", "000", "0000000"}: {"add", "R"}, 
	{"0110011", "000", "0100000"}: {"sub", "R"}, 
	{"0110011", "001", "0000000"}: {"sll", "R"}, 
	{"0110011", "010", "0000000"}: {"slt", "R"}, 
	{"0110011", "011", "0000000"}: {"sltu", "R"}, 
	{"0110011", "100", "0000000"}: {"xor", "R"}, 
	{"0110011", "101", "0000000"}: {"srl", "R"}, 
	{"0110011", "101", "0100000"}: {"sra", "R"}, 
	{"0110011", "110", "0000000"}: {"or", "R"}, 
	{"0110011", "111", "0000000"}: {"and", "R"}, 
	{"0001111", "000", ""}: {"fence", ""}, 
	{"", "", ""}: {"fence.tso", ""}, 
	{"", "", ""}: {"pause", ""}, 
	{"1110011", "000", ""}: {"ecall", "I"}, 
	{"", "", ""}: {"ebreak", ""}, 
}

var registerMap = map[int]string{
	0: "zero",
	1:   "ra",
	2:   "sp",
	3:   "gp",
	4:   "tp",
	5:   "t0",
	6:   "t1",
	7:   "t2",
	8:   "s0",
	9:   "s1",
	10:   "a0",
	11:   "a1",
	12:   "a2",
	13:   "a3",
	14:   "a4",
	15:   "a5",
	16:   "a6",
	17:   "a7",
	18:   "s2",
	19:   "s3",
	20:   "s4",
	21:   "s5",
	22:   "s6",
	23:   "s7",
	24:   "s8",
	25:   "s9",
	26:  "s10",
	27:  "s11",
	28:   "t3",
	29:   "t4",
	30:   "t5:",
	31:   "t6",
}
