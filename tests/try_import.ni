try {
    import './includes/math.ni'
} catch e {
    println("Test Failed: ", e)
}

if !isFunc(math.add) {
    println("math.add is not a function")
}
