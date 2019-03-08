package compiler

import (
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
		name:      ccb.name,
		inLoop:    ccb.inLoop,
	}

	for _, f := range class.Fields {
		compile(ccb2, f)
	}
	compileLoadNull(ccb2)
	ccb2.code.addInst(opcode.Return)

	code := ccb2.code
	props := &CodeBlock{
		Name:         fmt.Sprintf("%s.__init", class.Name),
		Filename:     ccb.filename,
		LocalCount:   len(ccb2.locals.table),
		Code:         code.Assemble(ccb2),
		Constants:    ccb2.constants.table,
		Names:        ccb2.names.table,
		Locals:       ccb2.locals.table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
	}

	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(props))

	if class.Parent == "" {
		compileLoadNull(ccb)
	} else {
		compile(ccb, &ast.Identifier{Value: class.Parent})
	}

	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.MakeStringObj(class.Name)))
	ccb.code.addInst(opcode.BuildClass, uint16(len(class.Methods)))
}

func compileTryCatch(ccb *codeBlockCompiler, try *ast.TryCatchExpression) {
	_, tryNoNil := try.Try.Statements[len(try.Try.Statements)-1].(*ast.ExpressionStatement)
	_, catchNoNil := try.Catch.Statements[len(try.Catch.Statements)-1].(*ast.ExpressionStatement)

	catchBlkLbl := randomLabel("catch_")
	endTryLbl := randomLabel("endTry_")

	ccb.code.addLabeledArgs(opcode.StartTry, catchBlkLbl)
	compile(ccb, try.Try)
	ccb.code.addLabeledArgs(opcode.JumpAbsolute, endTryLbl)

	ccb.code.addLabel(catchBlkLbl)
	if try.Symbol == nil {
		ccb.code.addInst(opcode.Pop)
	} else {
		ccb.code.addInst(opcode.Define, ccb.locals.indexOf(try.Symbol.Value))
	}

	compile(ccb, try.Catch)
	if try.Symbol != nil {
		ccb.code.addInst(opcode.DeleteFast, ccb.locals.indexOf(try.Symbol.Value))
	}

	if catchNoNil && !tryNoNil {
		ccb.code.addInst(opcode.JumpForward, 3)
	}

	if !tryNoNil || !catchNoNil {
		compileLoadNull(ccb)
	}

	ccb.code.addLabel(endTryLbl)
	ccb.code.addInst(opcode.EndBlock)
}

func compileBlock(ccb *codeBlockCompiler, block *ast.BlockStatement) {
	l := len(block.Statements) - 1
	for i, s := range block.Statements {
		compile(ccb, s)
		if i < l {
			if _, ok := s.(*ast.ExpressionStatement); ok {
				ccb.code.addInst(opcode.Pop)
			}
		}
	}
}

func compileFunction(ccb *codeBlockCompiler, fn *ast.FunctionLiteral, inClass, hasParent bool) {
	var body *CodeBlock
	if fn.Native {
		body = &CodeBlock{
			Name:     ccb.name + "." + fn.FQName,
			Filename: ccb.filename,
			Native:   true,
		}
	} else {
		ccb2 := &codeBlockCompiler{
			constants: newConstantTable(),
			locals:    newStringTable(),
			names:     newStringTable(),
			code:      NewInstSet(),
			filename:  ccb.filename,
			name:      ccb.name,
			inLoop:    ccb.inLoop,
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

			if !ccb2.code.last().Is(opcode.Return) {
				ccb2.code.addInst(opcode.Return)
			}
		} else {
			compileLoadNull(ccb2)
			ccb2.code.addInst(opcode.Return)
		}

		code := ccb2.code
		body = &CodeBlock{
			Name:         ccb.name + "." + fn.FQName,
			Filename:     ccb.filename,
			LocalCount:   len(ccb2.locals.table),
			Code:         code.Assemble(ccb2),
			Constants:    ccb2.constants.table,
			Names:        ccb2.names.table,
			Locals:       ccb2.locals.table,
			MaxStackSize: calculateStackSize(code),
			MaxBlockSize: calculateBlockSize(code),
		}
	}

	body.ClassMethod = inClass

	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(body))

	for _, p := range fn.Parameters {
		ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.MakeStringObj(p.Value)))
	}
	ccb.code.addInst(opcode.MakeArray, uint16(len(fn.Parameters)))

	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.MakeStringObj(fn.Name)))

	ccb.code.addInst(opcode.MakeFunction)
}

func compileIfStatement(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	if ifs.Alternative == nil {
		compileIfStatementNoElse(ccb, ifs)
		return
	}

	compile(ccb, ifs.Condition)

	_, trueNoNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	falseBrnLbl := randomLabel("false_")
	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, falseBrnLbl)
	compile(ccb, ifs.Consequence)
	if !trueNoNil {
		compileLoadNull(ccb)
	}

	_, falseNoNil := ifs.Alternative.Statements[len(ifs.Alternative.Statements)-1].(*ast.ExpressionStatement)
	afterIfStmt := randomLabel("afterIf_")
	ccb.code.addLabeledArgs(opcode.JumpAbsolute, afterIfStmt)
	ccb.code.addLabel(falseBrnLbl)
	compile(ccb, ifs.Alternative)
	ccb.code.addLabel(afterIfStmt)
	if !falseNoNil {
		compileLoadNull(ccb)
	}
}

