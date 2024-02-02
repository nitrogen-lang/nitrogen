fn native contains(s, substr)
fn native count(s, substr)
fn native dedup(s, char)
fn native format(s)
fn native hasPrefix(s, prefix)
fn native hasSuffix(s, suffix)
fn native replace(s, old, str, n)
fn native split(s, sep)
fn native splitN(s, sep, n)
fn native trimSpace(s)

class String {
    let str = ""

    fn init(s) {
        this.str = s
    }

    fn native contains(substr)
    fn native count(substr)
    fn native dedup(char)
    fn native format()
    fn native hasPrefix(prefix)
    fn native hasSuffix(suffix)
    fn native replace(old, str, n)
    fn native split(sep)
    fn native splitN(sep, n)
    fn native trimSpace()
}

return {
    "contains": contains,
    "count": count,
    "dedup": dedup,
    "format": format,
    "hasPrefix": hasPrefix,
    "hasSuffix": hasSuffix,
    "replace": replace,
    "split": split,
    "splitN": splitN,
    "trimSpace": trimSpace,
    "String": String,
}
