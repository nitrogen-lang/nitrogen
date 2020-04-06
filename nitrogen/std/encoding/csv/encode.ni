import "std/collections"
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
        this.cfile = f
    }

    fn writeRecord(record) {
        this.cfile.write(join(",", map(record, this.csvQuote)))
        this.cfile.write("\n")
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
