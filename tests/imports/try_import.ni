import "std/test"

test.run("Import directory", func(assert) {
    try {
        import '../../testdata/math.ni'
    } catch e {
        println("Test Failed: ", e)
        exit(1)
    }

    assert.isTrue(isFunc(math.add))
})
