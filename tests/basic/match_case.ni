import "std/test"

test.run("Match expression", fn(assert, check) {
    const item = "thing"

    const item2 = match (item) {
        "thing" => "thing2",
        "other thing" => "other thing2",
        _ => "",
    }

    check(assert.isEq(item2, "thing2"))
})

test.run("Match expression no parens", fn(assert, check) {
    const item = "thing"

    const item2 = match item {
        "thing" => "thing2",
        "other thing" => "other thing2",
        _ => "",
    }

    check(assert.isEq(item2, "thing2"))
})

test.run("Match expression default case", fn(assert, check) {
    const item = "thing3"

    const item2 = match (item) {
        "thing" => "thing2",
        "other thing" => "other thing2",
        _ => "default",
    }

    check(assert.isEq(item2, "default"))
})

test.run("Match expression 2 branches", fn(assert, check) {
    const item = "thing"

    const item2 = match (item) {
        "thing" => "thing2",
        _ => "default",
    }

    check(assert.isEq(item2, "thing2"))
})

test.run("Match expression 1 branch", fn(assert, check) {
    const item = "thing"

    const item2 = match (item) {
        "thing" => "thing2",
    }

    check(assert.isEq(item2, "thing2"))
})

test.run("Match expression 1 branch default", fn(assert, check) {
    const item = "thing"

    const item2 = match (item) {
        _ => "default",
    }

    check(assert.isEq(item2, "default"))
})
