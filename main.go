package main

import (
	"errors"
	"fmt"
	"os"
)

type Instruction struct {
	operator int
	operand  int
}

const TAPE_MAX_SIZE = 4 * 1024
const STACK_MAX_SIZE = 512

const (
	BF_OP_MOVE_RIGHT = iota
	BF_OP_MOVE_LEFT
	BF_OP_INCREMENT
	BF_OP_DECREMENT
	BF_OP_WRITE
	BF_OP_READ
	BF_OP_JUMP_IF_ZERO
	BF_OP_JUMP_UNLESS_ZERO
)

func CompileBrainfuck(instructions []byte) (program []Instruction, error error) {
	var pc, sp int = 0, 0
	var stack = []int{}

	for index, op := range instructions {
		switch op {
		case '>':
			program = append(program, Instruction{BF_OP_MOVE_RIGHT, 0})
			break
		case '<':
			program = append(program, Instruction{BF_OP_MOVE_LEFT, 0})
			break
		case '+':
			program = append(program, Instruction{BF_OP_INCREMENT, 0})
			break
		case '-':
			program = append(program, Instruction{BF_OP_DECREMENT, 0})
			break
		case '.':
			program = append(program, Instruction{BF_OP_WRITE, 0})
			break
		case ',':
			program = append(program, Instruction{BF_OP_READ, 0})
			break
		case '[':
			program = append(program, Instruction{BF_OP_JUMP_IF_ZERO, 0})
			stack = append(stack, pc)
			break
		case ']':
			// If stack is clear, it means ] does not have a correspondent
			// ] instruction.
			if len(stack) == 0 {
				return nil, errors.New(fmt.Sprintf("\n\tindex %d: Stack is clear", index))
			}

			sp = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			program = append(program, Instruction{BF_OP_JUMP_UNLESS_ZERO, sp})

			// Add current PC to correspondent JMP instruction
			program[sp].operand = pc
			break
		default:
			pc--
		}
		pc++
	}

	// If stack is not clear, then a error happened.
	if len(stack) != 0 {
		return nil, errors.New(fmt.Sprintf("\n\tStack is not empty at index %d", len(stack)))
	}

	return
}

func main() {
	args := os.Args

	if len(args) < 2 {
		panic(fmt.Sprintf("Usage: %s <rom-file>\n", os.Args[0]))
	}

	file, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Printf("Error on reading file: %d\n", err)
		return
	}

	program, err := CompileBrainfuck(file)
	if err != nil {
		fmt.Printf("An error occured: %d", err)
		return
	}

	for index, instruction := range program {
		fmt.Printf("%d: %d\t %d\n", index, instruction.operator, instruction.operand)
	}
}
