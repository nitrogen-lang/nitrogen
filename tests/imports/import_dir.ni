import "std/test"

test.run("Import directory", fn(assert, check) {
    import '../../testdata/math2' as math
    check(assert.isTrue(isFunc(math.add)))
    check(assert.isEq(math.add(2, 4), 6))
})
