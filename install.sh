#!/bin/bash
NAME=n2t-asm
go build -o ${NAME} cmd/${NAME}/main.go && sudo mv ${NAME} /usr/local/bin/${NAME}
