let math

try {
    import './includes/math.ni' as math2
    math = math2 // Required to hoist the module outside the try/catch scope
} catch e {
    println("Test Failed: ", e)
}

if !isFunc(math.add) {
    println("math.add is not a function")
}
