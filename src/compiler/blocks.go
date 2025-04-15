package compiler

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/token"
	"github.com/nitrogen-lang/nitrogen/src/vm/opcode"
)

func compileClassLiteral(ccb *codeBlockCompiler, class *ast.ClassLiteral) {
	ccb.linenum = class.Token.Pos.Line

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
		linenum:   ccb.linenum,
	}

	for _, f := range class.Fields {
		compile(ccb2, f)
	}
	compileLoadNull(ccb2)
	ccb2.code.addInst(opcode.Return, ccb2.linenum)

	code := ccb2.code
	assembledCode, lineOffsets := code.Assemble(ccb2)
	props := &CodeBlock{
		Name:         fmt.Sprintf("%s.__init", class.Name),
		Filename:     ccb.filename,
		LocalCount:   len(ccb2.locals.table),
		Code:         assembledCode,
		Constants:    ccb2.constants.table,
		Names:        ccb2.names.table,
		Locals:       ccb2.locals.table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
		LineOffsets:  lineOffsets,
	}

	ccb.linenum = ccb2.linenum
	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(props))

	if class.Parent == "" {
		compileLoadNull(ccb)
	} else {
		compile(ccb, &ast.Identifier{Value: class.Parent})
	}

	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.MakeStringObj(class.Name)))
	ccb.code.addInst(opcode.BuildClass, ccb.linenum, uint16(len(class.Methods)))
}

func compileBlock(ccb *codeBlockCompiler, block *ast.BlockStatement) {
	ccb.linenum = block.Token.Pos.Line
	l := len(block.Statements) - 1
	for i, s := range block.Statements {
		compile(ccb, s)
		if i < l {
			if _, ok := s.(*ast.ExpressionStatement); ok {
				ccb.code.addInst(opcode.Pop, ccb.linenum)
			}
		}
	}
}

func compileFunction(ccb *codeBlockCompiler, fn *ast.FunctionLiteral, inClass, hasParent bool) {
	ccb.linenum = fn.Token.Pos.Line
	var body *CodeBlock
	if fn.Native {
		body = &CodeBlock{
			Name:        ccb.name + "." + fn.FQName,
			Filename:    ccb.filename,
			Native:      true,
			LineOffsets: []uint16{0, uint16(ccb.linenum)},
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
			linenum:   ccb.linenum,
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
				ccb2.code.addInst(opcode.Return, ccb2.linenum)
			}
		} else {
			compileLoadNull(ccb2)
			ccb2.code.addInst(opcode.Return, ccb2.linenum)
		}

		code := ccb2.code
		assembledCode, lineOffsets := code.Assemble(ccb2)
		body = &CodeBlock{
			Name:         ccb.name + "." + fn.FQName,
			Filename:     ccb.filename,
			LocalCount:   len(ccb2.locals.table),
			Code:         assembledCode,
			Constants:    ccb2.constants.table,
			Names:        ccb2.names.table,
			Locals:       ccb2.locals.table,
			MaxStackSize: calculateStackSize(code),
			MaxBlockSize: calculateBlockSize(code),
			LineOffsets:  lineOffsets,
		}
		ccb.linenum = ccb2.linenum
	}

	body.ClassMethod = inClass

	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(body))

	for _, p := range fn.Parameters {
		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.MakeStringObj(p.Value)))
	}
	ccb.code.addInst(opcode.MakeArray, ccb.linenum, uint16(len(fn.Parameters)))

	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.MakeStringObj(fn.Name)))

	ccb.code.addInst(opcode.MakeFunction, ccb.linenum)
}

func compileIfStatement(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	ccb.linenum = ifs.Token.Pos.Line
	if ifs.Alternative == nil {
		compileIfStatementNoElse(ccb, ifs)
		return
	}

	compile(ccb, ifs.Condition)

	ccb.linenum = ifs.Consequence.Token.Pos.Line
	_, trueNoNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	falseBrnLbl := randomLabel("false_")
	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, ccb.linenum, falseBrnLbl)
	compile(ccb, ifs.Consequence)
	if !trueNoNil {
		compileLoadNull(ccb)
	}

	ccb.linenum = ifs.Alternative.Token.Pos.Line
	_, falseNoNil := ifs.Alternative.Statements[len(ifs.Alternative.Statements)-1].(*ast.ExpressionStatement)
	afterIfStmt := randomLabel("afterIf_")
	ccb.code.addLabeledArgs(opcode.JumpAbsolute, ccb.linenum, afterIfStmt)
	ccb.code.addLabel(falseBrnLbl, ccb.linenum)
	compile(ccb, ifs.Alternative)
	ccb.code.addLabel(afterIfStmt, ccb.linenum)
	if !falseNoNil {
		compileLoadNull(ccb)
	}
}

