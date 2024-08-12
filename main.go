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
				return nil, fmt.Errorf("Error on parsing index %d: Stack is clear", index)
			}

			jmp_pc := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			program = append(program, Instruction{BF_OP_JUMP_UNLESS_ZERO, jmp_pc})

			// Add current PC to correspondent JMP instruction
			program[jmp_pc].operand = pc
			break
		default:
			pc--
		}
		pc++
	}

	// If stack is not clear, then a error happened.
	if len(stack) != 0 {
		return nil, fmt.Errorf("Stack is %d", len(stack))
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
		panic(fmt.Sprint("Error on reading file:", err))
	}

}
