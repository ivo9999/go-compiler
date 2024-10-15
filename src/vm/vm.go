package vm

import (
	"fmt"
	"go-interp/code"
	"go-interp/compiler"
	"go-interp/object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.Push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.OpAdd:
			left := vm.Pop()
			right := vm.Pop()

			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value

			result := rightValue + leftValue
			vm.Push(&object.Integer{Value: result})
		}
	}
	return nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) Pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}
