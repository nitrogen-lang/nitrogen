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

test.run("file.readFile", fn(assert, check) {
    const data = file.readFile(filename)
    const expected = "Hello, world!\n"
    check(assert.isEq(data, expected))
})
