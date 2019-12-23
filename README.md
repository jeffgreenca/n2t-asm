# nand2tetris assembler in go

An assembler for the nand2tetris assembly language, written in Go.

## coverage
```
∅ cmd/n2t-asm
✓ internal/pkg/lex (3ms) (coverage: 75.0% of statements)
✓ internal/pkg/assembler (4ms) (coverage: 72.1% of statements)
✓ internal/pkg/parser (5ms) (coverage: 89.2% of statements)
✓ tests (4ms)

DONE 15 tests in 0.847s
```

## things I like

- test coverage made it easy to track down when problems occurred
- it works

## things I don't like

- global variables in the parser - would rather implement as recievers with an appropriate struct
- each of `CmdA`, `CmdC`, etc. should implement an abstract `Cmd` interface for use in the assembler package
- concepts that span packages should be split out - for example, commands, the lexers symbol definitions, and so on
- error handling is pretty rough, used a lot of panics rather than letting errors bubble up

# License

MIT License.
