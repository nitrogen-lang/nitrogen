let otherFile = import('./another.ni')

func main() {
    otherFile()
}

println("Calling main() from ", _FILE)
main()
