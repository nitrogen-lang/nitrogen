# csv.ni

Read and write CSV encoded data.

To use: `import 'std/encoding/csv'`

## class reader(f: resource)

A reader takes a [file](file.ni.md) resource. This class implementes the
iterator interface so it can be used in for..in loops.

### Methods

#### readRecord(): array|nil

Reads a single record from the file. The record is returned as an array of strings.
`readRecord()` will return nil if there's no more data to return.

#### delimiter(c: string)

Sets the delimiter or field separator. Defaults to comma `,`.

#### quote(c: string)

Sets the field quote character. Defaults to double quote `"`.

## Example

```
import "std/encoding/csv"
import "std/file"

const csvFile = new file.File('data.csv', 'r')
const reader = new csv.reader(csvFile)

for record in reader {
    println(record)
}

file.close(csvFile)
```

## class writer(f: resource)

A writer takes a [file](file.ni.md) resource.

### Fields

#### delimiter: string = ','

Field separator. Defaults to comma `,`.

#### quote: string = '"'

Field quote character. Defaults to double quote `"`.

### Methods

#### writeRecord(record: array): int

Write a record to the file. Fields will be quoted if needed.

## Example

```
import "std/encoding/csv"
import "std/file"

const csvFile = new file.File('data.csv', 'w')
const writer = new csv.writer(csvFile)

const records = [
    ["seq","name","age","state","zip","dollar","pick"],
    [1,"Keith,Jackson",27,"MN",81521,"$354.79",'"GREEN"'],
    [2,"Frances,Wheeler",58,"NY",21838,"$1322.39",'"YELLOW"'],
    [3,"Miguel,Hopkins",35,"GA",91111,"$522.29",'"WHITE"'],
    [4,"Noah,Spencer",22,"DE",94024,"$8178.92",'"GREEN"'],
    [5,"Lettie,\nLopez",64,"RI",39463,"$6219.12",'"GREEN"'],
]

for record in records {
    writer.writeRecord(record)
}

file.close(csvFile)
```

## interface WriterIfc

### Fields

#### write(data)

## interface ReaderIfc

### Fields

#### readChar(): string
