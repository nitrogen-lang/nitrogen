package vm

import (
	"testing"
)

func TestBlockStack(t *testing.T) {
	f := &Frame{
		blockStack: make([]block, 3),
		bp:         3,
	}

	loopBlock := &forLoopBlock{}
	f.blockStack[0] = loopBlock
	f.blockStack[1] = &recoverBlock{}
	f.blockStack[2] = &recoverBlock{}

	front := f.popBlockUntil(loopBlockT)
	if front != loopBlock {
		t.Fatal("Didn't receive loop block")
	}

	if f.bp != 1 {
		t.Fatalf("Block pointer isn't right. Got %d, wanted %d", f.bp, 1)
	}

	// Reset
	f.bp = 3
	loopBlock = &forLoopBlock{}
	f.blockStack[0] = &forLoopBlock{}
	f.blockStack[1] = &recoverBlock{}
	f.blockStack[2] = loopBlock

	front = f.popBlockUntil(loopBlockT)
	if front != loopBlock {
		t.Fatal("Didn't receive loop block")
	}

	if f.bp != 3 {
		t.Fatalf("Block pointer isn't right. Got %d, wanted %d", f.bp, 1)
	}
}
