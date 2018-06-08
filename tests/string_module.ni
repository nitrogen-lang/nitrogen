if !modulesSupported() {
    return
}

always str = import('../built-modules/strings.so')

let testStr = "abcefguvwxyz"

if !str.hasPrefix(testStr, "abc") {
    println("Test Failed: Test string expected to have prefix but didn't")
}
if !str.hasSuffix(testStr, "xyz") {
    println("Test Failed: Test string expected to have suffix but didn't")
}

testStr = "name:    John"

let dedupStr = str.dedup(testStr, " ")
if dedupStr != "name: John" {
    println("Test Failed: Expected 'name: John', got ", dedupStr)
}

testStr = "       test        "

let trimStr = str.trimSpace(testStr)
if trimStr != "test" {
    println("Test Failed: Expected 'test', got '", trimStr, "'")
}

testStr = "one,two,three"

let splits = str.split(testStr, ",")
if len(splits) != 3 {
    println("Test Failed: Expected len of 3, got ", len(splits))
}

splits = str.splitN(testStr, ",", 2)
if len(splits) != 2 {
    println("Test Failed: Expected len of 2, got ", len(splits))
}

if !str.contains(testStr, "two") {
    println("Test Failed: String expected to contain 'two'")
}

let commaCount = str.count(testStr, ",")
if commaCount != 2 {
    println("Test Failed: commaCount expected to be 2, got ", commaCount)
}

try {
    str.count(testStr, "")
    println("Test Failed: count didn't throw")
} catch {pass}

testStr = "oink oink oink"

let replaced = str.replace(testStr, "k", "ky", 2)
if replaced != "oinky oinky oink" {
    println("Test Failed: repalced not right, got ", replaced)
}

replaced = str.replace(testStr, "oink", "moo", -1)
if replaced != "moo moo moo" {
    println("Test Failed: repalced not right, got ", replaced)
}

try {
    str.replace(testStr, "", "moo", -1)
    println("Test Failed: replace didn't throw")
} catch {pass}
