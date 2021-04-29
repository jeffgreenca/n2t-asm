#!/bin/bash
N2T_PATH=${N2T_PATH:-$HOME/Documents/nand2tetris}
ASM_1=${N2T_PATH}/tools/Assembler.sh
ASM_2=./n2t-asm

TARGETS=$(find ${N2T_PATH} -type f -name \*.asm)

WORKDIR=${PWD}/temp
mkdir -p ${WORKDIR} && rm -rf ${WORKDIR}/ && mkdir ${WORKDIR}

FAIL=0
for t in ${TARGETS}; do
	fn=$(basename ${t})
	cp ${t} ${WORKDIR}/${fn}
	echo temp/${fn}
	${ASM_1} ${WORKDIR}/${fn} > /dev/null
	${ASM_2} ${WORKDIR}/${fn} > ${WORKDIR}/${fn:0:-4}.hack2
	diff ${WORKDIR}/${fn:0:-4}.hack2 ${WORKDIR}/${fn:0:-4}.hack || export FAIL=1
done

if [[ $FAIL -eq 0 ]]; then
	echo "No differences found"
	rm -rf ${WORKDIR}
fi
