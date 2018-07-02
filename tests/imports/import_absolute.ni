import "test"

test.run("Non-existant import", func(assert) {
    assert.shouldThrow(func() {
        import './includes/_not_exist.ni'
    })
})

test.run("Absolute import", func(assert) {
    import '../../testdata/math.ni'
    assert.isTrue(isFunc(math.add))
    assert.isEq(math.add(2, 4), 6)
})
