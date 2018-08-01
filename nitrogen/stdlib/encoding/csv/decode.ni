import "stdlib/file"

const DEFAULT_DELIM = ','
const DEFAULT_QUOTE = '"'

class lexer {
    let source
    let curChar
    let peekChar
    let delimiter = DEFAULT_DELIM
    let quote = DEFAULT_QUOTE

    func init(file) {
        this.source = file
        this.readChar()
        this.readChar()
    }

    func setDelimiter(c) {
        this.delimiter = c
    }

    func setQuote(c) {
        this.quote = c
    }

    func readChar() {
        this.curChar = this.peekChar
        this.peekChar = file.readChar(this.source)
    }

    func readQuotedString() {
        let str = ""

        this.readChar() // Move pass open quote
        for {
            if this.curChar == this.quote {
                if this.peekChar == this.quote {
                    this.readChar()
                } else {
                    break
                }
            }

            str += this.curChar
            this.readChar()
        }
        this.readChar() // Move pass close quote

        return str
    }

    func readToDelim() {
        let str = ""

        while this.curChar != this.delimiter and this.curChar != "\n" {
            str += this.curChar
            this.readChar()
        }

        return str
    }

    func readField() {
        if this.curChar == this.quote: return this.readQuotedString()
        return this.readToDelim()
    }

    func readRecord() {
        if isNull(this.curChar): return nil

        let fields = []

        for {
            const field = this.readField()
            fields = push(fields, field)
            if isNull(this.curChar): break
            if this.curChar == this.delimiter: this.readChar()
            if this.curChar == "\n" {
                this.readChar()
                break
            }
        }

        return fields
    }
}

class fileReader {
    let l

    func init(f) {
        if !isResource(f) or resourceID(f) != file.fileResourceID {
            throw "fileReader expected a file resource"
        }

        this.l = new lexer(f)
    }

    func readRecord() {
        return this.l.readRecord()
    }

    func readAllRecords() {
        let records = [];

        let record = this.l.readRecord()
        while !isNull(record) {
            records = push(records, record)
            record = this.l.readRecord()
        }

        return records
    }

    func delimiter(c) {
        this.l.setDelimiter(c)
    }

    func quote(c) {
        this.l.setQuote(c)
    }
}

return {
    "fileReader": fileReader,
}
