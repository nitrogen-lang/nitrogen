package compiler

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/elemental/compile"
	"github.com/nitrogen-lang/nitrogen/src/elemental/object"
	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func compileClassLiteral(ccb *compile.CodeBlockCompiler, class *ast.ClassLiteral) {
	ccb.Linenum = class.Token.Pos.Line

	for _, f := range class.Methods {
		f.FQName = fmt.Sprintf("%s.%s", class.Name, f.Name)
		compileFunction(ccb, f, true, class.Parent != "")
	}

	ccb2 := &compile.CodeBlockCompiler{
		Constants: compile.NewConstantTable(),
		Locals:    compile.NewStringTable(),
		Names:     compile.NewStringTable(),
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		InLoop:    ccb.InLoop,
		Linenum:   ccb.Linenum,
	}

	for _, f := range class.Fields {
		compileMain(ccb2, f)
	}
	compileLoadNull(ccb2)
	ccb2.Code.AddInst(opcode.Return, ccb2.Linenum)

	code := ccb2.Code
	assembledCode, lineOffsets := code.Assemble(ccb2)
	props := &compile.CodeBlock{
		Name:         fmt.Sprintf("%s.__init", class.Name),
		Filename:     ccb.Filename,
		LocalCount:   len(ccb2.Locals.Table),
		Code:         assembledCode,
		Constants:    ccb2.Constants.Table,
		Names:        ccb2.Names.Table,
		Locals:       ccb2.Locals.Table,
		MaxStackSize: calculateStackSize(code),
		MaxBlockSize: calculateBlockSize(code),
		LineOffsets:  lineOffsets,
	}

	ccb.Linenum = ccb2.Linenum
	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(props))

	if class.Parent == "" {
		compileLoadNull(ccb)
	} else {
		compileMain(ccb, &ast.Identifier{Value: class.Parent})
	}

	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.MakeStringObj(class.Name)))
	ccb.Code.AddInst(opcode.BuildClass, ccb.Linenum, uint16(len(class.Methods)))
}

func compileBlock(ccb *compile.CodeBlockCompiler, block *ast.BlockStatement) {
	ccb.Linenum = block.Token.Pos.Line
	l := len(block.Statements) - 1
	for i, s := range block.Statements {
		compileMain(ccb, s)
		if i < l {
			if _, ok := s.(*ast.ExpressionStatement); ok {
				ccb.Code.AddInst(opcode.Pop, ccb.Linenum)
			}
		}
	}
}

func compileFunction(ccb *compile.CodeBlockCompiler, fn *ast.FunctionLiteral, inClass, hasParent bool) {
	ccb.Linenum = fn.Token.Pos.Line
	var body *compile.CodeBlock
	if fn.Native {
		body = &compile.CodeBlock{
			Name:        ccb.Name + "." + fn.FQName,
			Filename:    ccb.Filename,
			Native:      true,
			LineOffsets: []uint16{0, uint16(ccb.Linenum)},
		}
	} else {
		ccb2 := &compile.CodeBlockCompiler{
			Constants: compile.NewConstantTable(),
			Locals:    compile.NewStringTable(),
			Names:     compile.NewStringTable(),
			Code:      compile.NewInstSet(),
			Filename:  ccb.Filename,
			Name:      ccb.Name,
			InLoop:    ccb.InLoop,
			Linenum:   ccb.Linenum,
		}

		for _, p := range fn.Parameters {
			ccb2.Locals.IndexOf(p.Value)
		}
		ccb2.Locals.IndexOf("arguments") // `arguments` holds any remaining arguments from a function call
		if inClass {
			ccb2.Locals.IndexOf("this")
			if hasParent {
				ccb2.Locals.IndexOf("parent")
			}
		}

		compileMain(ccb2, fn.Body)

		if len(fn.Body.Statements) > 0 {
			switch fn.Body.Statements[len(fn.Body.Statements)-1].(type) {
			case *ast.ExpressionStatement:
				break
			case *ast.ReturnStatement:
				break
			default:
				compileLoadNull(ccb2)
			}

			if !ccb2.Code.Last().Is(opcode.Return) {
				ccb2.Code.AddInst(opcode.Return, ccb2.Linenum)
			}
		} else {
			compileLoadNull(ccb2)
			ccb2.Code.AddInst(opcode.Return, ccb2.Linenum)
		}

		code := ccb2.Code
		assembledCode, lineOffsets := code.Assemble(ccb2)
		body = &compile.CodeBlock{
			Name:         ccb.Name + "." + fn.FQName,
			Filename:     ccb.Filename,
			LocalCount:   len(ccb2.Locals.Table),
			Code:         assembledCode,
			Constants:    ccb2.Constants.Table,
			Names:        ccb2.Names.Table,
			Locals:       ccb2.Locals.Table,
			MaxStackSize: calculateStackSize(code),
			MaxBlockSize: calculateBlockSize(code),
			LineOffsets:  lineOffsets,
		}
		ccb.Linenum = ccb2.Linenum
	}

	body.ClassMethod = inClass

	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(body))

	for _, p := range fn.Parameters {
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.MakeStringObj(p.Value)))
	}
	ccb.Code.AddInst(opcode.MakeArray, ccb.Linenum, uint16(len(fn.Parameters)))

	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.MakeStringObj(fn.Name)))

	ccb.Code.AddInst(opcode.MakeFunction, ccb.Linenum)
}

