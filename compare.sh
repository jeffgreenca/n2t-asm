#!/bin/bash
ASM_1=~/Documents/nand2tetris/tools/Assembler.sh
ASM_2=n2t-asm

TARGETS=$(find ~/Documents/nand2tetris/ -name \*.asm -type f)

WORKDIR=${PWD}/temp
mkdir ${WORKDIR} || rm -rf ${WORKDIR} && mkdir ${WORKDIR}

for t in ${TARGETS}; do
	fn=$(basename ${t})
	cp ${t} ${WORKDIR}/${fn}
	${ASM_1} ${WORKDIR}/${fn}
	diff --color=always <(${ASM_2} ${WORKDIR}/${fn}) ${WORKDIR}/${fn:0:-4}.hack
done
