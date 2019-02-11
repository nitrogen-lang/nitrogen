import "std/test"
import "std/string"

use string.String
use string.String as str

test.run("Use statement", fn(assert) {
    assert.isTrue(isDefined("String"))
    assert.isTrue(isDefined("str"))

    assert.isTrue(isClass(String))
    assert.isTrue(isClass(str))
})
