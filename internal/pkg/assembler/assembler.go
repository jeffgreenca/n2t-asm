package assembler

import (
	"fmt"
	"strings"

	"bitbucket.org/jeffgreenca/n2t-asm/internal/pkg/parser"
)

// Assemble converts a series of parsed Commands into HACK machine language instructions
func Assemble(program []parser.Command) ([]string, error) {
	var instructions []string

	// this is a simple, one-pass approach that doesn't support symbols
	for _, c := range program {
		switch c.Type {
		case parser.L_COMMAND:
			// not implemented
		case parser.C_COMMAND:
			cmd := c.C.(parser.CmdC)
			hack := C(cmd)
			instructions = append(instructions, hack)
		case parser.A_COMMAND:
			cmd := c.C.(parser.CmdA)
			if !cmd.Final {
				panic("Nonfinalized A command encountered")
			}
			hack := fmt.Sprintf("0%015b", cmd.Address)
			instructions = append(instructions, hack)
		default:
			panic("Unknown instruction type")
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
		panic("Unknown computation!")
	}

	// convert dest part
	d := []string{"0", "0", "0"}
	if cmd.D.A {
		d[0] = "1"
	}
	if cmd.D.D {
		d[1] = "1"
	}
	if cmd.D.M {
		d[2] = "1"
	}
	dest := strings.Join(d, "")

	// convert jump part
	jump, ok := JUMP[cmd.J]
	if !ok {
		jump = "000"
	}

	// combine
	return c + comp + dest + jump
}
