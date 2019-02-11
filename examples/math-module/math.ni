const exports = {}

exports.add = fn(a, b) { a + b }
exports.sub = fn(a, b) { a - b }
exports.mul = fn(a, b) { a * b }
exports.div = fn(a, b) { a / b }

return exports
