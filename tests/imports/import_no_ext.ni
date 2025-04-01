import "std/test"

test.run("Import with no extension", fn(assert, check) {
    import '../../testdata/math' as math
    check(assert.isTrue(isFunc(math.add)))
    check(assert.isEq(math.add(2, 4), 6))
})
