import './another.ni' as otherFile

const main = fn() {
    otherFile()
}

println("Calling main() from ", _FILE)
main()
