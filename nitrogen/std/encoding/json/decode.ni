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

    const init = fn(t, v) {
        this.type = t
        this.value = v
    }
}

class lexer {
    let source
    let curIndex
    let curChar

    const init = fn(str) {
        this.source = str
        this.curIndex = 0
        this.readChar()
    }

    const readChar = fn() {
        if this.curIndex == len(this.source) + 1 {
            this.curChar = nil
            return
        }

        this.curChar = this.source[this.curIndex]
        this.curIndex += 1
    }

    const readKeyword = fn() {
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

    const readString = fn() {
        let str = ""

        this.readChar() // Move pass open quote
        while this.curChar != '"' {
            str += this.curChar
            this.readChar()
        }
        this.readChar() // Move pass close quote

        new token(STRING, str)
    }

    const readNumber = fn() {
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

    const nextToken = fn() {
        this.skipWhitespace()
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

    const skipWhitespace = fn() {
        while this.isWhitespace(this.curChar) {
            this.readChar()
        }
    }

    const isDigit = fn(c) { (c >= "0" and c <= "9") or c == "." or c == "-" }
    const isLetter = fn(c) { (c >= "a" and c <= "z") or (c >= "A" and c <= "Z") }
    const beginsKeyword = fn(c) { c == "f" or c == "t" or c == "n" }
    const isWhitespace = fn(c) { c == " " or c == "\t" or c == "\r" or c == "\n" }
}

class parser {
    let lexer
    let curToken

    const init = fn(l) {
        this.lexer = l
        this.nextToken()
    }

    const nextToken = fn() {
        this.curToken = this.lexer.nextToken()
    }

    const parse = fn() {
        const ct = this.curToken.type

        if ct == LCURLY: return this.parseObject()
        if ct == LSQUARE: return this.parseArray()
        if ct == TRUE: return this.curToken.value
        if ct == FALSE: return this.curToken.value
        if ct == NULL: return this.curToken.value
        if ct == STRING: return this.curToken.value
        if ct == NUMBER: return this.curToken.value

        throw "Invalid JSON, expected { [ true false null \" or a number"
    }

    const parseArray = fn() {
        this.nextToken()
        let arr = []

        loop {
            if this.curToken.type == RSQUARE: break
            arr = push(arr, this.parse())
            this.nextToken()

            if this.curToken.type == RSQUARE: break
            if this.curToken.type != COMMA: throw "Invalid JSON array, expected a comma"
            this.nextToken()
        }

        arr
    }

    const parseObject = fn() {
        this.nextToken()
        let obj = {}

        loop {
            if this.curToken.type == RCURLY: break
            if this.curToken.type != STRING: throw "Invalid JSON object key, expected a string"
            const key = this.curToken.value

            this.nextToken()
            if this.curToken.type != COLON: throw "Invalid JSON object value pair, expected a colon"

            this.nextToken()
            obj[key] = this.parse()

            this.nextToken()
            if this.curToken.type == RCURLY: break
            if this.curToken.type != COMMA: throw "Invalid JSON object, expected a comma"
            this.nextToken()
        }

        obj
    }
}

const decode = fn(str) {
    const l = new lexer(str)
    const p = new parser(l)
    p.parse()
}

return {
    "decode": decode,
}
