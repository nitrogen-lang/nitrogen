import "std/test"
import "std/collections" as col

class MyCollectionIter {
    let collection
    let i = 0

    const init = fn(collection) {
        this.collection = collection
    }

    const _next = fn() {
        if this.i == len(this.collection.items) {
            return nil
        }

        const i = this.i
        const item = this.collection.items[i]
        this.i += 1
        return [i, item]
    }
}


class MyCollection {
    let items

    const init = fn(items) {
        this.items = items
    }

    const _iter = fn() {
        new MyCollectionIter(this)
    }
}

test.run("Class iterator", fn(assert, check) {
    const input = ["computer", "mouse", "keyboard"]
    const things = new MyCollection(input)

    let copy = []
    for thing in things {
        copy = push(copy, thing)
    }

    check(assert.isTrue(col.arrayMatch(input, copy)))
})
