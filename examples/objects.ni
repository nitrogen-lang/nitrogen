/*
 * Copyright (c) 2017, Lee Keitel
 * This file is released under the BSD 3-Clause license.
 *
 * This file demonstrates "object-oriented" programming
 * using hashmaps.
 */

const account = func(name) {
    // The "object" is declared inside the function making it unique
    // Since oBalance is not in the dispatch map, it can't be modified
    // except through the exposed functions.
    let oBalance = 0
    let dispatch = {
        "withdraw": nil,
        "deposit": nil,
        "balance": nil,
    }

    // "methods" are added as key pairs to the map
    dispatch.withdraw = func(amount) {
        if amount > oBalance {
            return 'Insufficent balance'
        }
        oBalance = oBalance - amount
        return oBalance
    }

    dispatch.deposit = func(amount) {
        oBalance = oBalance + amount
        return oBalance
    }

    dispatch.balance = func() { oBalance }
    dispatch.name = name
    return dispatch
}

const main = func() {
    let me = account("John Smith")
    println(me.name)
    println(me.balance())
    println(me.deposit(200)) // "methods" are invoked by simply looking up the key in the map
    println(me.withdraw(50))
    println(me.withdraw(250))
    println(me.balance())

    // Multiple invocations of account() create new, unique "objects"
    let me2 = account("John Smith2")
    println(me2.name)
    println(me2.balance())
    println(me2.deposit(90))

    // Show original account was not modified
    println(me.name)
    println(me.balance())
}

main()
