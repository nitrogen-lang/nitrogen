package compiler

import (
	"bytes"
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/token"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func compileClassLiteral(ccb *codeBlockCompiler, class *ast.ClassLiteral) {
	for _, f := range class.Methods {
		f.FQName = fmt.Sprintf("%s.%s", class.Name, f.Name)
		compileFunction(ccb, f, true, class.Parent != "")
	}

	ccb2 := &codeBlockCompiler{
		constants: newConstantTable(),
		locals:    newStringTable(),
		names:     newStringTable(),
		code:      NewInstSet(),
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
		catchSymOffset = 6 // Catch binds exception to a variable
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

	ccb.code.WriteByte(opcode.JumpForward.ToByte()) // No exception, jump past catch block
	ccb.code.Write(uint16ToBytes(uint16(catchBlock.Len() + catchSymOffset)))

	if try.Symbol == nil {
		ccb.code.WriteByte(opcode.Pop.ToByte())
	} else {
		ccb.code.WriteByte(opcode.Define.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(try.Symbol.Value)))
	}

	ccb.code.Write(catchBlock.Bytes())
	if try.Symbol != nil {
		ccb.code.WriteByte(opcode.DeleteFast.ToByte())
		ccb.code.Write(uint16ToBytes(ccb.locals.indexOf(try.Symbol.Value)))
	}
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
				ccb.code.addInst(opcode.Pop)
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
	ccb2.locals.indexOf("arguments") // `arguments` holds any remaining arguments from a function call
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

func compileInnerBlock(ccb *codeBlockCompiler, node ast.Node, extraOffset int) *bytes.Buffer {
	mainCode := ccb.code
	oldOffset := ccb.offset

	ccb.offset = ccb.code.Len() + ccb.offset + extraOffset
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

	compile(ccb, ifs.Condition)

	_, trueNoNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	_, falseNoNil := ifs.Alternative.Statements[len(ifs.Alternative.Statements)-1].(*ast.ExpressionStatement)

	trueBranch := compileInnerBlock(ccb, ifs.Consequence, 3)

	falseBranchLoc := ccb.code.Len() + trueBranch.Len() + ccb.offset + 6
	if !trueNoNil {
		// 3 = 1 opcode + 2 byte arg (implicit nil from true branch)
		falseBranchLoc += 3
	}

	extraOffset := trueBranch.Len() + 6
	if !trueNoNil {
		extraOffset += 3
	}
	falseBranch := compileInnerBlock(ccb, ifs.Alternative, extraOffset)

	jmpForw := falseBranch.Len()
	if !falseNoNil {
		jmpForw += 3
	}

	ccb.code.WriteByte(opcode.PopJumpIfFalse.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(falseBranchLoc)))
	ccb.code.Write(trueBranch.Bytes())
	if !trueNoNil {
		compileLoadNull(ccb)
	}

	ccb.code.WriteByte(opcode.JumpForward.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(jmpForw)))
	ccb.code.Write(falseBranch.Bytes())
	if !falseNoNil {
		compileLoadNull(ccb)
	}
}

func compileIfStatementNoElse(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	compile(ccb, ifs.Condition)

	// 3 to compensate for implicit nil return
	trueBranch := compileInnerBlock(ccb, ifs.Consequence, 3)

	_, noNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)

	// 6 = 2 opcodes + 2 x 2 byte args
	afterIfStmt := ccb.code.Len() + trueBranch.Len() + ccb.offset + 3
	if noNil {
		// 3 = 1 opcode + 2 byte arg
		afterIfStmt += 3
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
	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.NullConst))
}

func compileCompareExpression(ccb *codeBlockCompiler, cmp *ast.CompareExpression) {
	compile(ccb, cmp.Left)

	afterCompareLabel := randomLabel("cmp_")

	if cmp.Token.Type == token.LAnd {
		ccb.code.addLabeledArgs(opcode.JumpIfFalseOrPop, args(0), argsl(afterCompareLabel))
	} else {
		ccb.code.addLabeledArgs(opcode.JumpIfTrueOrPop, args(0), argsl(afterCompareLabel))
	}

	compile(ccb, cmp.Right)
	ccb.code.addLabel(afterCompareLabel)
}

func compileLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	if loop.Init == nil {
		if loop.Condition == nil {
			compileInfiniteLoop(ccb, loop)
		} else {
			compileWhileLoop(ccb, loop)
		}
		return
	}

	// A loop begins with a PREPARE_BLOCK opcode this creates the first layer environment
	ccb.code.WriteByte(opcode.OpenScope.ToByte())
	// Initialization is done in this first layer
	compile(ccb, loop.Init)

	// The loop operates in an isolated environment so locals need to be handled carefully
	condCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
		offset:    ccb.code.Len() + ccb.offset,
	}

	// Compile the loop's condition check code
	compile(condCCB, loop.Condition)
	condition := condCCB.code

	// Prepare for main body
	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
		// 8 = 2 x opcode + 3 x 2 byte args
		offset: ccb.code.Len() + condition.Len() + ccb.offset + 8,
	}

	// Compile main body of loop
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.WriteByte(opcode.Pop.ToByte())
	}
	loopBody := bodyCCB.code

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	// Prepare for iteration code
	iterCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
		// 3 = 1 opcode + 2 byte arg
		offset: ccb.code.Len() + condition.Len() + loopBody.Len() + ccb.offset + 3,
	}

	// Compile iteration
	compile(iterCCB, loop.Iter)
	iterator := iterCCB.code

	// Again, copy over the locals for indexing
	for _, n := range iterCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	// Generate and build the full loop code

	// 9 = 3 opcodes + 3 x 2 byte args
	endBlock := ccb.code.Len() + condition.Len() + loopBody.Len() + iterator.Len() + ccb.offset + 9
	// 8 = 2 opcode + 3 x 2 byte args
	iterBlock := ccb.code.Len() + condition.Len() + loopBody.Len() + ccb.offset + 8

	ccb.code.WriteByte(opcode.StartLoop.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))
	ccb.code.Write(uint16ToBytes(uint16(iterBlock)))

	ccb.code.Write(condition.Bytes())
	ccb.code.WriteByte(opcode.PopJumpIfFalse.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))

	ccb.code.Write(loopBody.Bytes())
	ccb.code.Write(iterator.Bytes())

	ccb.code.WriteByte(opcode.NextIter.ToByte())
	ccb.code.WriteByte(opcode.EndBlock.ToByte())
	ccb.code.WriteByte(opcode.CloseScope.ToByte())
	ccb.code.WriteByte(opcode.CloseScope.ToByte())
}

func compileInfiniteLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	// loopBody := compileInnerBlock(ccb, loop.Body, 0)
	// loopBody.WriteByte(opcode.Continue.ToByte())

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
		// 8 = 2 x opcode + 3 x 2 byte args
		offset: ccb.code.Len() + ccb.offset + 5,
	}
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.WriteByte(opcode.Pop.ToByte())
	}
	loopBody := bodyCCB.code

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	// 9 = 3 opcodes + 3 x 2 byte args
	endBlock := ccb.code.Len() + loopBody.Len() + ccb.offset + 6
	// 8 = 2 opcode + 3 x 2 byte args
	iterBlock := ccb.code.Len() + loopBody.Len() + ccb.offset + 5

	ccb.code.WriteByte(opcode.StartLoop.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))
	ccb.code.Write(uint16ToBytes(uint16(iterBlock)))

	ccb.code.Write(loopBody.Bytes())

	ccb.code.WriteByte(opcode.NextIter.ToByte())
	ccb.code.WriteByte(opcode.EndBlock.ToByte())
	ccb.code.WriteByte(opcode.CloseScope.ToByte())
}

func compileWhileLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	// loopBody := compileInnerBlock(ccb, loop.Body, 0)
	// loopBody.WriteByte(opcode.Continue.ToByte())

	condCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
		offset:    ccb.code.Len() + ccb.offset,
	}

	// Compile the loop's condition check code
	compile(condCCB, loop.Condition)
	condition := condCCB.code

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      new(bytes.Buffer),
		filename:  ccb.filename,
		offset:    ccb.code.Len() + ccb.offset + condition.Len() + 12,
	}
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.WriteByte(opcode.Pop.ToByte())
	}
	loopBody := bodyCCB.code

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	// 9 = 3 opcodes + 3 x 2 byte args
	endBlock := ccb.code.Len() + condition.Len() + loopBody.Len() + ccb.offset + 9
	// 8 = 2 opcode + 3 x 2 byte args
	iterBlock := ccb.code.Len() + condition.Len() + loopBody.Len() + ccb.offset + 8

	ccb.code.WriteByte(opcode.StartLoop.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))
	ccb.code.Write(uint16ToBytes(uint16(iterBlock)))

	ccb.code.Write(condition.Bytes())
	ccb.code.WriteByte(opcode.PopJumpIfFalse.ToByte())
	ccb.code.Write(uint16ToBytes(uint16(endBlock)))

	ccb.code.Write(loopBody.Bytes())

	ccb.code.WriteByte(opcode.NextIter.ToByte())
	ccb.code.WriteByte(opcode.EndBlock.ToByte())
	ccb.code.WriteByte(opcode.CloseScope.ToByte())
}
