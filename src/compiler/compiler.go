package compiler

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func Compile(tree *ast.Program, name string) *CodeBlock {
	return compileFrame(&ast.BlockStatement{Statements: tree.Statements}, name, tree.Filename)
}

func compileFrame(node ast.Node, name, filename string) *CodeBlock {
	ccb := &codeBlockCompiler{
		constants: newConstantTable(),
		locals:    newStringTable(),
		names:     newStringTable(),
		code:      NewInstSet(),
		filename:  filename,
		name:      name,
	}

	compile(ccb, node)
	if !ccb.code.last().Is(opcode.Return) {
		ccb.code.addInst(opcode.Return)
	}

	code := ccb.code
	c := &CodeBlock{
		Name:         name,
		Filename:     filename,
		LocalCount:   len(ccb.locals.table),
		Code:         code.Assemble(ccb),
		Constants:    ccb.constants.table,
		Names:        ccb.names.table,
		Locals:       ccb.locals.table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
	}

	return c
}

func compile(ccb *codeBlockCompiler, node ast.Node) {
	if node == nil {
		compileLoadNull(ccb)
		return
	}

	switch node := node.(type) {

	case *ast.ExpressionStatement:
		compile(ccb, node.Expression)

	case *ast.BlockStatement:
		compileBlock(ccb, node)

	case *ast.DoExpression:
		compileDoBlock(ccb, node)

	// Literals
	case *ast.IntegerLiteral:
		i := object.MakeIntObj(node.Value)
		ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(i))

	case *ast.NullLiteral:
		compileLoadNull(ccb)

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(str))

	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(float))

	case *ast.Boolean:
		b := object.NativeBoolToBooleanObj(node.Value)
		ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(b))

	case *ast.Array:
		for _, e := range node.Elements {
			compile(ccb, e)
		}
		ccb.code.addInst(opcode.MakeArray, uint16(len(node.Elements)))

	case *ast.HashLiteral:
		for k, v := range node.Pairs {
			compile(ccb, v)
			compile(ccb, k)
		}
		ccb.code.addInst(opcode.MakeMap, uint16(len(node.Pairs)))

	// Expressions
	case *ast.Identifier:
		if ccb.locals.contains(node.Value) {
			ccb.code.addInst(opcode.LoadFast, ccb.locals.indexOf(node.Value))
		} else {
			ccb.code.addInst(opcode.LoadGlobal, ccb.names.indexOf(node.Value))
		}

	case *ast.PrefixExpression:
		compile(ccb, node.Right)

		switch node.Operator {
		case "!":
			ccb.code.addInst(opcode.UnaryNot)
		case "-":
			ccb.code.addInst(opcode.UnaryNeg)
		}

	case *ast.InfixExpression:
		compile(ccb, node.Left)
		compile(ccb, node.Right)

		switch node.Operator {
		case "+":
			ccb.code.addInst(opcode.BinaryAdd)
		case "-":
			ccb.code.addInst(opcode.BinarySub)
		case "*":
			ccb.code.addInst(opcode.BinaryMul)
		case "/":
			ccb.code.addInst(opcode.BinaryDivide)
		case "%":
			ccb.code.addInst(opcode.BinaryMod)
		case "<<":
			ccb.code.addInst(opcode.BinaryShiftL)
		case ">>":
			ccb.code.addInst(opcode.BinaryShiftR)
		case "&":
			ccb.code.addInst(opcode.BinaryAnd)
		case "&^":
			ccb.code.addInst(opcode.BinaryAndNot)
		case "|":
			ccb.code.addInst(opcode.BinaryOr)
		case "^":
			ccb.code.addInst(opcode.BinaryNot)
		case "<":
			ccb.code.addInst(opcode.Compare, uint16(opcode.CmpLT))
		case ">":
			ccb.code.addInst(opcode.Compare, uint16(opcode.CmpGT))
		case "==":
			ccb.code.addInst(opcode.Compare, uint16(opcode.CmpEq))
		case "!=":
			ccb.code.addInst(opcode.Compare, uint16(opcode.CmpNotEq))
		case "<=":
			ccb.code.addInst(opcode.Compare, uint16(opcode.CmpLTEq))
		case ">=":
			ccb.code.addInst(opcode.Compare, uint16(opcode.CmpGTEq))
		}

	case *ast.CallExpression:
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compile(ccb, node.Arguments[i])
		}
		compile(ccb, node.Function)
		ccb.code.addInst(opcode.Call, uint16(len(node.Arguments)))

	case *ast.ReturnStatement:
		compile(ccb, node.Value)
		ccb.code.addInst(opcode.Return)

	case *ast.DefStatement:
		compile(ccb, node.Value)

		if node.Const {
			ccb.code.addInst(opcode.StoreConst, ccb.locals.indexOf(node.Name.Value))
		} else {
			ccb.code.addInst(opcode.Define, ccb.locals.indexOf(node.Name.Value))
		}

	case *ast.AssignStatement:
		compile(ccb, node.Value)

		if indexed, ok := node.Left.(*ast.IndexExpression); ok {
			compile(ccb, indexed.Index)
			compile(ccb, indexed.Left)
			ccb.code.addInst(opcode.StoreIndex)
			break
		}

		if attrib, ok := node.Left.(*ast.AttributeExpression); ok {
			compile(ccb, attrib.Left)
			ccb.code.addInst(opcode.StoreAttribute, ccb.names.indexOf(attrib.Index.String()))
			break
		}

		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			panic("Assignment to non ident or index")
		}

		if ccb.locals.contains(ident.Value) {
			ccb.code.addInst(opcode.StoreFast, ccb.locals.indexOf(ident.Value))
		} else {
			ccb.code.addInst(opcode.StoreGlobal, ccb.names.indexOf(ident.Value))
		}

	case *ast.DeleteStatement:
		ccb.code.addInst(opcode.DeleteFast, ccb.locals.indexOf(node.Name))

	case *ast.IfExpression:
		compileIfStatement(ccb, node)

	case *ast.CompareExpression:
		compileCompareExpression(ccb, node)

	case *ast.ImportStatement:
		str := &object.String{Value: node.Path.Value}
		ccb.code.addInst(opcode.Import, ccb.constants.indexOf(str))
		ccb.code.addInst(opcode.Define, ccb.locals.indexOf(node.Name.Value))

	case *ast.FunctionLiteral:
		compileFunction(ccb, node, false, false)

	case *ast.IndexExpression:
		compile(ccb, node.Index)
		compile(ccb, node.Left)
		ccb.code.addInst(opcode.LoadIndex)

	case *ast.LoopStatement:
		compileLoop(ccb, node)

	case *ast.IterLoopStatement:
		compileIterLoop(ccb, node)

	case *ast.ContinueStatement:
		if !ccb.inLoop {
			panic("continue used in non-loop block")
		}
		ccb.code.addInst(opcode.Continue)

	case *ast.BreakStatement:
		if !ccb.inLoop {
			panic("break used in non-loop block")
		}
		ccb.code.addInst(opcode.Break)

	case *ast.TryCatchExpression:
		compileTryCatch(ccb, node)

	case *ast.ThrowStatement:
		compile(ccb, node.Expression)
		ccb.code.addInst(opcode.Throw)

	case *ast.ClassLiteral:
		compileClassLiteral(ccb, node)

	case *ast.NewInstance:
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compile(ccb, node.Arguments[i])
		}
		compile(ccb, node.Class)

		ccb.code.addInst(opcode.MakeInstance, uint16(len(node.Arguments)))

	case *ast.AttributeExpression:
		compile(ccb, node.Left)
		ccb.code.addInst(opcode.LoadAttribute, ccb.names.indexOf(node.Index.String()))

	case *ast.PassStatement:
		// Ignore

	// Not implemented yet
	case *ast.Program:
		panic("ast.Program Not implemented yet")

	default:
		panic(fmt.Sprintf("Node type not implemented: %T", node))
	}
}
