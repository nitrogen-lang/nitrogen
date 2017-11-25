package compiler

import (
	"bytes"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/token"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func Compile(tree *ast.Program) *CodeBlock {
	return compileFrame(&ast.BlockStatement{Statements: tree.Statements}, "<module>")
}

func compileFrame(node ast.Node, name string) *CodeBlock {
	ccb := &codeBlockCompiler{
		constants: newConstantTable(),
		locals:    newStringTable(),
		names:     newStringTable(),
		code:      new(bytes.Buffer),
	}

	compile(ccb, node)
	if ccb.code.Bytes()[ccb.code.Len()-1] != opcode.Return {
		ccb.code.WriteByte(opcode.Return)
	}

	filename := ""
	if program, ok := node.(*ast.Program); ok {
		filename = program.Filename
	}

	code := ccb.code.Bytes()
	c := &CodeBlock{
		Name:         name,
		Filename:     filename,
		LocalCount:   len(ccb.locals.table),
		Code:         code,
		Constants:    ccb.constants.table,
		Names:        ccb.names.table,
		Locals:       ccb.locals.table,
		MaxStackSize: calculateStackSize(code),
	}

	return c
}

func calculateStackSize(c []byte) int {
	return 0
}

func compile(ccb *codeBlockCompiler, node ast.Node) {
	switch node := node.(type) {
	case *ast.ExpressionStatement:
		compile(ccb, node.Expression)
		if _, ok := node.Expression.(*ast.CallExpression); ok {
			ccb.code.WriteByte(opcode.Pop)
		}
	case *ast.BlockStatement:
		compileBlock(ccb, node)

	// Literals
	case *ast.IntegerLiteral:
		i := &object.Integer{Value: node.Value}
		ccb.code.WriteByte(opcode.LoadConst)
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(i)))
	case *ast.NullLiteral:
		ccb.code.WriteByte(opcode.LoadConst)
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.NullConst)))
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		ccb.code.WriteByte(opcode.LoadConst)
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(str)))
	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		ccb.code.WriteByte(opcode.LoadConst)
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(float)))
	case *ast.Boolean:
		b := object.FalseConst
		if node.Value {
			b = object.TrueConst
		}
		ccb.code.WriteByte(opcode.LoadConst)
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(b)))

	case *ast.Array:
		for _, e := range node.Elements {
			compile(ccb, e)
		}
		ccb.code.WriteByte(opcode.MakeArray)
		ccb.code.Write(uint16ToBytes(uint16(len(node.Elements))))
	case *ast.HashLiteral:
		for k, v := range node.Pairs {
			compile(ccb, v)
			compile(ccb, k)
		}
		ccb.code.WriteByte(opcode.MakeMap)
		ccb.code.Write(uint16ToBytes(uint16(len(node.Pairs))))

	// Expressions
	case *ast.Identifier:
		if ccb.locals.contains(node.Value) {
			ccb.code.WriteByte(opcode.LoadFast)
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(node.Value)))
		} else {
			ccb.code.WriteByte(opcode.LoadGlobal)
			ccb.code.Write(uint16ToBytes(ccb.names.indexOf(node.Value)))
		}
	case *ast.PrefixExpression:
		compile(ccb, node.Right)

		switch node.Operator {
		case "!":
			ccb.code.WriteByte(opcode.UnaryNot)
		case "-":
			ccb.code.WriteByte(opcode.UnaryNeg)
		}
	case *ast.InfixExpression:
		compile(ccb, node.Left)
		compile(ccb, node.Right)

		switch node.Operator {
		case "+":
			ccb.code.WriteByte(opcode.BinaryAdd)
		case "-":
			ccb.code.WriteByte(opcode.BinarySub)
		case "*":
			ccb.code.WriteByte(opcode.BinaryMul)
		case "/":
			ccb.code.WriteByte(opcode.BinaryDivide)
		case "%":
			ccb.code.WriteByte(opcode.BinaryMod)
		case "<<":
			ccb.code.WriteByte(opcode.BinaryShiftL)
		case ">>":
			ccb.code.WriteByte(opcode.BinaryShiftR)
		case "&":
			ccb.code.WriteByte(opcode.BinaryAnd)
		case "&^":
			ccb.code.WriteByte(opcode.BinaryAndNot)
		case "|":
			ccb.code.WriteByte(opcode.BinaryOr)
		case "^":
			ccb.code.WriteByte(opcode.BinaryNot)
		case "<":
			ccb.code.WriteByte(opcode.Compare)
			ccb.code.WriteByte(opcode.CmpLT)
		case ">":
			ccb.code.WriteByte(opcode.Compare)
			ccb.code.WriteByte(opcode.CmpGT)
		case "==":
			ccb.code.WriteByte(opcode.Compare)
			ccb.code.WriteByte(opcode.CmpEq)
		case "!=":
			ccb.code.WriteByte(opcode.Compare)
			ccb.code.WriteByte(opcode.CmpNotEq)
		case "<=":
			ccb.code.WriteByte(opcode.Compare)
			ccb.code.WriteByte(opcode.CmpLTEq)
		case ">=":
			ccb.code.WriteByte(opcode.Compare)
			ccb.code.WriteByte(opcode.CmpGTEq)
		}
	case *ast.CallExpression:
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compile(ccb, node.Arguments[i])
		}
		compile(ccb, node.Function)
		ccb.code.WriteByte(opcode.Call)
		ccb.code.Write(uint16ToBytes(uint16(len(node.Arguments))))
	case *ast.Program:
		panic("Not implemented yet")
	case *ast.ReturnStatement:
		compile(ccb, node.Value)
		ccb.code.WriteByte(opcode.Return)
	case *ast.DefStatement:
		compile(ccb, node.Value)

		if node.Const {
			ccb.code.WriteByte(opcode.StoreConst)
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(node.Name.Value)))
		} else {
			ccb.code.WriteByte(opcode.StoreFast)
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(node.Name.Value)))
		}
	case *ast.AssignStatement:
		compile(ccb, node.Value)

		if indexed, ok := node.Left.(*ast.IndexExpression); ok {
			compile(ccb, indexed.Index)
			compile(ccb, indexed.Left)
			ccb.code.WriteByte(opcode.StoreIndex)
			break
		}

		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			panic("Assignment to non ident or index")
		}

		if ccb.locals.contains(ident.Value) {
			ccb.code.WriteByte(opcode.StoreFast)
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(ident.Value)))
		} else {
			ccb.code.WriteByte(opcode.StoreGlobal)
			ccb.code.Write(uint16ToBytes(ccb.names.indexOf(ident.Value)))
		}
	case *ast.IfExpression:
		compileIfStatement(ccb, node)
	case *ast.CompareExpression:
		compileCompareExpression(ccb, node)

	case *ast.FunctionLiteral:
		compileFunction(ccb, node)

	case *ast.IndexExpression:
		compile(ccb, node.Index)
		compile(ccb, node.Left)
		ccb.code.WriteByte(opcode.LoadIndex)

	// Not implemented yet
	case *ast.ThrowStatement:
		panic("Not implemented yet")
	case *ast.ForLoopStatement:
		panic("Not implemented yet")
	case *ast.ContinueStatement:
		panic("Not implemented yet")
	case *ast.BreakStatement:
		panic("Not implemented yet")
	case *ast.TryCatchExpression:
		panic("Not implemented yet")
	case *ast.ClassLiteral:
		panic("Not implemented yet")
	case *ast.MakeInstance:
		panic("Not implemented yet")
	}
}

