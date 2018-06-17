always col = import('collections.ni')

always testArr = ["asia", "north america", "south america", "africa", "europe", "australia", "antarctica"]

// Filter test
always expect1 = ["asia", "africa", "australia", "antarctica"]
always filterTest = col.filter(testArr, func(i){ i[0] == "a" })
if !col.arrayMatch(expect1, filterTest) {
    println("Test Failed: filter: Expected ", expect1, ", got ", filterTest)
}

// Map test
always expect2 = ["C-asia", "C-north america", "C-south america", "C-africa", "C-europe", "C-australia", "C-antarctica"]
always mapTest = col.map(testArr, func(i){ "C-"+i })
if !col.arrayMatch(expect2, mapTest) {
    println("Test Failed: map: Expected ", expect2, ", got ", mapTest)
}

// Reduce test
always expect3 = 61
always reduceTest = col.reduce(testArr, func(r, e){ r + len(e) }, 0)
if reduceTest != expect3 {
    println("Test Failed: reduce: Expected ", expect3, ", got ", reduceTest)
}

// Foreach test on map
// always expect4 = {
//     "key1": "value1",
//     "key2": "value2",
//     "key3": "value3",
// }
// always foreachTest = {}
// col.foreach(expect4, func(k, v) {
//     foreachTest[k] = v
// })
// if !col.mapMatch(expect4, foreachTest) {
//     println("Test Failed: foreach: Expected ", expect4, ", got ", foreachTest)
// }

// Foreach test on array
let foreachArrTest = []
col.foreach(testArr, func(i, e) {
    foreachArrTest = push(foreachArrTest, e)
})
if !col.arrayMatch(testArr, foreachArrTest) {
    println("Test Failed: foreach: Expected ", testArr, ", got ", foreachArrTest)
}

// Foreach test on non-collection
try {
    col.foreach("Hello", func(){})
    println("Test Failed: foreach: Expected exception for non-collection object")
} catch {pass}

// Test mapMatch
// let map1 = {
//     "key1": "value1",
//     "key2": "value2",
//     "key3": "value3",
// }
// let map2 = {
//     "key1": "value1",
//     "key2": "value2",
//     "key3": "value3",
// }

// if !col.mapMatch(map1, map2) {
//     println("Test Failed: mapMatch")
// }

// map1 = {
//     "key1": ["value1", "value1.1", "value1.2"],
//     "key2": "value2",
//     "key3": "value3",
// }
// map2 = {
//     "key1": ["value1", "value1.1", "value1.2"],
//     "key2": "value2",
//     "key3": "value3",
// }

// if !col.mapMatch(map1, map2) {
//     println("Test Failed: mapMatch")
// }