func compileIfStatement(ccb *compile.CodeBlockCompiler, ifs *ast.IfExpression) {
	ccb.Linenum = ifs.Token.Pos.Line
	if ifs.Alternative == nil {
		compileIfStatementNoElse(ccb, ifs)
		return
	}

	compileMain(ccb, ifs.Condition)

	ccb.Linenum = ifs.Consequence.Token.Pos.Line
	_, trueNoNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	falseBrnLbl := randomLabel("false_")
	ccb.Code.AddLabeledArgs(opcode.PopJumpIfFalse, ccb.Linenum, falseBrnLbl)
	compileMain(ccb, ifs.Consequence)
	if !trueNoNil {
		compileLoadNull(ccb)
	}

	ccb.Linenum = ifs.Alternative.Token.Pos.Line
	_, falseNoNil := ifs.Alternative.Statements[len(ifs.Alternative.Statements)-1].(*ast.ExpressionStatement)
	afterIfStmt := randomLabel("afterIf_")
	ccb.Code.AddLabeledArgs(opcode.JumpAbsolute, ccb.Linenum, afterIfStmt)
	ccb.Code.AddLabel(falseBrnLbl, ccb.Linenum)
	compileMain(ccb, ifs.Alternative)
	ccb.Code.AddLabel(afterIfStmt, ccb.Linenum)
	if !falseNoNil {
		compileLoadNull(ccb)
	}
}

func compileIfStatementNoElse(ccb *compile.CodeBlockCompiler, ifs *ast.IfExpression) {
	compileMain(ccb, ifs.Condition)

	ccb.Linenum = ifs.Consequence.Token.Pos.Line
	_, noNil := ifs.Consequence.Statements[len(ifs.Consequence.Statements)-1].(*ast.ExpressionStatement)
	falseBrnLbl := randomLabel("false_")
	afterIfStmt := randomLabel("afterIf_")

	ccb.Code.AddLabeledArgs(opcode.PopJumpIfFalse, ccb.Linenum, falseBrnLbl)
	compileMain(ccb, ifs.Consequence)
	if !noNil {
		compileLoadNull(ccb)
	}

	ccb.Code.AddLabeledArgs(opcode.JumpAbsolute, ccb.Linenum, afterIfStmt)
	ccb.Code.AddLabel(falseBrnLbl, ccb.Linenum)
	compileLoadNull(ccb)
	ccb.Code.AddLabel(afterIfStmt, ccb.Linenum)
}

func compileLoadNull(ccb *compile.CodeBlockCompiler) {
	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.NullConst))
}

func compileCompareExpression(ccb *compile.CodeBlockCompiler, cmp *ast.CompareExpression) {
	ccb.Linenum = cmp.Token.Pos.Line
	compileMain(ccb, cmp.Left)

	afterCompareLabel := randomLabel("cmp_")

	if cmp.Token.Type == token.LAnd {
		ccb.Code.AddLabeledArgs(opcode.JumpIfFalseOrPop, ccb.Linenum, afterCompareLabel)
	} else {
		ccb.Code.AddLabeledArgs(opcode.JumpIfTrueOrPop, ccb.Linenum, afterCompareLabel)
	}

	compileMain(ccb, cmp.Right)
	ccb.Code.AddLabel(afterCompareLabel, ccb.Linenum)
}