func compileIfStatementNoElse(ccb *codeBlockCompiler, ifs *ast.IfExpression) {
	compile(ccb, ifs.Condition)

	ccb.linenum = ifs.Consequence.Token.Pos.Line
	_, noNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	falseBrnLbl := randomLabel("false_")
	afterIfStmt := randomLabel("afterIf_")

	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, ccb.linenum, falseBrnLbl)
	compile(ccb, ifs.Consequence)
	if !noNil {
		compileLoadNull(ccb)
	}

	ccb.code.addLabeledArgs(opcode.JumpAbsolute, ccb.linenum, afterIfStmt)
	ccb.code.addLabel(falseBrnLbl, ccb.linenum)
	compileLoadNull(ccb)
	ccb.code.addLabel(afterIfStmt, ccb.linenum)
}

func compileLoadNull(ccb *codeBlockCompiler) {
	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.NullConst))
}

func compileCompareExpression(ccb *codeBlockCompiler, cmp *ast.CompareExpression) {
	ccb.linenum = cmp.Token.Pos.Line
	compile(ccb, cmp.Left)

	afterCompareLabel := randomLabel("cmp_")

	if cmp.Token.Type == token.LAnd {
		ccb.code.addLabeledArgs(opcode.JumpIfFalseOrPop, ccb.linenum, afterCompareLabel)
	} else {
		ccb.code.addLabeledArgs(opcode.JumpIfTrueOrPop, ccb.linenum, afterCompareLabel)
	}

	compile(ccb, cmp.Right)
	ccb.code.addLabel(afterCompareLabel, ccb.linenum)
}

func compileLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	ccb.linenum = loop.Token.Pos.Line
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
	ccb.code.addInst(opcode.StartBlock, ccb.linenum)
	// Initialization is done in this first layer
	compile(ccb, loop.Init)

	condCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		linenum:   ccb.linenum,
	}

	// Compile the loop's condition check code
	compile(condCCB, loop.Condition)
	ccb.linenum = condCCB.linenum

	// Prepare for main body
	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
		linenum:   ccb.linenum,
	}

	// Compile main body of loop
	compile(bodyCCB, loop.Body)
	ccb.linenum = bodyCCB.linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop, ccb.linenum)
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
		linenum:   ccb.linenum,
	}

	// Compile iteration
	compile(iterCCB, loop.Iter)
	ccb.linenum = iterCCB.linenum

	// Again, copy over the locals for indexing
	for _, n := range iterCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	ccb.code.addLabeledArgs(opcode.StartLoop, ccb.linenum, endBlockLbl, iterBlockLbl)

	ccb.code.merge(condCCB.code)
	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, ccb.linenum, endBlockLbl)
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl, ccb.linenum)
	ccb.code.merge(iterCCB.code)
	ccb.code.addInst(opcode.NextIter, ccb.linenum)
	ccb.code.addLabel(endBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.EndBlock, ccb.linenum)
	ccb.code.addInst(opcode.EndBlock, ccb.linenum)
}

func compileInfiniteLoop(ccb *codeBlockCompiler, loop *ast.LoopStatement) {
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")

	ccb.code.addLabeledArgs(opcode.StartLoop, ccb.linenum, endBlockLbl, iterBlockLbl)

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
		linenum:   ccb.linenum,
	}
	compile(bodyCCB, loop.Body)
	ccb.linenum = bodyCCB.linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop, ccb.linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.NextIter, ccb.linenum)
	ccb.code.addLabel(endBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.EndBlock, ccb.linenum)
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
		linenum:   ccb.linenum,
	}

	// Compile the loop's condition check code
	compile(condCCB, loop.Condition)
	ccb.linenum = condCCB.linenum

	// Prepare for main body
	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
		linenum:   ccb.linenum,
	}

	// Compile main body of loop
	compile(bodyCCB, loop.Body)
	ccb.linenum = bodyCCB.linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop, ccb.linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}

	ccb.code.addLabeledArgs(opcode.StartLoop, ccb.linenum, endBlockLbl, iterBlockLbl)

	ccb.code.merge(condCCB.code)
	ccb.code.addLabeledArgs(opcode.PopJumpIfFalse, ccb.linenum, endBlockLbl)
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.NextIter, ccb.linenum)
	ccb.code.addLabel(endBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.EndBlock, ccb.linenum)
}

