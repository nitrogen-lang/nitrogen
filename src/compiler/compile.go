package compiler

import (
	"bytes"
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/token"
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
		code:      new(bytes.Buffer),
		filename:  filename,
	}

	compile(ccb, node)
	if opcode.Opcode(ccb.code.Bytes()[ccb.code.Len()-1]) != opcode.Return {
		ccb.code.WriteByte(opcode.Return.ToByte())
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
		MaxBlockSize: calculateBlockSize(code),
	}

	return c
}

type maxsizer struct {
	max, current int
}

func (s *maxsizer) sub(delta int) {
	s.current -= delta
	if s.current > s.max { // Delta can be negative which would add to the size
		s.max = s.current
	}
}
func (s *maxsizer) add(delta int) {
	s.current += delta
	if s.current > s.max {
		s.max = s.current
	}
}

func calculateStackSize(c []byte) int {
	offset := 0
	stackSize := &maxsizer{}
	for offset < len(c) {
		code := opcode.Opcode(c[offset])
		offset++

		switch code {
		case opcode.LoadConst, opcode.LoadFast, opcode.LoadGlobal, opcode.StartTry:
			stackSize.add(1)
		case opcode.StoreIndex:
			stackSize.sub(3)
		case opcode.BinaryAdd, opcode.BinarySub, opcode.BinaryMul, opcode.BinaryDivide, opcode.BinaryMod, opcode.BinaryShiftL,
			opcode.BinaryShiftR, opcode.BinaryAnd, opcode.BinaryOr, opcode.BinaryNot, opcode.BinaryAndNot,
			opcode.StoreConst, opcode.StoreFast, opcode.Define, opcode.StoreGlobal, opcode.LoadIndex, opcode.Compare,
			opcode.Return, opcode.Pop, opcode.PopJumpIfTrue, opcode.PopJumpIfFalse, opcode.Throw:
			stackSize.sub(1)
		case opcode.Call:
			params := int(bytesToUint16(c[offset], c[offset+1]))
			stackSize.sub(params)
		case opcode.MakeArray:
			l := int(bytesToUint16(c[offset], c[offset+1]))
			stackSize.sub(l - 1)
		case opcode.BuildClass:
			l := int(bytesToUint16(c[offset], c[offset+1]))
			stackSize.sub(l + 2)
		case opcode.MakeMap:
			l := int(bytesToUint16(c[offset], c[offset+1]))
			stackSize.sub(l*2 - 1)
		case opcode.MakeFunction, opcode.StoreAttribute:
			stackSize.sub(2)
		}

		if opcode.HasOneByteArg[code] {
			offset++
		} else if opcode.HasTwoByteArg[code] {
			offset += 2
		} else if opcode.HasFourByteArg[code] {
			offset += 4
		}
	}
	return stackSize.max
}

