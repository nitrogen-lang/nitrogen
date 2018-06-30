if !modulesSupported() {
    return
}

import 'strings.so' as str

let testStr = "abcefguvwxyz"

if !str.hasPrefix(testStr, "abc") {
    println("Test Failed: Test string expected to have prefix but didn't")
    exit(1)
}
if !str.hasSuffix(testStr, "xyz") {
    println("Test Failed: Test string expected to have suffix but didn't")
    exit(1)
}

testStr = "name:    John"

let dedupStr = str.dedup(testStr, " ")
if dedupStr != "name: John" {
    println("Test Failed: Expected 'name: John', got ", dedupStr)
    exit(1)
}

testStr = "       test        "

let trimStr = str.trimSpace(testStr)
if trimStr != "test" {
    println("Test Failed: Expected 'test', got '", trimStr, "'")
    exit(1)
}

testStr = "one,two,three"

let splits = str.split(testStr, ",")
if len(splits) != 3 {
    println("Test Failed: Expected len of 3, got ", len(splits))
    exit(1)
}

splits = str.splitN(testStr, ",", 2)
if len(splits) != 2 {
    println("Test Failed: Expected len of 2, got ", len(splits))
    exit(1)
}

if !str.contains(testStr, "two") {
    println("Test Failed: String expected to contain 'two'")
    exit(1)
}

let commaCount = str.count(testStr, ",")
if commaCount != 2 {
    println("Test Failed: commaCount expected to be 2, got ", commaCount)
    exit(1)
}

try {
    str.count(testStr, "")
    println("Test Failed: count didn't throw")
    exit(1)
} catch {pass}

testStr = "oink oink oink"

let replaced = str.replace(testStr, "k", "ky", 2)
if replaced != "oinky oinky oink" {
    println("Test Failed: repalced not right, got ", replaced)
    exit(1)
}

replaced = str.replace(testStr, "oink", "moo", -1)
if replaced != "moo moo moo" {
    println("Test Failed: repalced not right, got ", replaced)
    exit(1)
}

try {
    str.replace(testStr, "", "moo", -1)
    println("Test Failed: replace didn't throw")
    exit(1)
} catch {pass}
