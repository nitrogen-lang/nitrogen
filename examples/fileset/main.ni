import './another.ni' as otherFile

const main = func() {
    otherFile()
}

println("Calling main() from ", _FILE)
main()
