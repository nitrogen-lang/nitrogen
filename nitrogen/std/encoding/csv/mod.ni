import 'std/encoding/csv/encode'
import 'std/encoding/csv/decode'

return {
    "writer": encode.fileWriter,
    "WriterIfc": encode.Writer,
    "reader": decode.reader,
    "ReaderIfc": decode.CharReader,
}