func compileIfStatementNoElse(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	compile(ccb, ifs.Condition)

	_, noNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	falseBrnLbl := randomLabel("false_")
	afterIfStmt := randomLabel("afterIf_")

	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, falseBrnLbl)
	compile(ccb, ifs.Consequence)
	if !noNil {
		compileLoadNull(ccb)
	}

	ccb.code.addLabeledArgs(opcode.JumpAbsolute, afterIfStmt)
	ccb.code.addLabel(falseBrnLbl)
	compileLoadNull(ccb)
	ccb.code.addLabel(afterIfStmt)
}

func compileLoadNull(ccb *codeBlockCompiler) {
	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.NullConst))
}

func compileCompareExpression(ccb *codeBlockCompiler, cmp *ast.CompareExpression) {
	compile(ccb, cmp.Left)

	afterCompareLabel := randomLabel("cmp_")

	if cmp.Token.Type == token.LAnd {
		ccb.code.addLabeledArgs(opcode.JumpIfFalseOrPop, afterCompareLabel)
	} else {
		ccb.code.addLabeledArgs(opcode.JumpIfTrueOrPop, afterCompareLabel)
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

	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")

	// A loop begins with a PREPARE_BLOCK opcode this creates the first layer environment
	ccb.code.addInst(opcode.OpenScope)
	// Initialization is done in this first layer
	compile(ccb, loop.Init)

	condCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
	}

	// Compile the loop's condition check code
	compile(condCCB, loop.Condition)

	// Prepare for main body
	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
	}

	// Compile main body of loop
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	// Prepare for iteration code
	iterCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
	}

	// Compile iteration
	compile(iterCCB, loop.Iter)

	// Again, copy over the locals for indexing
	for _, n := range iterCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	ccb.code.addLabeledArgs(opcode.StartLoop, endBlockLbl, iterBlockLbl)

	ccb.code.merge(condCCB.code)
	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, endBlockLbl)
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl)
	ccb.code.merge(iterCCB.code)
	ccb.code.addInst(opcode.NextIter)
	ccb.code.addLabel(endBlockLbl)
	ccb.code.addInst(opcode.EndBlock)
	ccb.code.addInst(opcode.CloseScope)
	ccb.code.addInst(opcode.CloseScope)
}

func compileInfiniteLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")

	ccb.code.addLabeledArgs(opcode.StartLoop, endBlockLbl, iterBlockLbl)

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
	}
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl)
	ccb.code.addInst(opcode.NextIter)
	ccb.code.addLabel(endBlockLbl)
	ccb.code.addInst(opcode.EndBlock)
	ccb.code.addInst(opcode.CloseScope)
}

func compileWhileLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")

	condCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
	}

	// Compile the loop's condition check code
	compile(condCCB, loop.Condition)

	// Prepare for main body
	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
	}

	// Compile main body of loop
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	ccb.code.addLabeledArgs(opcode.StartLoop, endBlockLbl, iterBlockLbl)

	ccb.code.merge(condCCB.code)
	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, endBlockLbl)
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl)
	ccb.code.addInst(opcode.NextIter)
	ccb.code.addLabel(endBlockLbl)
	ccb.code.addInst(opcode.EndBlock)
	ccb.code.addInst(opcode.CloseScope)
}

func compileIterLoop(ccb *codeBlockCompiler, loop *ast.IterLoopStatement) {
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")
	endIterLbl := randomLabel("end_iter_")

	compile(ccb, loop.Iter)
	ccb.code.addInst(opcode.GetIter)

	ccb.code.addLabeledArgs(opcode.StartLoop, endBlockLbl, iterBlockLbl)
	ccb.code.addInst(opcode.Dup)
	ccb.code.addInst(opcode.LoadAttribute, ccb.names.indexOf("_next"))
	ccb.code.addInst(opcode.Call, 0)

	ccb.code.addInst(opcode.Dup)
	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.NullConst))
	ccb.code.addInst(opcode.Compare, uint16(opcode.CmpEq))
	ccb.code.addLabeledArgs(opcode.PopJumpIfTrue, endIterLbl)
	ccb.code.addInst(opcode.JumpForward, 4)

	ccb.code.addLabel(endIterLbl)
	ccb.code.addInst(opcode.Pop) // Duplicated return from _next()
	ccb.code.addLabeledArgs(opcode.JumpAbsolute, endBlockLbl)

	bodyStrTable := newStringTableOffset(len(ccb.locals.table))

	if loop.Key != nil {
		ccb.code.addInst(opcode.Dup)
		ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.MakeIntObj(0)))
		ccb.code.addInst(opcode.LoadIndex)
		ccb.code.addInst(opcode.Define, ccb.locals.indexOf(loop.Key.Value))
		bodyStrTable.indexOf(loop.Key.Value)
	}

	ccb.code.addInst(opcode.LoadConst, ccb.constants.indexOf(object.MakeIntObj(1)))
	ccb.code.addInst(opcode.LoadIndex)
	ccb.code.addInst(opcode.Define, ccb.locals.indexOf(loop.Value.Value))
	bodyStrTable.indexOf(loop.Value.Value)

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    bodyStrTable,
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
	}
	compile(bodyCCB, loop.Body)

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl)
	ccb.code.addInst(opcode.NextIter)
	ccb.code.addLabel(endBlockLbl)
	ccb.code.addInst(opcode.EndBlock)
	ccb.code.addInst(opcode.CloseScope)
	ccb.code.addInst(opcode.Pop) // Iterator object
}

func compileDoBlock(ccb *codeBlockCompiler, node *ast.DoExpression) {
	ccb.code.addInst(opcode.OpenScope)

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    ccb.inLoop,
	}
	compile(bodyCCB, node.Statements)

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}
	ccb.code.merge(bodyCCB.code)

	ccb.code.addInst(opcode.CloseScope)
}