func calculateBlockSize(c []byte) int {
	offset := 0
	blockLen := &maxsizer{}
	for offset < len(c) {
		code := opcode.Opcode(c[offset])
		offset++

		switch code {
		case opcode.StartLoop, opcode.StartTry:
			blockLen.add(1)
		case opcode.EndBlock:
			blockLen.sub(1)
		}

		if opcode.HasOneByteArg[code] {
			offset++
		} else if opcode.HasTwoByteArg[code] {
			offset += 2
		} else if opcode.HasFourByteArg[code] {
			offset += 4
		}
	}
	return blockLen.max
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

	// Literals
	case *ast.IntegerLiteral:
		i := &object.Integer{Value: node.Value}
		ccb.code.WriteByte(opcode.LoadConst.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(i)))
	case *ast.NullLiteral:
		compileLoadNull(ccb)
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		ccb.code.WriteByte(opcode.LoadConst.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(str)))
	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		ccb.code.WriteByte(opcode.LoadConst.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(float)))
	case *ast.Boolean:
		b := object.FalseConst
		if node.Value {
			b = object.TrueConst
		}
		ccb.code.WriteByte(opcode.LoadConst.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(b)))

	case *ast.Array:
		for _, e := range node.Elements {
			compile(ccb, e)
		}
		ccb.code.WriteByte(opcode.MakeArray.ToByte())
		ccb.code.Write(uint16ToBytes(uint16(len(node.Elements))))
	case *ast.HashLiteral:
		for k, v := range node.Pairs {
			compile(ccb, v)
			compile(ccb, k)
		}
		ccb.code.WriteByte(opcode.MakeMap.ToByte())
		ccb.code.Write(uint16ToBytes(uint16(len(node.Pairs))))

	// Expressions
	case *ast.Identifier:
		if ccb.locals.contains(node.Value) {
			ccb.code.WriteByte(opcode.LoadFast.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(node.Value)))
		} else {
			ccb.code.WriteByte(opcode.LoadGlobal.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.names.indexOf(node.Value)))
		}
	case *ast.PrefixExpression:
		compile(ccb, node.Right)

		switch node.Operator {
		case "!":
			ccb.code.WriteByte(opcode.UnaryNot.ToByte())
		case "-":
			ccb.code.WriteByte(opcode.UnaryNeg.ToByte())
		}
	case *ast.InfixExpression:
		compile(ccb, node.Left)
		compile(ccb, node.Right)

		switch node.Operator {
		case "+":
			ccb.code.WriteByte(opcode.BinaryAdd.ToByte())
		case "-":
			ccb.code.WriteByte(opcode.BinarySub.ToByte())
		case "*":
			ccb.code.WriteByte(opcode.BinaryMul.ToByte())
		case "/":
			ccb.code.WriteByte(opcode.BinaryDivide.ToByte())
		case "%":
			ccb.code.WriteByte(opcode.BinaryMod.ToByte())
		case "<<":
			ccb.code.WriteByte(opcode.BinaryShiftL.ToByte())
		case ">>":
			ccb.code.WriteByte(opcode.BinaryShiftR.ToByte())
		case "&":
			ccb.code.WriteByte(opcode.BinaryAnd.ToByte())
		case "&^":
			ccb.code.WriteByte(opcode.BinaryAndNot.ToByte())
		case "|":
			ccb.code.WriteByte(opcode.BinaryOr.ToByte())
		case "^":
			ccb.code.WriteByte(opcode.BinaryNot.ToByte())
		case "<":
			ccb.code.WriteByte(opcode.Compare.ToByte())
			ccb.code.WriteByte(opcode.CmpLT)
		case ">":
			ccb.code.WriteByte(opcode.Compare.ToByte())
			ccb.code.WriteByte(opcode.CmpGT)
		case "==":
			ccb.code.WriteByte(opcode.Compare.ToByte())
			ccb.code.WriteByte(opcode.CmpEq)
		case "!=":
			ccb.code.WriteByte(opcode.Compare.ToByte())
			ccb.code.WriteByte(opcode.CmpNotEq)
		case "<=":
			ccb.code.WriteByte(opcode.Compare.ToByte())
			ccb.code.WriteByte(opcode.CmpLTEq)
		case ">=":
			ccb.code.WriteByte(opcode.Compare.ToByte())
			ccb.code.WriteByte(opcode.CmpGTEq)
		}
	case *ast.CallExpression:
		for i := len(node.Arguments) - 1; i >= 0; i-- {
			compile(ccb, node.Arguments[i])
		}
		compile(ccb, node.Function)
		ccb.code.WriteByte(opcode.Call.ToByte())
		ccb.code.Write(uint16ToBytes(uint16(len(node.Arguments))))
	case *ast.ReturnStatement:
		compile(ccb, node.Value)
		ccb.code.WriteByte(opcode.Return.ToByte())
	case *ast.DefStatement:
		compile(ccb, node.Value)

		if node.Const {
			ccb.code.WriteByte(opcode.StoreConst.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(node.Name.Value)))
		} else {
			ccb.code.WriteByte(opcode.Define.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(node.Name.Value)))
		}
		compileLoadNull(ccb)
	case *ast.AssignStatement:
		compile(ccb, node.Value)

		if indexed, ok := node.Left.(*ast.IndexExpression); ok {
			compile(ccb, indexed.Index)
			compile(ccb, indexed.Left)
			ccb.code.WriteByte(opcode.StoreIndex.ToByte())
			compileLoadNull(ccb)
			break
		}

		if attrib, ok := node.Left.(*ast.AttributeExpression); ok {
			compile(ccb, attrib.Left)
			ccb.code.WriteByte(opcode.StoreAttribute.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.names.indexOf(attrib.Index.Value)))
			compileLoadNull(ccb)
			break
		}

		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			panic("Assignment to non ident or index")
		}

		if ccb.locals.contains(ident.Value) {
			ccb.code.WriteByte(opcode.StoreFast.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(ident.Value)))
		} else {
			ccb.code.WriteByte(opcode.StoreGlobal.ToByte())
			ccb.code.Write(uint16ToBytes(ccb.names.indexOf(ident.Value)))
		}
		compileLoadNull(ccb)
	case *ast.IfExpression:
		compileIfStatement(ccb, node)
	case *ast.CompareExpression:
		compileCompareExpression(ccb, node)

	case *ast.FunctionLiteral:
		compileFunction(ccb, node, false, false)

	case *ast.IndexExpression:
		compile(ccb, node.Index)
		compile(ccb, node.Left)
		ccb.code.WriteByte(opcode.LoadIndex.ToByte())

	case *ast.ForLoopStatement:
		compileLoop(ccb, node)
	case *ast.ContinueStatement:
		ccb.code.WriteByte(opcode.Continue.ToByte())
	case *ast.BreakStatement:
		ccb.code.WriteByte(opcode.Break.ToByte())

	case *ast.TryCatchExpression:
		compileTryCatch(ccb, node)
	case *ast.ThrowStatement:
		compile(ccb, node.Expression)
		ccb.code.WriteByte(opcode.Throw.ToByte())
	case *ast.ClassLiteral:
		compileClassLiteral(ccb, node)
	case *ast.MakeInstance:
		for _, p := range node.Arguments {
			compile(ccb, p)
		}
		compile(ccb, node.Class)

		ccb.code.WriteByte(opcode.MakeInstance.ToByte())
		ccb.code.Write(uint16ToBytes(uint16(len(node.Arguments))))
	case *ast.AttributeExpression:
		compile(ccb, node.Left)
		ccb.code.WriteByte(opcode.LoadAttribute.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.names.indexOf(node.Index.Value)))
	case *ast.PassStatement:
		// Ignore

	// Not implemented yet
	case *ast.Program:
		panic("ast.Program Not implemented yet")
	default:
		panic(fmt.Sprintf("Node type not implemented: %T", node))
	}
}