func compileBlock(ccb *codeBlockCompiler, block *ast.BlockStatement) {
	for _, s := range block.Statements {
		compile(ccb, s)
	}
}

func compileFunction(ccb *codeBlockCompiler, fn *ast.FunctionLiteral) {
	ccb2 := &codeBlockCompiler{
		constants: newConstantTable(),
		locals:    newStringTable(),
		names:     newStringTable(),
		code:      new(bytes.Buffer),
	}

	for _, p := range fn.Parameters {
		ccb2.locals.indexOf(p.Value)
	}

	compile(ccb2, fn.Body)
	if ccb2.code.Bytes()[ccb2.code.Len()-1] != opcode.Return {
		ccb2.code.WriteByte(opcode.Return)
	}

	code := ccb2.code.Bytes()
	body := &CodeBlock{
		Name:         fn.Name,
		Filename:     "",
		LocalCount:   len(ccb2.locals.table),
		Code:         code,
		Constants:    ccb2.constants.table,
		Names:        ccb2.names.table,
		Locals:       ccb2.locals.table,
		MaxStackSize: calculateStackSize(code),
	}

	ccb.code.WriteByte(opcode.LoadConst)
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(body)))

	for _, p := range fn.Parameters {
		ccb.code.WriteByte(opcode.LoadConst)
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.MakeStringObj(p.Value))))
	}
	ccb.code.WriteByte(opcode.MakeArray)
	ccb.code.Write(uint16ToBytes(uint16(len(fn.Parameters))))

	ccb.code.WriteByte(opcode.LoadConst)
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.MakeStringObj(fn.Name))))

	ccb.code.WriteByte(opcode.MakeFunction)
}

