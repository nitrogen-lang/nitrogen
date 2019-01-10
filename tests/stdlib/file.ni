import "std/file"
import "std/filepath"
import "std/test"

const testdataDir = _ENV['TESTDATA_DIR']
if isNil(testdataDir) {
    println("TESTDATA_DIR not set")
    exit(1)
}

const filename = filepath.join(testdataDir, 'test.txt')

test.run("file.readAll", func(assert) {
    const data = file.readAll(filename)
    const expected = "Hello, world!\n"
    assert.isEq(data, expected)
})
