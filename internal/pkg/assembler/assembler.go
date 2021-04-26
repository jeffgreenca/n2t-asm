package assembler

import (
	"fmt"
	"strconv"

	"github.com/jeffgreenca/n2t-asm/internal/pkg/parser"
)

// Assemble converts a series of parsed Commands into HACK machine language instructions
func Assemble(program []parser.Command) ([]string, error) {
	// init symbol table
	symbols := map[string]int{
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

	// pass one: scan program for labels, adding to symbol table
	pos := 0
	for _, c := range program {
		switch c.Type {
		case parser.L_COMMAND:
			cmd := c.C.(parser.CmdL)
			symbols[cmd.Symbol] = pos
		case parser.C_COMMAND, parser.A_COMMAND:
			pos++
		default:
			return []string{}, fmt.Errorf("unknown instruction type in command '%+v': %v", c, c.Type)
		}
	}

	// pass two: if encountering an @SYMBOL
	//		if an existing symbol, finalize the CmdA struct
	//		if a new symbol, add to symbol table as a new user defined variable and finalize CmdA struct
	var instructions []string
	userVarPos := 0x010
	for _, c := range program {
		switch c.Type {
		case parser.C_COMMAND:
			cmd := c.C.(parser.CmdC)
			hack, err := cTos(cmd)
			if err != nil {
				return []string{}, fmt.Errorf("failed parsing C cmd: %v", err)
			}
			instructions = append(instructions, hack)
		case parser.A_COMMAND:
			cmd := c.C.(parser.CmdA)
			if !cmd.Final {
				loc, ok := symbols[cmd.Symbol]
				if !ok {
					loc = userVarPos
					symbols[cmd.Symbol] = loc
					userVarPos++
				}
				cmd.Address = loc
				cmd.Final = true
			}
			hack := fmt.Sprintf("0%015b", cmd.Address)
			instructions = append(instructions, hack)
		}
	}

	return instructions, nil
}

var (
	JUMP = map[string]int{
		"JGT": 0b001,
		"JEQ": 0b010,
		"JGE": 0b011,
		"JLT": 0b100,
		"JNE": 0b101,
		"JLE": 0b110,
		"JMP": 0b111,
	}

	COMP = map[string]int{
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

func cTos(cmd parser.CmdC) (string, error) {
	// C command prefix - 3 bits
	p := 0b111

	// comp flags - 7 bits
	c, ok := COMP[cmd.C]
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
	j, ok := JUMP[cmd.J]
	if !ok {
		j = 0
	}

	// combine - 3 + 7 + 3 + 3 = 16 bit instruction
	instruction := p<<(16-3) + c<<(16-3-7) + d<<(16-3-7-3) + j
	// FormatUint is slightly faster than fmt.Sprintf("%016b",i), per benchmarking
	return strconv.FormatUint(uint64(instruction), 2), nil
}
