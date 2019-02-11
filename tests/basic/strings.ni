import "std/test"

test.run("Check equality of raw and escaped strings", fn(assert) {
	const str1 = 'Hello,
World!'
	const str2 = "Hello,\nWorld!"
	assert.isEq(str1, str2)
})

test.run("Check changing string index", fn(assert) {
	const str1 = "Hello, world!"
	const expected = "Hello, World!"

	str1[7] = "W"

	assert.isEq(str1, expected)
})

test.run("Indexing UTF-8 string", fn(assert) {
	const str1 = "Hello, 世界!"
	const expected = "世"
	assert.isEq(str1[7], expected)
})