func compileClassLiteral(ccb *codeBlockCompiler, class *ast.ClassLiteral) {
	for _, f := range class.Methods {
		f.FQName = fmt.Sprintf("%s.%s", class.Name, f.Name)
		compileFunction(ccb, f, true, class.Parent != "")
	}

	ccb2 := &codeBlockCompiler{
		constants: newConstantTable(),
		locals:    newStringTable(),
		names:     newStringTable(),
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
	}

	for _, f := range class.Fields {
		compile(ccb2, f)
	}
	compileLoadNull(ccb2)
	ccb2.code.WriteByte(opcode.Return.ToByte())

	code := ccb2.code.Bytes()
	props := &CodeBlock{
		Name:         fmt.Sprintf("%s.__init", class.Name),
		Filename:     ccb.filename,
		LocalCount:   len(ccb2.locals.table),
		Code:         code,
		Constants:    ccb2.constants.table,
		Names:        ccb2.names.table,
		Locals:       ccb2.locals.table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
	}

	ccb.code.WriteByte(opcode.LoadConst.ToByte())
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(props)))

	if class.Parent == "" {
		compileLoadNull(ccb)
	} else {
		compile(ccb, &ast.Identifier{Value: class.Parent})
	}

	ccb.code.WriteByte(opcode.LoadConst.ToByte())
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.MakeStringObj(class.Name))))

	ccb.code.WriteByte(opcode.BuildClass.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(len(class.Methods))))
}

