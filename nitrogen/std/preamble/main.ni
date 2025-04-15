import "std/preamble/os"
import "std/preamble/io"
import "std/preamble/collection"

export fn native toInt()
export fn native toFloat()
export fn native toString()
export fn native toByteString()
export fn native parseInt()
export fn native parseFloat()
export fn native varType()
export fn native isDefined()
export fn native isFloat()
export fn native isInt()
export fn native isBool()
export fn native isNull()
export fn native isNil()
export fn native isFunc()
export fn native isString()
export fn native isByteString()
export fn native isArray()
export fn native isMap()
export fn native isError()
export fn native isException()
export fn native isResource()
export fn native isClass()
export fn native isInstance()
export fn native isModule()
export fn native errorVal()
export fn native resourceID()
export fn native modulesSupported()
export fn native error()
export fn native instanceOf()
export fn native classOf()

export fn strToArray(str) {
    let result = []
    for char in str {
        result = push(result, char)
    }
    return result
}

export const exit = os.exit

export const print = io.print
export const printlnb = io.printlnb
export const println = io.println
export const printerr = io.printerr
export const printerrln = io.printerrln
export const printenv = io.printenv
export const varDump = io.varDump
export const readline = io.readline

export const len = collection.len
export const first = collection.first
export const last = collection.last
export const rest = collection.rest
export const pop = collection.pop
export const push = collection.push
export const prepend = collection.prepend
export const splice = collection.splice
export const slice = collection.slice
export const sort = collection.sort
export const hashMerge = collection.hashMerge
export const hashKeys = collection.hashKeys
export const hasKey = collection.hasKey
export const range = collection.range
