func testStringsTypes() {
	let str1 = 'Hello,
World!'

	let str2 = "Hello,\nWorld!"

	if str1 != str2 {
		println("ERROR: Expected raw string and interpreted string to be the same.")
        println("str1: \"", str1, "\"")
        println("str2: \"", str2, "\"")
	}
}
testStringsTypes();

func testStringTypes2() {
	let str1 = "Hello, world!"
	let expected = "Hello, World!"

	str1[7] = "W"

	if str1 != expected {
		println("String Test Failed: Expected ", expected, ", got ", str1)
	}
}
testStringTypes2();

func testStringTypes3() {
	const str1 = "Hello, 世界!"
	const expected = "世"

	const indexed = str1[7]

	if indexed != expected {
		println("String Test Failed: Expected ", expected, ", got ", indexed)
	}
}
testStringTypes3();
