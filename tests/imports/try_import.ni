import "std/test"

test.run("Import directory", fn(assert, check) {
    const math = recover {
        import '../../testdata/math.ni'
        math
    }

    if !isMap(math) {
        println("Import test failed: ", varType(math))
        exit(1)
    }

    check(assert.isTrue(isFunc(math.add)))
})
