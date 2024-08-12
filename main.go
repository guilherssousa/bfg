package main

import (
	"fmt"
	"os"
	"strconv"
)

const TAPE_MAX_SIZE = 4 * 1024

type Tape = [TAPE_MAX_SIZE]uint8

type Program struct {
	Tape *Tape
	PC   int
}

func (p *Program) Run(instructions []byte) {
	for _, instruction := range instructions {
		switch instruction {
		case '>':
			p.PC++
			break
		case '<':
			p.PC++
			break
		case '+':
			p.Tape[p.PC]++
			break
		case '-':
			p.Tape[p.PC]++
			break
		case '.':
			fmt.Print(strconv.Itoa(int(p.Tape[p.PC])))
			break
		case ',':
			// TODO: Implement ,
			break
		case '[':
			break
		case ']':
			break
		default:
			break
		}
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		panic(fmt.Sprintf("Usage: %s <rom-file>\n", os.Args[0]))
	}

	file, err := os.ReadFile(args[1])
	if err != nil {
		panic(fmt.Sprint("Error on reading file:", err))
	}

	p := Program{Tape: &Tape{}, PC: 0}
	p.Run(file)
}
