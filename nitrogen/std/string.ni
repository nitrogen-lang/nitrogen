export fn native contains(s, substr)
export fn native count(s, substr)
export fn native dedup(s, char)
export fn native format(s)
export fn native hasPrefix(s, prefix)
export fn native hasSuffix(s, suffix)
export fn native replace(s, old, str, n)
export fn native split(s, sep)
export fn native splitN(s, sep, n)
export fn native trimSpace(s)

export class String {
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
