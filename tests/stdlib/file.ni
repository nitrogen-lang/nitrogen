import "std/file"
import "std/filepath"
import "std/test"
import "std/os"

const testdataDir = os.env()['TESTDATA_DIR']
if isNil(testdataDir) {
    println("TESTDATA_DIR not set")
    exit(1)
}

const filename = filepath.join(testdataDir, 'test.txt')

test.run("file.readAll", fn(assert) {
    const data = file.readAll(filename)
    const expected = "Hello, world!\n"
    assert.isEq(data, expected)
})
