import "std/collections" as col
import "std/encoding/csv"
import "std/file"
import "std/filepath"
import "std/test"

const testdataDir = _ENV['TESTDATA_DIR']
if isNil(testdataDir) {
    println("TESTDATA_DIR not set")
    exit(1)
}

const filename = filepath.join(testdataDir, 'test.csv')

const csvFile = file.open(filename, 'r')
const reader = new csv.fileReader(csvFile)

const records = reader.readAllRecords()

file.close(csvFile)

test.run("Check row count", func(assert) {
    assert.isTrue(isArray(records))
    assert.isEq(len(records), 6)
})

test.run("Check field count", func(assert) {
    col.foreach(records, func(i, v) {
        assert.isEq(7, len(v))
    });
})

test.run("Check field with quotes", func(assert) {
    const row = records[1]
    assert.isEq('"GREEN"', row[6])
})

test.run("Check field with newline", func(assert) {
    const row = records[5]
    assert.isEq("Lettie,\nLopez", row[1])
})
