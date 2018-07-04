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

    func readKeyword() {
        let str = ""

        while this.isLetter(this.curChar) {
            str += this.curChar
            this.readChar()
        }

        if str == "true" {
            return new token(TRUE, true)
        } elif str == "false" {
            return new token(FALSE, false)
        } elif str == "null" {
            return new token(NULL, nil)
        }

        new token(INVALID, "")
    }

    func readString() {
        let str = ""

        this.readChar() // Move pass open quote
        while this.curChar != '"' {
            str += this.curChar
            this.readChar()
        }
        this.readChar() // Move pass close quote

        new token(STRING, str)
    }

    func readNumber() {
        // TODO: Implement signs and e notation
        let str = ""
        let isFloat = false

        while this.isDigit(this.curChar) {
            if this.curChar == ".": isFloat = true
            str += this.curChar
            this.readChar()
        }

        const num = if isFloat {
            parseFloat(str)
        } else {
            parseInt(str)
        }

        if isNull(num): return (new token(INVALID, ""))
        new token(NUMBER, num)
    }

    func nextToken() {
        const char = this.curChar

        // Punctuation and delimiters
        if char == '{'{
            this.readChar()
            return new token(LCURLY, '{')
        }
        if char == '}'{
            this.readChar()
            return new token(RCURLY, '}')
        }
        if char == '['{
            this.readChar()
            return new token(LSQUARE, '[')
        }
        if char == ']'{
            this.readChar()
            return new token(RSQUARE, ']')
        }
        if char == ':'{
            this.readChar()
            return new token(COLON, ':')
        }
        if char == ','{
            this.readChar()
            return new token(COMMA, ',')
        }

        // Concrete tokens
        if char == '"': return this.readString()
        if this.isDigit(char): return this.readNumber()
        if this.beginsKeyword(char): return this.readKeyword()

        new token(INVALID, "")
    }

    func isDigit(c) { (c >= "0" and c <= "9") or c == "." }
    func isLetter(c) { (c >= "a" and c <= "z") or (c >= "A" and c <= "Z") }
    func beginsKeyword(c) { c == "f" or c == "t" or c == "n" }
}

class parser {
    let lexer
    let curToken

    func init(l) {
        this.lexer = l
        this.nextToken()
    }

    func nextToken() {
        this.curToken = this.lexer.nextToken()
    }

    func parse() {
        const ct = this.curToken.type

        if ct == LCURLY: return this.parseObject()
        if ct == LSQUARE: return this.parseArray()
        if ct == TRUE: return this.curToken.value
        if ct == FALSE: return this.curToken.value
        if ct == NULL: return this.curToken.value
        if ct == STRING: return this.curToken.value
        if ct == NUMBER: return this.curToken.value

        throw "Invalid JSON"
    }

    func parseArray() {
        this.nextToken()
        let arr = []

        for {
            if this.curToken.type == RSQUARE: break
            arr = push(arr, this.parse())
            this.nextToken()

            if this.curToken.type == RSQUARE: break
            if this.curToken.type != COMMA: throw "Invalid JSON array"
            this.nextToken()
        }

        arr
    }

    func parseObject() {
        this.nextToken()
        let obj = {}

        for {
            if this.curToken.type == RCURLY: break
            if this.curToken.type != STRING: throw "Invalid JSON object key"
            const key = this.curToken.value

            this.nextToken()
            if this.curToken.type != COLON: throw "Invalid JSON object value pair"

            this.nextToken()
            obj[key] = this.parse()

            this.nextToken()
            if this.curToken.type == RCURLY: break
            if this.curToken.type != COMMA: throw "Invalid JSON object"
            this.nextToken()
        }

        obj
    }
}

func decode(str) {
    const l = new lexer(str)
    const p = new parser(l)
    p.parse()
}

return {
    "decode": decode,
}
