import "std/preamble/os"
import "std/preamble/io"
import "std/preamble/collection"

fn native toInt()
fn native toFloat()
fn native toString()
fn native toByteString()
fn native parseInt()
fn native parseFloat()
fn native varType()
fn native isDefined()
fn native isFloat()
fn native isInt()
fn native isBool()
fn native isNull()
fn native isNil()
fn native isFunc()
fn native isString()
fn native isByteString()
fn native isArray()
fn native isMap()
fn native isError()
fn native isException()
fn native isResource()
fn native isClass()
fn native isInstance()
fn native errorVal()
fn native resourceID()
fn native modulesSupported()
fn native error()
fn native instanceOf()
fn native classOf()

fn strToArray(str) {
    let result = []
    for char in str {
        result = push(result, char)
    }
    return result
}

return {
    toInt,
    toFloat,
    toString,
    toByteString,
    parseInt,
    parseFloat,
    varType,
    isDefined,
    isFloat,
    isInt,
    isBool,
    isNull,
    isNil,
    isFunc,
    isString,
    isByteString,
    isArray,
    isMap,
    isError,
    isException,
    isResource,
    isClass,
    isInstance,
    errorVal,
    resourceID,
    modulesSupported,
    error,
    instanceOf,
    classOf,

    strToArray,

    "exit": os.exit,

    "print": io.print,
    "printlnb": io.printlnb,
    "println": io.println,
    "printerr": io.printerr,
    "printerrln": io.printerrln,
    "printenv": io.printenv,
    "varDump": io.varDump,
    "readline": io.readline,

    "len": collection.len,
    "first": collection.first,
    "last": collection.last,
    "rest": collection.rest,
    "pop": collection.pop,
    "push": collection.push,
    "prepend": collection.prepend,
    "splice": collection.splice,
    "slice": collection.slice,
    "sort": collection.sort,
    "hashMerge": collection.hashMerge,
    "hashKeys": collection.hashKeys,
    "hasKey": collection.hasKey,
    "range": collection.range,
}
