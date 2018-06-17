const thing = 42

try {
    thing = 43
    println("Test Failed: Constant reassigned")
} catch {pass}
