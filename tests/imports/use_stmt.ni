import "std/test"
import "std/string"

use string.String
use string.String as str

test.run("Use statement", fn(assert, check) {
    check(assert.isTrue(isDefined("String")))
    check(assert.isTrue(isDefined("str")))

    check(assert.isTrue(isClass(String)))
    check(assert.isTrue(isClass(str)))
})
