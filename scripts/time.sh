#!/bin/bash
N2T_PATH=${N2T_PATH:-$HOME/Documents/nand2tetris}
target=$N2T_PATH/projects/06/pong/Pong.asm
time ./n2t-asm $target > /dev/null
