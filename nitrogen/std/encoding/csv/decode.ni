const DEFAULT_DELIM = ','
const DEFAULT_QUOTE = '"'

interface CharReader {
    readChar()
}

class lexer {
    let source
    let curChar
    let peekChar
    let delimiter = DEFAULT_DELIM
    let quote = DEFAULT_QUOTE

    const init = fn(file) {
        if ! file implements CharReader {
            throw "f must be a CharReader"
        }

        this.source = file
        this.readChar()
        this.readChar()
    }

    const setDelimiter = fn(c) {
        this.delimiter = c
    }

    const setQuote = fn(c) {
        this.quote = c
    }

    const readChar = fn() {
        this.curChar = this.peekChar
        this.peekChar = this.source.readChar()
    }

    const readQuotedString = fn() {
        let str = ""

        this.readChar() // Move pass open quote
        loop {
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

    const readToDelim = fn() {
        let str = ""

        while this.curChar != this.delimiter and this.curChar != "\n" {
            str += this.curChar
            this.readChar()
        }

        return str
    }

    const readField = fn() {
        if this.curChar == this.quote: return this.readQuotedString()
        return this.readToDelim()
    }

    const readRecord = fn() {
        if isNull(this.curChar): return nil

        let fields = []

        loop {
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

    const init = fn(f) {
        this.l = new lexer(f)
    }

    const readRecord = fn() {
        return this.l.readRecord()
    }

    const readAllRecords = fn() {
        let records = [];

        let record = this.l.readRecord()
        while !isNull(record) {
            records = push(records, record)
            record = this.l.readRecord()
        }

        return records
    }

    const delimiter = fn(c) {
        this.l.setDelimiter(c)
    }

    const quote = fn(c) {
        this.l.setQuote(c)
    }
}

return {
    "fileReader": fileReader,
    "CharReader": CharReader,
}