func compileLoop(ccb *compile.CodeBlockCompiler, loop *ast.LoopStatement) {
	ccb.Linenum = loop.Token.Pos.Line
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
	ccb.Code.AddInst(opcode.StartBlock, ccb.Linenum)
	// Initialization is done in this first layer
	compileMain(ccb, loop.Init)

	condCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		Linenum:   ccb.Linenum,
	}

	// Compile the loop's condition check code
	compileMain(condCCB, loop.Condition)
	ccb.Linenum = condCCB.Linenum

	// Prepare for main body
	bodyCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		InLoop:    true,
		Linenum:   ccb.Linenum,
	}

	// Compile main body of loop
	compileMain(bodyCCB, loop.Body)
	ccb.Linenum = bodyCCB.Linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.Code.AddInst(opcode.Pop, ccb.Linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.Locals.Table[len(ccb.Locals.Table):] {
		ccb.Locals.IndexOf(n)
	}

	// Prepare for iteration code
	iterCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		Linenum:   ccb.Linenum,
	}

	// Compile iteration
	compileMain(iterCCB, loop.Iter)
	ccb.Linenum = iterCCB.Linenum

	// Again, copy over the locals for indexing
	for _, n := range iterCCB.Locals.Table[len(ccb.Locals.Table):] {
		ccb.Locals.IndexOf(n)
	}

	ccb.Code.AddLabeledArgs(opcode.StartLoop, ccb.Linenum, endBlockLbl, iterBlockLbl)

	ccb.Code.Merge(condCCB.Code)
	ccb.Code.AddLabeledArgs(opcode.PopJumpIfFalse, ccb.Linenum, endBlockLbl)
	ccb.Code.Merge(bodyCCB.Code)

	ccb.Code.AddLabel(iterBlockLbl, ccb.Linenum)
	ccb.Code.Merge(iterCCB.Code)
	ccb.Code.AddInst(opcode.NextIter, ccb.Linenum)
	ccb.Code.AddLabel(endBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.EndBlock, ccb.Linenum)
	ccb.Code.AddInst(opcode.EndBlock, ccb.Linenum)
}

func compileInfiniteLoop(ccb *compile.CodeBlockCompiler, loop *ast.LoopStatement) {
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")

	ccb.Code.AddLabeledArgs(opcode.StartLoop, ccb.Linenum, endBlockLbl, iterBlockLbl)

	bodyCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		InLoop:    true,
		Linenum:   ccb.Linenum,
	}
	compileMain(bodyCCB, loop.Body)
	ccb.Linenum = bodyCCB.Linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.Code.AddInst(opcode.Pop, ccb.Linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.Locals.Table[len(ccb.Locals.Table):] {
		ccb.Locals.IndexOf(n)
	}
	ccb.Code.Merge(bodyCCB.Code)

	ccb.Code.AddLabel(iterBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.NextIter, ccb.Linenum)
	ccb.Code.AddLabel(endBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.EndBlock, ccb.Linenum)
}

func compileWhileLoop(ccb *compile.CodeBlockCompiler, loop *ast.LoopStatement) {
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")

	condCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		Linenum:   ccb.Linenum,
	}

	// Compile the loop's condition check code
	compileMain(condCCB, loop.Condition)
	ccb.Linenum = condCCB.Linenum

	// Prepare for main body
	bodyCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		InLoop:    true,
		Linenum:   ccb.Linenum,
	}

	// Compile main body of loop
	compileMain(bodyCCB, loop.Body)
	ccb.Linenum = bodyCCB.Linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.Code.AddInst(opcode.Pop, ccb.Linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.Locals.Table[len(ccb.Locals.Table):] {
		ccb.Locals.IndexOf(n)
	}

	ccb.Code.AddLabeledArgs(opcode.StartLoop, ccb.Linenum, endBlockLbl, iterBlockLbl)

	ccb.Code.Merge(condCCB.Code)
	ccb.Code.AddLabeledArgs(opcode.PopJumpIfFalse, ccb.Linenum, endBlockLbl)
	ccb.Code.Merge(bodyCCB.Code)

	ccb.Code.AddLabel(iterBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.NextIter, ccb.Linenum)
	ccb.Code.AddLabel(endBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.EndBlock, ccb.Linenum)
}

