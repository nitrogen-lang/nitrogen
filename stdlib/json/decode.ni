// Token Types:
const INVALID = "INVALID"
const LCURLY = "LCURLY"
const RCURLY = "RCURLY"
const COLON = "COLON"
const COMMA = "COMMA"
const LSQUARE = "LSQUARE"
const RSQUARE = "RSQUARE"
const TRUE = "TRUE"
const FALSE = "FALSE"
const NULL = "NULL"
const STRING = "STRING"
const NUMBER = "NUMBER"

class token {
    let type
    let value

    func init(t, v) {
        this.type = t
        this.value = v
    }
}

class lexer {
    let source
    let curIndex
    let curChar

    func init(str) {
        this.source = str
        this.curIndex = 0
        this.readChar()
    }

    func readChar() {
        if this.curIndex == len(this.source) + 1 {
            this.curChar = nil
            return
        }

        this.curChar = this.source[this.curIndex]
        this.curIndex += 1
    }

    func readObject() {

    }

    func readArray() {

    }

    func readIdent() {

    }

    func readString() {
        let str = ""

        for {
            this.readChar()

            if this.curChar == '"' {
                break
            }

            str += this.curChar
        }
        this.readChar()

        make token(STRING, str)
    }

    func readNumber() {

    }

    func nextToken() {
        if this.curChar == '"' {
            return this.readString()
        }
        make token(INVALID, "")
    }
}

class parser {
    let lexer
    let curToken

    func init(l) {
        this.lexer = l
    }

    func nextToken() {
        this.curToken = this.lexer.nextToken()
    }

    func parse() {
        this.nextToken()

        if this.curToken.type == INVALID {
            throw "Invalid JSON"
        }

        this.curToken.value
    }
}

func decode(str, str2) {
    const l = make lexer(str)
    const p = make parser(l)
    p.parse()
}

return {
    "decode": decode,
}
