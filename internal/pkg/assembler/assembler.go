package assembler

import (
	"fmt"

	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/parser"
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
			panic("Unknown instruction type")
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
			hack := C(cmd)
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
	JUMP = map[string]string{
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}

	COMP = map[string]string{
		"0":   "0101010",
		"1":   "0111111",
		"-1":  "0111010",
		"D":   "0001100",
		"A":   "0110000",
		"!D":  "0001101",
		"!A":  "0110001",
		"-D":  "0001111",
		"-A":  "0110011",
		"D+1": "0011111",
		"A+1": "0110111",
		"D-1": "0001110",
		"A-1": "0110010",
		"D+A": "0000010",
		"D-A": "0010011",
		"A-D": "0000111",
		"D&A": "0000000",
		"D|A": "0010101",
		"_a":  "1101010",
		"_b":  "1111111",
		"_c":  "1111010",
		"_d":  "1001100",
		"M":   "1110000",
		"_e":  "1001101",
		"!M":  "1110001",
		"_f":  "1001111",
		"-M":  "1110011",
		"_g":  "1011111",
		"M+1": "1110111",
		"_h":  "1001110",
		"M-1": "1110010",
		"D+M": "1000010",
		"M+D": "1000010",
		"D-M": "1010011",
		"M-D": "1000111",
		"D&M": "1000000",
		"D|M": "1010101",
	}
)

func C(cmd parser.CmdC) string {
	// C command prefix
	c := "111"
	// convert comp part
	comp, ok := COMP[cmd.C]
	if !ok {
		panic(fmt.Sprintf("Unknown computation: %s", cmd.C))
	}

	// convert dest part
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
	dest := fmt.Sprintf("%03b", d)

	// convert jump part
	jump, ok := JUMP[cmd.J]
	if !ok {
		jump = "000"
	}

	// combine
	return c + comp + dest + jump
}
