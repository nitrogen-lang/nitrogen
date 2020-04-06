import "std/collections" as col
import "std/encoding/csv"
import "std/file"
import "std/filepath"
import "std/test"
import "std/os"

const testdataDir = os.env()['TESTDATA_DIR']
if isNil(testdataDir) {
    println("TESTDATA_DIR not set")
    exit(1)
}

const decodeFilename = filepath.join(testdataDir, 'test.csv')

const decodeCsvFile = file.open(decodeFilename, 'r')
const reader = new csv.fileReader(decodeCsvFile)

const records = reader.readAllRecords()

file.close(decodeCsvFile)

test.run("Check row count", fn(assert) {
    assert.isTrue(isArray(records))
    assert.isEq(len(records), 6)
})

test.run("Check field count", fn(assert) {
    col.foreach(records, fn(i, v) {
        assert.isEq(7, len(v))
    });
})

test.run("Check field with quotes", fn(assert) {
    const row = records[1]
    assert.isEq('"GREEN"', row[6])
})

test.run("Check field with newline", fn(assert) {
    const row = records[5]
    assert.isEq("Lettie,\nLopez", row[1])
})


const encodeFilename = filepath.join(testdataDir, "tmp-test.csv")
const encodeCsvFile = file.open(encodeFilename, "w")
const writer = new csv.fileWriter(encodeCsvFile)

const newRecords = [
    ["seq","name","age","state","zip","dollar","pick"],
    [1,"Keith,Jackson",27,"MN",81521,"$354.79",'"GREEN"'],
    [2,"Frances,Wheeler",58,"NY",21838,"$1322.39",'"YELLOW"'],
    [3,"Miguel,Hopkins",35,"GA",91111,"$522.29",'"WHITE"'],
    [4,"Noah,Spencer",22,"DE",94024,"$8178.92",'"GREEN"'],
    [5,"Lettie,\nLopez",64,"RI",39463,"$6219.12",'"GREEN"'],
]

for record in newRecords {
    writer.writeRecord(record)
}

file.close(encodeCsvFile)

test.run("Check encoded file was written", fn(assert) {
    const fileData = file.readAll(encodeFilename)
    assert.isEq(fileData, 'seq,name,age,state,zip,dollar,pick
1,"Keith,Jackson",27,MN,81521,$354.79,"""GREEN"""
2,"Frances,Wheeler",58,NY,21838,$1322.39,"""YELLOW"""
3,"Miguel,Hopkins",35,GA,91111,$522.29,"""WHITE"""
4,"Noah,Spencer",22,DE,94024,$8178.92,"""GREEN"""
5,"Lettie,
Lopez",64,RI,39463,$6219.12,"""GREEN"""
')
}, fn() {
    file.remove(encodeFilename)
})
