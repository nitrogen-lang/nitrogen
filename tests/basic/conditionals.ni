// This test mainly checks the parser to ensure single line if statements are parsed correctly

let counter = 0

if true {
    counter += 1
}

if (true) {
    counter += 1
}

if true: counter += 1
if (true): counter += 1

if counter != 4 {
    println("Test Failed: counter should be 4, got ", counter)
}
