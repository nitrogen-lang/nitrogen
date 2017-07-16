/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates "object-oriented" programming.
 * Nitrogen doesn't have objects as such, not yet anyway,
 * but it can still be achieved using hash maps and factory
 * methods like below.
 */

func account(name) {
    // The "object" is declared inside the function making it unique
    let dispatch = {
        "name": name,
        "withdraw": nil,
        "deposit": nil,
        "balance": 0,
    }

    // "methods" are added as key pairs to the map
    dispatch["withdraw"] = func(amount) {
        if (amount > dispatch["balance"]) {
            return 'Insufficent balance'
        }
        dispatch["balance"] = dispatch["balance"] - amount
        return dispatch["balance"]
    }

    dispatch["deposit"] = func(amount) {
        dispatch["balance"] = dispatch["balance"] + amount
        return dispatch["balance"]
    }

    return dispatch
}

func main() {
    let me = account("John Smith")
    println(me["balance"])
    println(me["deposit"](200)) // "methods" are invoked by simply looking up the key in the map
    println(me["withdraw"](50))
    println(me["withdraw"](250))
    println(me["balance"])

    // Multiple invocations of account() create new, unique "objects"
    let me2 = account("John Smith2")
    println(me2["balance"])
    println(me2["deposit"](90))
}

main()