func compileTryCatch(ccb *codeBlockCompiler, try *ast.TryCatchExpression) {
	ccb.code.WriteByte(opcode.PrepareBlock.ToByte())

	_, tryNoNil := try.Try.Statements[len(try.Try.Statements)-1].(*ast.ExpressionStatement)
	_, catchNoNil := try.Catch.Statements[len(try.Catch.Statements)-1].(*ast.ExpressionStatement)

	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = mainCode.Len() + ccb.offset
	ccb.code = new(bytes.Buffer)
	compile(ccb, try.Try)
	tryBlock := ccb.code

	// 6 = 2 opcodes + 2 x 2 byte args (START_TRY and JUMP_FORWARD)
	catchBlockLoc := mainCode.Len() + tryBlock.Len() + 6
	ccb.offset = catchBlockLoc
	ccb.code = new(bytes.Buffer)
	if try.Symbol != nil {
		ccb.locals.indexOf(try.Symbol.Value)
	}
	compile(ccb, try.Catch)
	catchBlock := ccb.code

	ccb.code = mainCode
	ccb.offset = oldOffset

	catchSymOffset := 1 // Catch doesn't bind exception, just pops it
	if try.Symbol != nil {
		catchSymOffset = 3 // Catch binds exception to a variable
	}
	if tryNoNil && !catchNoNil {
		catchSymOffset += 3 // Skip load nil from catch block
	}
	if !tryNoNil && catchNoNil {
		if try.Symbol == nil {
			catchSymOffset += 3 // Skip load nil from try block
		} else {
			catchSymOffset += 6 // Skip load nil catch block and symbol bind
		}
	}

	ccb.code.WriteByte(opcode.StartTry.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(catchBlockLoc)))
	ccb.code.Write(tryBlock.Bytes())

	ccb.code.WriteByte(opcode.JumpForward.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(catchBlock.Len() + catchSymOffset)))

	if try.Symbol == nil {
		ccb.code.WriteByte(opcode.Pop.ToByte())
	} else {
		ccb.code.WriteByte(opcode.Define.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(try.Symbol.Value)))
	}

	ccb.code.Write(catchBlock.Bytes())
	if catchNoNil && !tryNoNil {
		ccb.code.WriteByte(opcode.JumpForward.ToByte())
		ccb.code.Write(uint16ToBytes(3))
	}
	if !tryNoNil || !catchNoNil {
		compileLoadNull(ccb)
	}
	ccb.code.WriteByte(opcode.EndBlock.ToByte())
}

func compileBlock(ccb *codeBlockCompiler, block *ast.BlockStatement) {
	l := len(block.Statements) - 1
	for i, s := range block.Statements {
		compile(ccb, s)
		if i < l {
			switch s.(type) {
			case *ast.ExpressionStatement:
				ccb.code.WriteByte(opcode.Pop.ToByte())
			case *ast.DefStatement:
				ccb.code.WriteByte(opcode.Pop.ToByte())
			}
		}
	}
}

