package assembler

import (
	"fmt"
	"strconv"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/command"
)

// table is the symbol table type
type table map[string]int

// static lookup tables from command string to instruction partial values
var (
	jump = map[string]int{
		"JGT": 0b001,
		"JEQ": 0b010,
		"JGE": 0b011,
		"JLT": 0b100,
		"JNE": 0b101,
		"JLE": 0b110,
		"JMP": 0b111,
	}

	comp = map[string]int{
		"0":   0b0101010,
		"1":   0b0111111,
		"-1":  0b0111010,
		"D":   0b0001100,
		"A":   0b0110000,
		"!D":  0b0001101,
		"!A":  0b0110001,
		"-D":  0b0001111,
		"-A":  0b0110011,
		"D+1": 0b0011111,
		"A+1": 0b0110111,
		"D-1": 0b0001110,
		"A-1": 0b0110010,
		"D+A": 0b0000010,
		"D-A": 0b0010011,
		"A-D": 0b0000111,
		"D&A": 0b0000000,
		"D|A": 0b0010101,
		"_a":  0b1101010,
		"_b":  0b1111111,
		"_c":  0b1111010,
		"_d":  0b1001100,
		"M":   0b1110000,
		"_e":  0b1001101,
		"!M":  0b1110001,
		"_f":  0b1001111,
		"-M":  0b1110011,
		"_g":  0b1011111,
		"M+1": 0b1110111,
		"_h":  0b1001110,
		"M-1": 0b1110010,
		"D+M": 0b1000010,
		"M+D": 0b1000010,
		"D-M": 0b1010011,
		"M-D": 0b1000111,
		"D&M": 0b1000000,
		"D|M": 0b1010101,
	}
)

type Program []command.Any

// Assemble commands into HACK machine language.
func Assemble(program Program) ([]string, error) {
	symbols, err := build(program)
	if err != nil {
		return []string{}, fmt.Errorf("build symbols: %v", err)
	}

	instructions, err := assemble(program, symbols)
	if err != nil {
		return []string{}, fmt.Errorf("assemble second pass: %v", err)
	}

	return instructions, nil
}

// build symbol table
func build(program Program) (table, error) {
	symbols := table{
		"SP":     0x0000,
		"LCL":    0x0001,
		"ARG":    0x0002,
		"THIS":   0x0003,
		"THAT":   0x0004,
		"R0":     0x0000,
		"R1":     0x0001,
		"R2":     0x0002,
		"R3":     0x0003,
		"R4":     0x0004,
		"R5":     0x0005,
		"R6":     0x0006,
		"R7":     0x0007,
		"R8":     0x0008,
		"R9":     0x0009,
		"R10":    0x00a,
		"R11":    0x00b,
		"R12":    0x00c,
		"R13":    0x00d,
		"R14":    0x00e,
		"R15":    0x00f,
		"SCREEN": 0x4000,
		"KBD":    0x6000,
	}

	// pass one, add labels to symbol table
	pos := 0
	for _, c := range program {
		if cmd, ok := c.(command.L); ok {
			symbols[cmd.Symbol] = pos
			continue
		}
		pos++
	}
	return symbols, nil
}

// assemble instructions from program and completed symbol table.
func assemble(program Program, symbols table) ([]string, error) {
	// pass two: if encountering an @SYMBOL
	//		if an existing symbol, finalize the CmdA struct
	//		if a new symbol, add to symbol table as a new user defined variable and finalize CmdA struct
	var instructions []string
	userVarPos := 0x010
	for _, c := range program {
		switch cmd := c.(type) {
		case command.C:
			hack, err := cTos(cmd)
			if err != nil {
				return []string{}, fmt.Errorf("failed parsing C cmd: %v", err)
			}
			instructions = append(instructions, hack)
		case command.A:
			// TODO simplify
			if !cmd.Static {
				loc, ok := symbols[cmd.Symbol]
				if !ok {
					loc = userVarPos
					symbols[cmd.Symbol] = loc
					userVarPos++
				}
				cmd.Address = loc
				// cmd.Final = true
			}
			hack := fmt.Sprintf("0%015b", cmd.Address)
			instructions = append(instructions, hack)
		}
	}

	return instructions, nil
}

func cTos(cmd command.C) (string, error) {
	// C command prefix - 3 bits
	p := 0b111

	// comp flags - 7 bits
	c, ok := comp[cmd.C]
	if !ok {
		return "", fmt.Errorf("unknown computation in '%+v': %s", cmd, cmd.C)
	}

	// destination flags - 3 bits
	d := 0
	if cmd.D.A {
		d += 1 << 2
	}
	if cmd.D.D {
		d += 1 << 1
	}
	if cmd.D.M {
		d += 1
	}

	// jump flags - 3 bits
	j, ok := jump[cmd.J]
	if !ok {
		j = 0
	}

	// combine - 3 + 7 + 3 + 3 = 16 bit instruction
	instruction := p<<(16-3) + c<<(16-3-7) + d<<(16-3-7-3) + j
	// FormatUint is slightly faster than fmt.Sprintf("%016b",i), per benchmarking
	return strconv.FormatUint(uint64(instruction), 2), nil
}
