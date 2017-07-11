let string = "Hello, world!"

func change_string(next) {
     string = next
}

if (string != "Hello, world!") {
   println("Test 1 failed!")
   return;
}

change_string("Hello, mars!")

if (string != "Hello, mars!") {
   println("Test 2 failed!")
}
