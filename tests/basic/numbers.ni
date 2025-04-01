import "std/test"

test.run("Check numeric separator", fn(assert, check) {
	const num_no_separators = 10000000
	const num_with_separators1 = 10_000_000
	const num_with_separators2 = 10__000__000
	const num_with_separators3 = 10_000_000_
	check(assert.isEq(num_no_separators, num_with_separators1))
	check(assert.isEq(num_no_separators, num_with_separators2))
	check(assert.isEq(num_no_separators, num_with_separators3))

	const float_no_separators = 10000.00
	const float_with_separators = 10_000.00
	check(assert.isEq(float_no_separators, float_with_separators))

	const hex_no_separators = 0xAB12EF58
	const hex_with_separators1 = 0xAB_12_EF_58
	const hex_with_separators2 = 0x_AB_12_EF_58
	check(assert.isEq(hex_no_separators, hex_with_separators1))
	check(assert.isEq(hex_no_separators, hex_with_separators2))
})

test.run("Test hex literal", fn(assert, check) {
	const dec = 10
	const hex = 0x0A
	check(assert.isEq(dec, hex))
})

test.run("Test binary literal", fn(assert, check) {
	const dec = 10
	const bin = 0b1010
	check(assert.isEq(dec, bin))
})

test.run("Test octal literal", fn(assert, check) {
	const dec = 10
	const oct = 0o12
	check(assert.isEq(dec, oct))
})