func compileIfStatement(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	if ifs.Alternative == nil {
		compileIfStatementNoElse(ccb, ifs)
		return
	}

	compile(ccb, ifs.Condition)

	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = ccb.code.Len() + ccb.offset
	ccb.code = new(bytes.Buffer)
	compile(ccb, ifs.Consequence)
	trueBranch := ccb.code

	// Prior code to if statement + size of true branch + faked offset + 3 (other opcodes)
	falseBranchLoc := mainCode.Len() + trueBranch.Len() + ccb.offset + 4
	ccb.offset = falseBranchLoc
	ccb.code = new(bytes.Buffer)
	compile(ccb, ifs.Alternative)
	falseBranch := ccb.code

	ccb.code = mainCode
	ccb.offset = oldOffset

	ccb.code.WriteByte(opcode.PopJumpIfFalse)
	ccb.code.Write(uint16ToBytes(uint16(falseBranchLoc)))
	ccb.code.Write(trueBranch.Bytes())
	ccb.code.WriteByte(opcode.Pop)
	ccb.code.WriteByte(opcode.JumpForward)
	ccb.code.Write(uint16ToBytes(uint16(falseBranch.Len())))
	ccb.code.Write(falseBranch.Bytes())
}

func compileIfStatementNoElse(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	compile(ccb, ifs.Condition)

	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = ccb.code.Len() + ccb.offset
	ccb.code = new(bytes.Buffer)
	compile(ccb, ifs.Consequence)
	trueBranch := ccb.code

	ccb.code = mainCode
	ccb.offset = oldOffset

	// Prior code to if statement + size of true branch + faked offset - 2 (IDK why 2, it just works)
	afterIfStmt := ccb.code.Len() + trueBranch.Len() + ccb.offset + 3

	ccb.code.WriteByte(opcode.PopJumpIfFalse)
	ccb.code.Write(uint16ToBytes(uint16(afterIfStmt)))
	ccb.code.Write(trueBranch.Bytes())
	ccb.code.WriteByte(opcode.Pop)
}

func compileCompareExpression(ccb *codeBlockCompiler, cmp *ast.CompareExpression) {
	compile(ccb, cmp.Left)

	rightBranchLoc := ccb.code.Len() + 3

	if cmp.Token.Type == token.LAnd {
		ccb.code.WriteByte(opcode.JumpIfFalseOrPop)
	} else {
		ccb.code.WriteByte(opcode.JumpIfTrueOrPop)
	}

	ccb.code.Write(uint16ToBytes(uint16(rightBranchLoc)))

	compile(ccb, cmp.Right)
}
