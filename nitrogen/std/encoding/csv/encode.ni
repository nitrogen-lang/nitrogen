import "std/collections"
import "std/file"
import "std/string"

use collections.join
use collections.map
use string.contains
use string.replace

class fileWriter {
    let cfile
    let delimiter = ','
    let quote = '"'

    fn init(f) {
        if !isResource(f) or resourceID(f) != file.fileResourceID {
            throw "fileWriter expected a file resource"
        }
        this.cfile = f
    }

    fn writeRecord(record) {
        file.write(this.cfile, join(",", map(record, this.csvQuote)))
        file.write(this.cfile, "\n")
    }

    fn csvQuote(item) {
        item = toString(item)
        if contains(item, this.quote) {
            item = replace(item, this.quote, this.quote+this.quote, -1)
            return this.quote + item + this.quote
        }
        if contains(item, this.delimiter) {
            return this.quote + item + this.quote
        }
        item
    }
}

return {
    "fileWriter": fileWriter,
}
