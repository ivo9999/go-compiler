package compiler

import (
	"fmt"
	"go-interp/ast"
	"go-interp/code"
	"go-interp/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.Emit(code.OpPop)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.Emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.Emit(code.OpAdd)
		case "-":
			c.Emit(code.OpSub)
		case "*":
			c.Emit(code.OpMul)
		case "/":
			c.Emit(code.OpDiv)
		case ">":
			c.Emit(code.OpGreaterThan)
		case "==":
			c.Emit(code.OpEqual)
		case "!=":
			c.Emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.Emit(code.OpConstant, c.AddConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.Emit(code.OpTrue)
		} else {
			c.Emit(code.OpFalse)
		}

	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) AddConstant(o object.Object) int {
	c.constants = append(c.constants, o)
	return len(c.constants) - 1
}

func (c *Compiler) Emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
