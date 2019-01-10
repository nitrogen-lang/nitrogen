import "std/test"
import "std/string"

use string.String
use string.String as str

test.run("Use statement", func(assert) {
    assert.isTrue(isDefined("String"))
    assert.isTrue(isDefined("str"))

    assert.isTrue(isClass(String))
    assert.isTrue(isClass(str))
})
