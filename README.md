[![Go](https://github.com/jeffgreenca/n2t-asm/actions/workflows/go.yml/badge.svg)](https://github.com/jeffgreenca/n2t-asm/actions/workflows/go.yml)

# n2t-asm

An assembler for [nand2tetris](https://www.nand2tetris.org/), written in Go.

# example usage

```
$ ./scripts/build.sh
$ ./n2t-asm program.asm > program.hack
# or, via stdin
$ cat program.asm | ./n2t-asm > program.hack
```

# testing

```
# run tests
$ go test ./...

# compare output to official assembler
# on all '*.asm' provided by nand2tetris project files
$ export N2T_PATH=/path/to/your/nand2tetris
$ ./scripts/build.sh && ./scripts/compare.sh
```

# License

MIT License.
