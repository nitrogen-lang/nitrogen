fn native readFile (path)
fn native remove (path)
fn native exists (path)
fn native rename (oldname, newname)
fn native dirlist (path)
fn native isdir (path)

class File {
    fn native init (path)
    fn native close ()
    fn native write (data)
    fn native readAll ()
    fn native readLine ()
    fn native readChar ()
    fn native remove ()
    fn native rename (newname)
}

return {
    readFile,
    remove,
    exists,
    rename,
    dirlist,
    isdir,
    File,
}
