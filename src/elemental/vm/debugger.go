package vm

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/elemental/vm/opcode"
)

var debugPrintInstructionsAlways = false

func debugPrompt(vm *VirtualMachine) {
	if debugPrintInstructionsAlways {
		printDebugFrameInfo(vm)
		return
	}

	done := false
	for !done {
		cmd := debugPrompInput()
		done = debugExecCmd(cmd, vm)
	}
}

func debugPrompInput() string {
	in := os.Stdin
	out := os.Stdout
	scanner := bufio.NewScanner(in)

	fmt.Fprint(out, "> ")
	scanned := scanner.Scan()
	if !scanned {
		return ""
	}

	return strings.TrimSpace(scanner.Text())
}

func printDebugFrameInfo(vm *VirtualMachine) {
	fmt.Fprintf(vm.GetStdout(), "================\n")
	fmt.Fprintf(vm.GetStdout(), "** Next Step:\n")
	fmt.Fprintf(vm.GetStdout(), "** PC = %d; OPCODE = %s\n", vm.currentFrame.pc-1, opcode.Names[vm.currentOpcode()])
	fmt.Fprintf(vm.GetStdout(), "** FRAME_MODULE = %s\n", vm.currentFrame.module)
	fmt.Fprintf(vm.GetStdout(), "** FRAME_FILENAME = %s:%d\n", vm.currentFrame.code.Filename, vm.currentFrame.lineno())
	fmt.Fprintf(vm.GetStdout(), "================\n")
}

func debugExecCmd(input string, vm *VirtualMachine) bool {
	switch input {
	case "quit":
		vm.Exit(0)
		return true
	case "continue":
		vm.breakpoint = false
		return true
	case "step":
		return true
	case "frame":
		printDebugFrameInfo(vm)
	case "contwithinstrs":
		debugPrintInstructionsAlways = true
		return true
	case "frames":
		vm.callStack.forEach(func(f *Frame) {
			fmt.Fprintf(vm.GetStdout(), "** %s:%d in module %s\n", f.code.Filename, f.lineno(), f.module)
		})
	case "env":
		vm.currentFrame.env.Print("  ")
	}
	return false
}
