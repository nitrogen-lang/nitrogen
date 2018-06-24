const thing = 42

try {
    thing = 43
    println("Test Failed: Constant reassigned")
} catch {pass}

try {
    let me_out = "please"
} catch { pass }

if !isDefined("me_out") {
    println("Test Failed: var in try block no available outside")
}
