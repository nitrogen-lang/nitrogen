func testStringsTypes() {
	let str1 = 'Hello,
World!';

	let str2 = "Hello,\nWorld!";

	if (str1 != str2) {
		println("ERROR: Expected raw string and interpreted string to be the same.");
        println("str1: \"", str1, "\"");
        println("str2: \"", str2, "\"");
	}
}

testStringsTypes();