func compileIterLoop(ccb *codeBlockCompiler, loop *ast.IterLoopStatement) {
	ccb.linenum = loop.Token.Pos.Line
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")
	endIterLbl := randomLabel("end_iter_")

	compile(ccb, loop.Iter)
	ccb.code.addInst(opcode.GetIter, ccb.linenum)

	ccb.code.addLabeledArgs(opcode.StartLoop, ccb.linenum, endBlockLbl, iterBlockLbl)
	ccb.code.addInst(opcode.Dup, ccb.linenum)
	ccb.code.addInst(opcode.LoadAttribute, ccb.linenum, ccb.names.indexOf("_next"))
	ccb.code.addInst(opcode.Call, ccb.linenum, 0)

	ccb.code.addInst(opcode.Dup, ccb.linenum)
	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.NullConst))
	ccb.code.addInst(opcode.Compare, ccb.linenum, uint16(opcode.CmpEq))
	ccb.code.addLabeledArgs(opcode.PopJumpIfTrue, ccb.linenum, endIterLbl)
	ccb.code.addInst(opcode.JumpForward, ccb.linenum, 4)

	ccb.code.addLabel(endIterLbl, ccb.linenum)
	ccb.code.addInst(opcode.Pop, ccb.linenum) // Duplicated return from _next()
	ccb.code.addLabeledArgs(opcode.JumpAbsolute, ccb.linenum, endBlockLbl)

	bodyStrTable := newStringTableOffset(len(ccb.locals.table))

	if loop.Key != nil {
		ccb.code.addInst(opcode.Dup, ccb.linenum)
		ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.MakeIntObj(0)))
		ccb.code.addInst(opcode.LoadIndex, ccb.linenum)
		ccb.code.addInst(opcode.Define, ccb.linenum, ccb.locals.indexOf(loop.Key.Value), 0)
		bodyStrTable.indexOf(loop.Key.Value)
	}

	ccb.code.addInst(opcode.LoadConst, ccb.linenum, ccb.constants.indexOf(object.MakeIntObj(1)))
	ccb.code.addInst(opcode.LoadIndex, ccb.linenum)
	ccb.code.addInst(opcode.Define, ccb.linenum, ccb.locals.indexOf(loop.Value.Value), 0)
	bodyStrTable.indexOf(loop.Value.Value)

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    bodyStrTable,
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    true,
		linenum:   ccb.linenum,
	}
	compile(bodyCCB, loop.Body)
	ccb.linenum = bodyCCB.linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.code.addInst(opcode.Pop, ccb.linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(iterBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.NextIter, ccb.linenum)
	ccb.code.addLabel(endBlockLbl, ccb.linenum)
	ccb.code.addInst(opcode.EndBlock, ccb.linenum)
	ccb.code.addInst(opcode.Pop, ccb.linenum) // Iterator object
}

func compileDoBlock(ccb *codeBlockCompiler, node *ast.DoExpression) {
	endBlockLabel := randomLabel("endBlk_")

	ccb.linenum = node.Token.Pos.Line
	if node.Recoverable {
		ccb.code.addLabeledArgs(opcode.Recover, ccb.linenum, endBlockLabel)
	} else {
		ccb.code.addInst(opcode.StartBlock, ccb.linenum)
	}

	bodyCCB := &codeBlockCompiler{
		constants: ccb.constants,
		locals:    newStringTableOffset(len(ccb.locals.table)),
		names:     ccb.names,
		code:      NewInstSet(),
		filename:  ccb.filename,
		name:      ccb.name,
		inLoop:    ccb.inLoop,
		linenum:   ccb.linenum,
	}
	compile(bodyCCB, node.Statements)
	ccb.linenum = bodyCCB.linenum

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.locals.table[len(ccb.locals.table):] {
		ccb.locals.indexOf(n)
	}
	ccb.code.merge(bodyCCB.code)

	ccb.code.addLabel(endBlockLabel, ccb.linenum)
	ccb.code.addInst(opcode.EndBlock, ccb.linenum)
}
