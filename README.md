[![Go](https://github.com/jeffgreenca/n2t-asm/actions/workflows/go.yml/badge.svg)](https://github.com/jeffgreenca/n2t-asm/actions/workflows/go.yml)

# n2t-asm

An assembler for [nand2tetris](https://www.nand2tetris.org/), written in Go.

# usage

```
$ ./scripts/build.sh
$ ./n2t-asm program.asm > program.hack
# or, via stdin
$ cat program.asm | ./n2t-asm > program.hack
```

# testing and building

```
$ go test ./...
$ ./scripts/build.sh
```

# validation

If you have nand2tetris software installed locally, this script will assmeble the corpus of `.asm` files from nand2tetris using both this assembler and the nand2tetris provided assembler, diffing the results.

```
$ export N2T_PATH=/path/to/your/nand2tetris
$ ./scripts/build.sh && ./scripts/compare.sh
```
