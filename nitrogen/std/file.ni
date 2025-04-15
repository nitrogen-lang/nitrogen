export fn native readFile (path)
export fn native remove (path)
export fn native exists (path)
export fn native rename (oldname, newname)
export fn native dirlist (path)
export fn native isdir (path)

export class File {
    fn native init (path)
    fn native close ()
    fn native write (data)
    fn native readAll ()
    fn native readLine ()
    fn native readChar ()
    fn native remove ()
    fn native rename (newname)
}
