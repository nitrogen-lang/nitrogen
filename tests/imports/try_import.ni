import "std/test"

test.run("Import directory", fn(assert) {
    try {
        import '../../testdata/math.ni'
    } catch e {
        println("Test Failed: ", e)
        exit(1)
    }

    assert.isTrue(isFunc(math.add))
})