func compileIterLoop(ccb *compile.CodeBlockCompiler, loop *ast.IterLoopStatement) {
	ccb.Linenum = loop.Token.Pos.Line
	endBlockLbl := randomLabel("end_")
	iterBlockLbl := randomLabel("iter_")
	endIterLbl := randomLabel("end_iter_")

	compileMain(ccb, loop.Iter)
	ccb.Code.AddInst(opcode.GetIter, ccb.Linenum)

	ccb.Code.AddLabeledArgs(opcode.StartLoop, ccb.Linenum, endBlockLbl, iterBlockLbl)
	ccb.Code.AddInst(opcode.Dup, ccb.Linenum)
	ccb.Code.AddInst(opcode.LoadAttribute, ccb.Linenum, ccb.Names.IndexOf("_next"))
	ccb.Code.AddInst(opcode.Call, ccb.Linenum, 0)

	ccb.Code.AddInst(opcode.Dup, ccb.Linenum)
	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.NullConst))
	ccb.Code.AddInst(opcode.Compare, ccb.Linenum, uint16(opcode.CmpEq))
	ccb.Code.AddLabeledArgs(opcode.PopJumpIfTrue, ccb.Linenum, endIterLbl)
	ccb.Code.AddInst(opcode.JumpForward, ccb.Linenum, 4)

	ccb.Code.AddLabel(endIterLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.Pop, ccb.Linenum) // Duplicated return from _next()
	ccb.Code.AddLabeledArgs(opcode.JumpAbsolute, ccb.Linenum, endBlockLbl)

	bodyStrTable := compile.NewStringTableOffset(len(ccb.Locals.Table))

	if loop.Key != nil {
		ccb.Code.AddInst(opcode.Dup, ccb.Linenum)
		ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.MakeIntObj(0)))
		ccb.Code.AddInst(opcode.LoadIndex, ccb.Linenum)
		ccb.Code.AddInst(opcode.Define, ccb.Linenum, ccb.Locals.IndexOf(loop.Key.Value), 0)
		bodyStrTable.IndexOf(loop.Key.Value)
	}

	ccb.Code.AddInst(opcode.LoadConst, ccb.Linenum, ccb.Constants.IndexOf(object.MakeIntObj(1)))
	ccb.Code.AddInst(opcode.LoadIndex, ccb.Linenum)
	ccb.Code.AddInst(opcode.Define, ccb.Linenum, ccb.Locals.IndexOf(loop.Value.Value), 0)
	bodyStrTable.IndexOf(loop.Value.Value)

	bodyCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    bodyStrTable,
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		InLoop:    true,
		Linenum:   ccb.Linenum,
	}
	compileMain(bodyCCB, loop.Body)
	ccb.Linenum = bodyCCB.Linenum

	// If the body ends in an expression, we need to pop it so the stack is correct
	if _, ok := loop.Body.Statements[len(loop.Body.Statements)-1].(*ast.ExpressionStatement); ok {
		bodyCCB.Code.AddInst(opcode.Pop, ccb.Linenum)
	}

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.Locals.Table[len(ccb.Locals.Table):] {
		ccb.Locals.IndexOf(n)
	}
	ccb.Code.Merge(bodyCCB.Code)

	ccb.Code.AddLabel(iterBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.NextIter, ccb.Linenum)
	ccb.Code.AddLabel(endBlockLbl, ccb.Linenum)
	ccb.Code.AddInst(opcode.EndBlock, ccb.Linenum)
	ccb.Code.AddInst(opcode.Pop, ccb.Linenum) // Iterator object
}

func compileDoBlock(ccb *compile.CodeBlockCompiler, node *ast.DoExpression) {
	endBlockLabel := randomLabel("endBlk_")

	ccb.Linenum = node.Token.Pos.Line
	if node.Recoverable {
		ccb.Code.AddLabeledArgs(opcode.Recover, ccb.Linenum, endBlockLabel)
	} else {
		ccb.Code.AddInst(opcode.StartBlock, ccb.Linenum)
	}

	bodyCCB := &compile.CodeBlockCompiler{
		Constants: ccb.Constants,
		Locals:    compile.NewStringTableOffset(len(ccb.Locals.Table)),
		Names:     ccb.Names,
		Code:      compile.NewInstSet(),
		Filename:  ccb.Filename,
		Name:      ccb.Name,
		InLoop:    ccb.InLoop,
		Linenum:   ccb.Linenum,
	}
	compileMain(bodyCCB, node.Statements)
	ccb.Linenum = bodyCCB.Linenum

	// This copies the local variables into the outer compile block for table indexing
	for _, n := range bodyCCB.Locals.Table[len(ccb.Locals.Table):] {
		ccb.Locals.IndexOf(n)
	}
	ccb.Code.Merge(bodyCCB.Code)

	ccb.Code.AddLabel(endBlockLabel, ccb.Linenum)
	ccb.Code.AddInst(opcode.EndBlock, ccb.Linenum)
}
