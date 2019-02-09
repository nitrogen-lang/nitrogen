# csv.ni

Read and write CSV (or similarly) encoded data.

To use: `import 'std/encoding/csv'`

## class fileReader(f: resource)

A fileReader takes a [file](file.ni.md) resource.

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

const csvFile = file.open(data.csv, 'r')
const reader = new csv.fileReader(csvFile)

let record = reader.readRecord()
while !isNull(record) {
    println(record)
    record = reader.readRecord()
}

file.close(csvFile)
```
