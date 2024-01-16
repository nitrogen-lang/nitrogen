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
		linenum:   0,
	}

	compile(ccb, node)
	if !ccb.code.last().Is(opcode.Return) {
		ccb.code.addInst(opcode.Return, ccb.linenum)
	}

	code := ccb.code
	assembledCode, lineOffsets := code.Assemble(ccb)
	c := &CodeBlock{
		Name:         name,
		Filename:     filename,
		LocalCount:   len(ccb.locals.table),
		Code:         assembledCode,
		Constants:    ccb.constants.table,
		Names:        ccb.names.table,
		Locals:       ccb.locals.table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
		LineOffsets:  lineOffsets,
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
		ccb.linenum = node.Token.Pos.Line
		i := object.MakeIntObj(node.Value)
		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(i))

	case *ast.NullLiteral:
		ccb.linenum = node.Token.Pos.Line
		compileLoadNull(ccb)

	case *ast.StringLiteral:
		ccb.linenum = node.Token.Pos.Line
		str := &object.String{Value: node.Value}
		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(str))

	case *ast.FloatLiteral:
		ccb.linenum = node.Token.Pos.Line
		float := &object.Float{Value: node.Value}
		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(float))

	case *ast.Boolean:
		ccb.linenum = node.Token.Pos.Line
		b := object.NativeBoolToBooleanObj(node.Value)
		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(b))

	case *ast.Array:
		ccb.linenum = node.Token.Pos.Line
		for _, e := range node.Elements {
			compile(ccb, e)
		}
		ccb.code.addInst(opcode.MakeArray, ccb.linenum, uint16(len(node.Elements)))

	case *ast.HashLiteral:
		ccb.linenum = node.Token.Pos.Line
		for k, v := range node.Pairs {
			compile(ccb, v)
			compile(ccb, k)
		}
		ccb.code.addInst(opcode.MakeMap, ccb.linenum, uint16(len(node.Pairs)))

	case *ast.InterfaceLiteral:
		ccb.linenum = node.Token.Pos.Line
		iface := &object.Interface{
			Name:    node.Name,
			Methods: make(map[string]*object.IfaceMethodDef, len(node.Methods)),
		}

		for name, def := range node.Methods {
			iface.Methods[name] = &object.IfaceMethodDef{
				Name:       def.Name,
				Parameters: def.Params,
			}
		}

		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(iface))

	// Expressions
	case *ast.Identifier:
		ccb.linenum = node.Token.Pos.Line
		if ccb.locals.contains(node.Value) {
			ccb.code.addInst(opcode.LoadFast, ccb.linenum, ccb.locals.indexOf(node.Value))
		} else {
			ccb.code.addInst(opcode.LoadGlobal, ccb.linenum, ccb.names.indexOf(node.Value))
		}

	case *ast.PrefixExpression:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Right)

		switch node.Operator {
		case "!":
			ccb.code.addInst(opcode.UnaryNot, ccb.linenum)
		case "-":
			ccb.code.addInst(opcode.UnaryNeg, ccb.linenum)
		}

	case *ast.InfixExpression:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Left)
		compile(ccb, node.Right)

		switch node.Operator {
		case "+":
			ccb.code.addInst(opcode.BinaryAdd, ccb.linenum)
		case "-":
			ccb.code.addInst(opcode.BinarySub, ccb.linenum)
		case "*":
			ccb.code.addInst(opcode.BinaryMul, ccb.linenum)
		case "/":
			ccb.code.addInst(opcode.BinaryDivide, ccb.linenum)
		case "%":
			ccb.code.addInst(opcode.BinaryMod, ccb.linenum)
		case "<<":
			ccb.code.addInst(opcode.BinaryShiftL, ccb.linenum)
		case ">>":
			ccb.code.addInst(opcode.BinaryShiftR, ccb.linenum)
		case "&":
			ccb.code.addInst(opcode.BinaryAnd, ccb.linenum)
		case "&^":
			ccb.code.addInst(opcode.BinaryAndNot, ccb.linenum)
		case "|":
			ccb.code.addInst(opcode.BinaryOr, ccb.linenum)
		case "^":
			ccb.code.addInst(opcode.BinaryNot, ccb.linenum)
		case "<":
			ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpLT))
		case ">":
			ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpGT))
		case "==":
			ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpEq))
		case "!=":
			ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpNotEq))
		case "<=":
			ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpLTEq))
		case ">=":
			ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpGTEq))
		case "implements":
			ccb.code.addInst(opcode.Implements, ccb.linenum)
		}

	case *ast.CallExpression:
		ccb.linenum = node.Token.Pos.Line
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compile(ccb, node.Arguments[i])
		}
		compile(ccb, node.Function)
		ccb.code.addInst(opcode.Call, ccb.linenum, uint16(len(node.Arguments)))

	case *ast.ReturnStatement:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Value)
		ccb.code.addInst(opcode.Return, ccb.linenum)

	case *ast.DefStatement:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Value)

		if node.Const {
			ccb.code.addInst(opcode.StoreConst, ccb.linenum, ccb.locals.indexOf(node.Name.Value))
		} else {
			ccb.code.addInst(opcode.Define, ccb.linenum, ccb.locals.indexOf(node.Name.Value))
		}

	case *ast.AssignStatement:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Value)

		if indexed, ok := node.Left.(*ast.IndexExpression); ok {
			compile(ccb, indexed.Index)
			compile(ccb, indexed.Left)
			ccb.code.addInst(opcode.StoreIndex, ccb.linenum)
			break
		}

		if attrib, ok := node.Left.(*ast.AttributeExpression); ok {
			compile(ccb, attrib.Left)
			ccb.code.addInst(opcode.StoreAttribute, ccb.linenum, ccb.names.indexOf(attrib.Index.String()))
			break
		}

		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			panic("Assignment to non ident or index")
		}

		if ccb.locals.contains(ident.Value) {
			ccb.code.addInst(opcode.StoreFast, ccb.linenum, ccb.locals.indexOf(ident.Value))
		} else {
			ccb.code.addInst(opcode.StoreGlobal, ccb.linenum, ccb.names.indexOf(ident.Value))
		}

	case *ast.DeleteStatement:
		ccb.linenum = node.Token.Pos.Line
		ccb.code.addInst(opcode.DeleteFast, ccb.linenum, ccb.locals.indexOf(node.Name))

	case *ast.IfExpression:
		compileIfStatement(ccb, node)

	case *ast.CompareExpression:
		compileCompareExpression(ccb, node)

	case *ast.ImportStatement:
		ccb.linenum = node.Token.Pos.Line
		str := &object.String{Value: node.Path.Value}
		ccb.code.addInst(opcode.Import, ccb.linenum, ccb.constants.indexOf(str))
		ccb.code.addInst(opcode.Define, ccb.linenum, ccb.locals.indexOf(node.Name.Value))

	case *ast.FunctionLiteral:
		compileFunction(ccb, node, false, false)

	case *ast.IndexExpression:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Left)
		compile(ccb, node.Index)
		ccb.code.addInst(opcode.LoadIndex, ccb.linenum)

	case *ast.LoopStatement:
		compileLoop(ccb, node)

	case *ast.IterLoopStatement:
		compileIterLoop(ccb, node)

	case *ast.ContinueStatement:
		ccb.linenum = node.Token.Pos.Line
		if !ccb.inLoop {
			panic("continue used in non-loop block")
		}
		ccb.code.addInst(opcode.Continue, ccb.linenum)

	case *ast.BreakStatement:
		ccb.linenum = node.Token.Pos.Line
		if !ccb.inLoop {
			panic("break used in non-loop block")
		}
		ccb.code.addInst(opcode.Break, ccb.linenum)

	case *ast.ClassLiteral:
		compileClassLiteral(ccb, node)

	case *ast.NewInstance:
		ccb.linenum = node.Token.Pos.Line
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compile(ccb, node.Arguments[i])
		}
		compile(ccb, node.Class)

		ccb.code.addInst(opcode.MakeInstance, ccb.linenum, uint16(len(node.Arguments)))

	case *ast.AttributeExpression:
		ccb.linenum = node.Token.Pos.Line
		compile(ccb, node.Left)
		ccb.code.addInst(opcode.LoadAttribute, ccb.linenum, ccb.names.indexOf(node.Index.String()))

	case *ast.PassStatement:
		ccb.linenum = node.Token.Pos.Line
		// Ignore

	// Not implemented yet
	case *ast.Program:
		panic("ast.Program Not implemented yet")

	default:
		panic(fmt.Sprintf("Node type not implemented: %T", node))
	}
}