func compileFunction(ccb *codeBlockCompiler, fn *ast.FunctionLiteral, inClass bool, hasParent bool) {
	ccb2 := &codeBlockCompiler{
		constants: newConstantTable(),
		locals:    newStringTable(),
		names:     newStringTable(),
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
	}

	for _, p := range fn.Parameters {
		ccb2.locals.indexOf(p.Value)
	}
	if inClass {
		ccb2.locals.indexOf("this")
		if hasParent {
			ccb2.locals.indexOf("parent")
		}
	}

	compile(ccb2, fn.Body)

	if len(fn.Body.Statements) > 0 {
		switch fn.Body.Statements[len(fn.Body.Statements)-1].(type) {
		case *ast.ExpressionStatement:
			break
		case *ast.ReturnStatement:
			break
		default:
			compileLoadNull(ccb2)
		}

		if opcode.Opcode(ccb2.code.Bytes()[ccb2.code.Len()-1]) != opcode.Return {
			ccb2.code.WriteByte(opcode.Return.ToByte())
		}
	} else {
		compileLoadNull(ccb2)
		ccb2.code.WriteByte(opcode.Return.ToByte())
	}

	code := ccb2.code.Bytes()
	body := &CodeBlock{
		Name:         fn.FQName,
		Filename:     ccb.filename,
		LocalCount:   len(ccb2.locals.table),
		Code:         code,
		Constants:    ccb2.constants.table,
		Names:        ccb2.names.table,
		Locals:       ccb2.locals.table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
	}

	ccb.code.WriteByte(opcode.LoadConst.ToByte())
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(body)))

	for _, p := range fn.Parameters {
		ccb.code.WriteByte(opcode.LoadConst.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.MakeStringObj(p.Value))))
	}
	ccb.code.WriteByte(opcode.MakeArray.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(len(fn.Parameters))))

	ccb.code.WriteByte(opcode.LoadConst.ToByte())
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.MakeStringObj(fn.Name))))

	ccb.code.WriteByte(opcode.MakeFunction.ToByte())
}

func compileInnerBlock(ccb *codeBlockCompiler, node ast.Node) *bytes.Buffer {
	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = ccb.code.Len() + ccb.offset
	ccb.code = new(bytes.Buffer)
	compile(ccb, node)
	block := ccb.code

	ccb.code = mainCode
	ccb.offset = oldOffset
	return block
}

func compileIfStatement(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	if ifs.Alternative == nil {
		compileIfStatementNoElse(ccb, ifs)
		return
	}

	_, trueNoNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	_, falseNoNil := ifs.Alternative.Statements[len(ifs.Alternative.Statements)-1].(*ast.ExpressionStatement)

	compile(ccb, ifs.Condition)

	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = ccb.code.Len() + ccb.offset
	ccb.code = new(bytes.Buffer)
	compile(ccb, ifs.Consequence)
	trueBranch := ccb.code

	// 1 = 1 opcode
	falseBranchLoc := mainCode.Len() + trueBranch.Len() + ccb.offset + 1
	if trueNoNil {
		// 3 = 1 opcode + 2 byte arg (implicit nil from true branch)
		falseBranchLoc -= 3
	}
	ccb.offset = falseBranchLoc
	ccb.code = new(bytes.Buffer)
	compile(ccb, ifs.Alternative)
	falseBranch := ccb.code

	ccb.code = mainCode
	ccb.offset = oldOffset

	ccb.code.WriteByte(opcode.PopJumpIfFalse.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(falseBranchLoc)))
	ccb.code.Write(trueBranch.Bytes())
	if !trueNoNil {
		compileLoadNull(ccb)
	}
	ccb.code.WriteByte(opcode.JumpForward.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(falseBranch.Len())))
	ccb.code.Write(falseBranch.Bytes())
	if !falseNoNil {
		compileLoadNull(ccb)
	}
}

func compileIfStatementNoElse(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	compile(ccb, ifs.Condition)

	trueBranch := compileInnerBlock(ccb, ifs.Consequence)

	_, noNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)

	// 6 = 2 opcodes + 2 x 2 byte args
	afterIfStmt := ccb.code.Len() + trueBranch.Len() + ccb.offset + 6
	if !noNil {
		// 3 = 1 opcode + 2 byte arg
		afterIfStmt -= 3
	}

	ccb.code.WriteByte(opcode.PopJumpIfFalse.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(afterIfStmt)))
	ccb.code.Write(trueBranch.Bytes())
	if noNil {
		ccb.code.WriteByte(opcode.JumpForward.ToByte())
		ccb.code.Write(uint16ToBytes(uint16(3))) // 3 = 1 opcode + 2 byte arg (for implicit nil)
	}
	compileLoadNull(ccb)
}

