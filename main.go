package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Instruction struct {
	operator uint8
	operand  uint8
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

var DEBUG_INSTRUCTION_NAMES = map[uint8]string{
	BF_OP_MOVE_RIGHT:       "MOVE_RIGHT",
	BF_OP_MOVE_LEFT:        "MOVE_LEFT",
	BF_OP_INCREMENT:        "INCREMENT",
	BF_OP_DECREMENT:        "DECREMENT",
	BF_OP_WRITE:            "WRITE",
	BF_OP_READ:             "READ",
	BF_OP_JUMP_IF_ZERO:     "JUMP_IF_ZERO",
	BF_OP_JUMP_UNLESS_ZERO: "JUMP_UNLESS_ZERO",
}

func CompileBrainfuck(instructions []byte) (program []Instruction, error error) {
	var pc, sp uint8 = 0, 0
	var stack = []uint8{}

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

func RunBrainfuck(instructions []Instruction, debug bool) {
	tape := make([]int8, TAPE_MAX_SIZE)
	var head uint8 = 0

	for pc := 0; pc < len(instructions); pc++ {
		i := instructions[pc]

		if debug {
			var sample_start int
			var sample_end int

			if int(head)-7 < 0 {
				sample_start = 0
				sample_end = 15
			} else {
				sample_start = int(head) - 7
				sample_end = int(head) + 8
			}

			fmt.Printf("\nHead: %d\tPC: %d\t\n", head, pc)
			for dI := sample_start; dI < sample_end; dI++ {
				if dI == int(head) {
					fmt.Printf("[ >%d ] ", tape[dI])
				} else {
					fmt.Printf("[ %d ] ", tape[dI])
				}
			}
			fmt.Printf("\nCurrent Instruction: %s\n", DEBUG_INSTRUCTION_NAMES[i.operator])

			// Block until input
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
		}

		switch i.operator {
		case BF_OP_MOVE_RIGHT:
			head++
		case BF_OP_MOVE_LEFT:
			head--
		case BF_OP_INCREMENT:
			tape[head]++
		case BF_OP_DECREMENT:
			tape[head]--
		case BF_OP_WRITE:
			fmt.Printf("%c", tape[head])
		case BF_OP_READ:
		case BF_OP_JUMP_IF_ZERO:
			if tape[head] == 0 {
				pc = int(i.operand)
			}
		case BF_OP_JUMP_UNLESS_ZERO:
			if tape[head] != 0 {
				pc = int(i.operand)
			}
		default:
			panic("Unknown operator")
		}
	}
}

func main() {
	args := os.Args
	debug := os.Getenv("DEBUG") == "true"

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

	if debug {
		fmt.Println("Running on debug mode")
	}
	RunBrainfuck(program, debug)
}
