import "std/test"

test.run("Non-existant import", fn(assert, check) {
    check(assert.shouldRecover(fn() {
        import './includes/_not_exist.ni'
    }))
})

test.run("Absolute import", fn(assert, check) {
    import '../../testdata/math.ni'
    check(assert.isTrue(isFunc(math.add)))
    check(assert.isEq(math.add(2, 4), 6))
})
