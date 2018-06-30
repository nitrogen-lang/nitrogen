let arr = ["one", "two"]

let prependArr = prepend(arr, "zero")
if prependArr[0] != "zero" {
    println("Test Failed: Expected \"zero\", got ", prependArr[0])
    exit(1)
}

let pushArr = push(arr, "three")
if len(pushArr) != 3 {
    println("Test Failed: Expected 3, got ", len(pushArr))
    exit(1)
}
if pushArr[2] != "three" {
    println("Test Failed: Expected \"three\", got ", pushArr[2])
    exit(1)
}

let popArr = pop(arr)
if len(popArr) != 1 {
    println("Test Failed: Expected 1, got ", len(popArr))
    exit(1)
}
if popArr[0] != "one" {
    println("Test Failed: Expected \"one\", got ", popArr[0])
    exit(1)
}

let arr2 = arr + ["three", "four"]
let spliceArr = splice(arr2, 1, 2)
if len(spliceArr) != 2 {
    println("Test Failed: Expected 2, got ", len(spliceArr))
    exit(1)
}
if spliceArr[0] != "one" {
    println("Test Failed: Expected \"one\", got ", spliceArr[0])
    exit(1)
}
if spliceArr[1] != "four" {
    println("Test Failed: Expected \"four\", got ", spliceArr[1])
    exit(1)
}

let spliceArr2 = splice(arr2, 2)
if len(spliceArr2) != 2 {
    println("Test Failed: Expected 2, got ", len(spliceArr2))
    exit(1)
}
if spliceArr2[0] != "one" {
    println("Test Failed: Expected \"one\", got ", spliceArr2[0])
    exit(1)
}
if spliceArr2[1] != "two" {
    println("Test Failed: Expected \"two\", got ", spliceArr2[1])
    exit(1)
}

let noopSplice = splice(arr2, 0)
if len(noopSplice) != 0 {
    println("Test Failed: Expected 0, got ", len(noopSplice))
    exit(1)
}

try {
    noopSplice = splice(arr2, -1)
    println("Test Failed: Negative args didn't throw")
    exit(1)
} catch {pass}

noopSplice = splice(arr2, 1, 0)
if len(noopSplice) != 4 {
    println("Test Failed: Expected 4, got ", len(noopSplice))
    exit(1)
}

try {
    noopSplice = splice(arr2, 1, -1)
    println("Test Failed: Negative args didn't throw")
    exit(1)
} catch {pass}

let sliceArr = slice(arr2, 0)
if len(sliceArr) != 4 {
    println("Test Failed: Expected 4, got ", len(sliceArr))
    exit(1)
}

sliceArr = slice(arr2, 1, 2)
if len(sliceArr) != 2 {
    println("Test Failed: Expected 2, got ", len(sliceArr))
    exit(1)
}
if sliceArr[0] != "two" {
    println("Test Failed: Expected \"one\", got ", sliceArr[0])
    exit(1)
}
if sliceArr[1] != "three" {
    println("Test Failed: Expected \"four\", got ", sliceArr[1])
    exit(1)
}
