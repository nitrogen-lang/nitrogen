try {
    import '../../testdata/math.ni'
} catch e {
    println("Test Failed: ", e)
    exit(1)
}

if !isFunc(math.add) {
    println("math.add is not a function")
    exit(1)
}
