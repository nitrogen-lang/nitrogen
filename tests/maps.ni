/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates looping over a hash map.
 */

func arrayHasValue(array, value) {
    for (i = 0; i < len(array); i += 1) {
        if (array[i] == value) {
	        return true
	    }
    }
    return false
}

func main() {
    let hash = {
        "key1": "value1",
        "key2": "value2",
        "key3": "value3",
        "key4": "value4",
    }

    let keys = hashKeys(hash);
    let expectedKeys = ["key1", "key2", "key3", "key4"]

    if (len(keys) != len(expectedKeys)) {
        println("Not enough keys. Expected ", len(expectedKeys), " got ", len(keys))
    }

    for (i = 0; i < len(expectedKeys); i += 1) {
        if (!arrayHasValue(keys, expectedKeys[i])) {
            println("hash keys doesn't have ", expectedKeys[i])
            return
        }
    }
}

main()
