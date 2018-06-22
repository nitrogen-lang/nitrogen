import './another.ni' as otherFile

func main() {
    otherFile()
}

println("Calling main() from ", _FILE)
main()
