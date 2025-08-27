package compiler

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/elemental/compile"
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
)

func Compile(tree *ast.Program, name string) *compile.CodeBlock {
	return compileFrame(&ast.BlockStatement{Statements: tree.Statements}, name, tree.Filename)
}

func compileFrame(node ast.Node, name, filename string) *compile.CodeBlock {
	ccb := &compile.CodeBlockCompiler{
		Constants: compile.NewConstantTable(),
		Locals:    compile.NewStringTable(),
		Names:     compile.NewStringTable(),
		Code:      compile.NewInstSet(),
		Filename:  filename,
		Name:      name,
		Linenum:   0,
	}

	compileMain(ccb, node)
	if !ccb.Code.Last().Is(opcode.Return) {
		ccb.Code.AddInst(opcode.Return, ccb.Linenum)
	}

	code := ccb.Code
	assembledCode, lineOffsets := code.Assemble(ccb)
	c := &compile.CodeBlock{
		Name:         name,
		Filename:     filename,
		LocalCount:   len(ccb.Locals.Table),
		Code:         assembledCode,
		Constants:    ccb.Constants.Table,
		Names:        ccb.Names.Table,
		Locals:       ccb.Locals.Table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
		LineOffsets:  lineOffsets,
	}

	return c
}

func compileMain(ccb *compile.CodeBlockCompiler, node ast.Node) {
	if node == nil {
		compileLoadNull(ccb)
		return
	}

	switch node := node.(type) {

	case *ast.ExpressionStatement:
		compileMain(ccb, node.Expression)

	case *ast.BlockStatement:
		compileBlock(ccb, node)

	case *ast.DoExpression:
		compileDoBlock(ccb, node)

	// Literals
	case *ast.IntegerLiteral:
		ccb.Linenum = node.Token.Pos.Line
		i := object.MakeIntObj(node.Value)
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(i))

	case *ast.NullLiteral:
		ccb.Linenum = node.Token.Pos.Line
		compileLoadNull(ccb)

	case *ast.StringLiteral:
		ccb.Linenum = node.Token.Pos.Line
		str := &object.String{Value: node.Value}
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(str))

	case *ast.ByteStringLiteral:
		ccb.Linenum = node.Token.Pos.Line
		str := &object.ByteString{Value: node.Value}
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(str))

	case *ast.FloatLiteral:
		ccb.Linenum = node.Token.Pos.Line
		float := &object.Float{Value: node.Value}
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(float))

	case *ast.Boolean:
		ccb.Linenum = node.Token.Pos.Line
		b := object.NativeBoolToBooleanObj(node.Value)
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(b))

	case *ast.Array:
		ccb.Linenum = node.Token.Pos.Line
		for _, e := range node.Elements {
			compileMain(ccb, e)
		}
		ccb.Code.AddInst(opcode.MakeArray, ccb.Linenum, uint16(len(node.Elements)))

	case *ast.HashLiteral:
		ccb.Linenum = node.Token.Pos.Line
		for k, v := range node.Pairs {
			compileMain(ccb, v)
			compileMain(ccb, k)
		}
		ccb.Code.AddInst(opcode.MakeMap, ccb.Linenum, uint16(len(node.Pairs)))

	case *ast.InterfaceLiteral:
		ccb.Linenum = node.Token.Pos.Line
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

		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(iface))

	// Expressions
	case *ast.Identifier:
		ccb.Linenum = node.Token.Pos.Line
		if ccb.Locals.Contains(node.Value) {
			ccb.Code.AddInst(opcode.LoadFast, ccb.Linenum, ccb.Locals.IndexOf(node.Value))
		} else {
			ccb.Code.AddInst(opcode.LoadGlobal, ccb.Linenum, ccb.Names.IndexOf(node.Value))
		}

	case *ast.PrefixExpression:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Right)

		switch node.Operator {
		case "!":
			ccb.Code.AddInst(opcode.UnaryNot, ccb.Linenum)
		case "-":
			ccb.Code.AddInst(opcode.UnaryNeg, ccb.Linenum)
		}

	case *ast.InfixExpression:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Left)
		compileMain(ccb, node.Right)

		switch node.Operator {
		case "+":
			ccb.Code.AddInst(opcode.BinaryAdd, ccb.Linenum)
		case "-":
			ccb.Code.AddInst(opcode.BinarySub, ccb.Linenum)
		case "*":
			ccb.Code.AddInst(opcode.BinaryMul, ccb.Linenum)
		case "/":
			ccb.Code.AddInst(opcode.BinaryDivide, ccb.Linenum)
		case "%":
			ccb.Code.AddInst(opcode.BinaryMod, ccb.Linenum)
		case "<<":
			ccb.Code.AddInst(opcode.BinaryShiftL, ccb.Linenum)
		case ">>":
			ccb.Code.AddInst(opcode.BinaryShiftR, ccb.Linenum)
		case "&":
			ccb.Code.AddInst(opcode.BinaryAnd, ccb.Linenum)
		case "&^":
			ccb.Code.AddInst(opcode.BinaryAndNot, ccb.Linenum)
		case "|":
			ccb.Code.AddInst(opcode.BinaryOr, ccb.Linenum)
		case "^":
			ccb.Code.AddInst(opcode.BinaryNot, ccb.Linenum)
		case "<":
			ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpLT))
		case ">":
			ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpGT))
		case "==":
			ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpEq))
		case "!=":
			ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpNotEq))
		case "<=":
			ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpLTEq))
		case ">=":
			ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpGTEq))
		case "implements":
			ccb.Code.AddInst(opcode.Implements, ccb.Linenum)
		}

	case *ast.CallExpression:
		ccb.Linenum = node.Token.Pos.Line
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compileMain(ccb, node.Arguments[i])
		}
		compileMain(ccb, node.Function)
		ccb.Code.AddInst(opcode.Call, ccb.Linenum, uint16(len(node.Arguments)))

	case *ast.ReturnStatement:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Value)
		ccb.Code.AddInst(opcode.Return, ccb.Linenum)

	case *ast.DefStatement:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Value)

		ccb.Code.AddInst(opcode.Define, ccb.Linenum,
			ccb.Locals.IndexOf(node.Name.Value),
			uint16(opcode.NewDefineFlag().WithConstant(node.Const).WithExport(node.Export)))

	case *ast.AssignStatement:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Value)

		if indexed, ok := node.Left.(*ast.IndexExpression); ok {
			compileMain(ccb, indexed.Index)
			compileMain(ccb, indexed.Left)
			ccb.Code.AddInst(opcode.StoreIndex, ccb.Linenum)
			break
		}

		if attrib, ok := node.Left.(*ast.AttributeExpression); ok {
			compileMain(ccb, attrib.Left)
			ccb.Code.AddInst(opcode.StoreAttribute, ccb.Linenum, ccb.Names.IndexOf(attrib.Index.String()))
			break
		}

		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			panic("Assignment to non ident or index")
		}

		if ccb.Locals.Contains(ident.Value) {
			ccb.Code.AddInst(opcode.StoreFast, ccb.Linenum, ccb.Locals.IndexOf(ident.Value))
		} else {
			ccb.Code.AddInst(opcode.StoreGlobal, ccb.Linenum, ccb.Names.IndexOf(ident.Value))
		}

	case *ast.DeleteStatement:
		ccb.Linenum = node.Token.Pos.Line
		ccb.Code.AddInst(opcode.DeleteFast, ccb.Linenum, ccb.Locals.IndexOf(node.Name))

	case *ast.IfExpression:
		compileIfStatement(ccb, node)

	case *ast.CompareExpression:
		compileCompareExpression(ccb, node)

	case *ast.ImportStatement:
		ccb.Linenum = node.Token.Pos.Line
		str := &object.String{Value: node.Path.Value}
		ccb.Code.AddInst(opcode.Import, ccb.Linenum, ccb.Constants.IndexOf(str))
		ccb.Code.AddInst(opcode.Define, ccb.Linenum, ccb.Locals.IndexOf(node.Name.Value), 0)

	case *ast.FunctionLiteral:
		compileFunction(ccb, node, false, false)

	case *ast.IndexExpression:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Left)
		compileMain(ccb, node.Index)
		ccb.Code.AddInst(opcode.LoadIndex, ccb.Linenum)

	case *ast.LoopStatement:
		compileLoop(ccb, node)

	case *ast.IterLoopStatement:
		compileIterLoop(ccb, node)

	case *ast.ContinueStatement:
		ccb.Linenum = node.Token.Pos.Line
		if !ccb.InLoop {
			panic("continue used in non-loop block")
		}
		ccb.Code.AddInst(opcode.Continue, ccb.Linenum)

	case *ast.BreakStatement:
		ccb.Linenum = node.Token.Pos.Line
		if !ccb.InLoop {
			panic("break used in non-loop block")
		}
		ccb.Code.AddInst(opcode.Break, ccb.Linenum)

	case *ast.ClassLiteral:
		compileClassLiteral(ccb, node)

	case *ast.NewInstance:
		ccb.Linenum = node.Token.Pos.Line
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compileMain(ccb, node.Arguments[i])
		}
		compileMain(ccb, node.Class)

		ccb.Code.AddInst(opcode.MakeInstance, ccb.Linenum, uint16(len(node.Arguments)))

	case *ast.AttributeExpression:
		ccb.Linenum = node.Token.Pos.Line
		compileMain(ccb, node.Left)
		ccb.Code.AddInst(opcode.LoadAttribute, ccb.Linenum, ccb.Names.IndexOf(node.Index.String()))

	case *ast.PassStatement:
		ccb.Linenum = node.Token.Pos.Line
		// Ignore

	case *ast.BreakpointStatement:
		ccb.Linenum = node.Token.Pos.Line
		ccb.Code.AddInst(opcode.Breakpoint, ccb.Linenum)

	// Not implemented yet
	case *ast.Program:
		panic("ast.Program Not implemented yet")

	default:
		panic(fmt.Sprintf("Node type not implemented: %T", node))
	}
}
