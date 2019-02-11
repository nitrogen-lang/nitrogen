import "std/test"

test.run("Non-existant import", fn(assert) {
    assert.shouldThrow(fn() {
        import './includes/_not_exist.ni'
    })
})

test.run("Absolute import", fn(assert) {
    import '../../testdata/math.ni'
    assert.isTrue(isFunc(math.add))
    assert.isEq(math.add(2, 4), 6)
})
