import "std/test"

test.run("Check equality of raw and escaped strings", fn(assert, check) {
	const str1 = 'Hello,
World!'
	const str2 = "Hello,\nWorld!"
	check(assert.isEq(str1, str2))
})

test.run("Check changing string index", fn(assert, check) {
	const str1 = "Hello, world!"
	const expected = "Hello, World!"

	str1[7] = "W"

	check(assert.isEq(str1, expected))
})

test.run("Indexing UTF-8 string", fn(assert, check) {
	const str1 = "Hello, 世界!"
	const expected = "世"
	check(assert.isEq(str1[7], expected))
})

test.run("String iterator", fn(assert, check) {
	const str1 = "Hello, 世界!"
	const expected = ["H", "e", "l", "l", "o", ",", " ", "世", "界", "!"]
	let result = []

	for char in str1 {
		result = push(result, char)
	}

	check(assert.isEq(result, expected))
})

test.run("Byte string iterator", fn(assert, check) {
	const str1 = b"Hello, 世界!"
	const expected = [72,101,108,108,111,44,32,228,184,150,231,149,140,33]
	let result = []

	for char in str1 {
		result = push(result, toInt(char))
	}

	check(assert.isEq(result, expected))
})

test.run("Byte string hex literals", fn(assert, check) {
	const str1 = b"Hello, \xE4\xB8\x96\xE7\x95\x8C!"
	const expected = [72,101,108,108,111,44,32,228,184,150,231,149,140,33]
	let result = []

	for char in str1 {
		result = push(result, toInt(char))
	}

	check(assert.isEq(result, expected))
})

test.run("Byte strings", fn(assert, check) {
	const str1 = b"Hello"
	check(assert.isEq(str1[3], b"l"))
})

test.run("Byte strings concat", fn(assert, check) {
	const str1 = b"Hello"
	const str2 = b", World"
	check(assert.isEq(str1 + str2, b"Hello, World"))
})

test.run("Byte strings index unicode", fn(assert, check) {
	const str1 = b"Hello, 世界!"
	const expected = 228
	check(assert.isEq(toInt(str1[7]), expected))
})

test.run("Check changing byte string index", fn(assert, check) {
	const str1 = b"Hello, world!"
	const expected = b"Hello, World!"

	str1[7] = b"W"

	check(assert.isEq(str1, expected))
})

test.run("Check string type conversions", fn(assert, check) {
	const str1 = b"Hello, world!"
	const str2 = "Hello, world!"

	check(assert.isEq(toString(str1), str2))
	check(assert.isEq(str1, toByteString(str2)))
})