func compileLoadNull(ccb *codeBlockCompiler) {
	ccb.code.WriteByte(opcode.LoadConst.ToByte())
	ccb.code.Write(uint16ToBytes(ccb.constants.indexOf(object.NullConst)))
}

func compileCompareExpression(ccb *codeBlockCompiler, cmp *ast.CompareExpression) {
	compile(ccb, cmp.Left)

	cntBranch := compileInnerBlock(ccb, cmp.Right)

	// 3 = 1 opcode + 1 x 2 byte arg
	afterCompare := ccb.code.Len() + cntBranch.Len() + ccb.offset + 3

	if cmp.Token.Type == token.LAnd {
		ccb.code.WriteByte(opcode.JumpIfFalseOrPop.ToByte())
	} else {
		ccb.code.WriteByte(opcode.JumpIfTrueOrPop.ToByte())
	}

	ccb.code.Write(uint16ToBytes(uint16(afterCompare)))
	ccb.code.Write(cntBranch.Bytes())
}

func compileLoop(ccb *codeBlockCompiler, loop *ast.ForLoopStatement) {
	if loop.Init == nil {
		compileInfiniteLoop(ccb, loop)
		return
	}

	ccb.code.WriteByte(opcode.PrepareBlock.ToByte())
	compile(ccb, loop.Init)

	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = ccb.code.Len() + ccb.offset
	ccb.code = new(bytes.Buffer)
	compile(ccb, loop.Condition)
	condition := ccb.code

	// 8 = 2 x opcode + 3 x 2 byte args
	ccb.offset = mainCode.Len() + condition.Len() + oldOffset + 8
	ccb.code = new(bytes.Buffer)
	compile(ccb, loop.Body)
	loopBody := ccb.code

	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		ccb.code.WriteByte(opcode.Pop.ToByte())
	}

	// 3 = 1 opcode + 2 byte arg
	ccb.offset = mainCode.Len() + condition.Len() + loopBody.Len() + ccb.offset + 3
	ccb.code = new(bytes.Buffer)
	compile(ccb, loop.Iter)
	iterator := ccb.code

	ccb.code = mainCode
	ccb.offset = oldOffset

	// 10 = 4 opcodes + 3 x 2 byte args
	endBlock := mainCode.Len() + condition.Len() + loopBody.Len() + iterator.Len() + ccb.offset + 10
	// 8 = 2 opcode + 3 x 2 byte args
	iterBlock := mainCode.Len() + condition.Len() + loopBody.Len() + ccb.offset + 8
	ccb.code.WriteByte(opcode.StartLoop.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))
	ccb.code.Write(uint16ToBytes(uint16(iterBlock)))

	ccb.code.Write(condition.Bytes())
	ccb.code.WriteByte(opcode.PopJumpIfFalse.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))

	ccb.code.Write(loopBody.Bytes())
	ccb.code.Write(iterator.Bytes())

	ccb.code.WriteByte(opcode.Pop.ToByte())
	ccb.code.WriteByte(opcode.NextIter.ToByte())
	ccb.code.WriteByte(opcode.EndBlock.ToByte())
}

func compileInfiniteLoop(ccb *codeBlockCompiler, loop *ast.ForLoopStatement) {
	loopBody := compileInnerBlock(ccb, loop.Body)
	loopBody.WriteByte(opcode.Continue.ToByte())

	// 3 = 1 opcode + 1 x 2 byte arg
	loopEnd := ccb.code.Len() + loopBody.Len() + 3
	ccb.code.WriteByte(opcode.StartLoop.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(loopEnd)))

	ccb.code.Write(loopBody.Bytes())
}
