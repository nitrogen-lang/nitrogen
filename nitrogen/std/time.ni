export fn native now()
export fn native now_ms()
export fn native now_ns()
export fn native now_unix()

export class Time {
    const u = 0

    fn native format(f)

    fn after(t) {
        this.u > t.u
    }

    fn before(t) {
        this.u < t.u
    }
}

// For now the time format is based on Go's native formatter. Eventually I would like
// to use the more standard formats using letters and numbers like yyyy-mm-dd.
export const Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
export const ANSIC       = "Mon Jan _2 15:04:05 2006"
export const UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
export const RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
export const RFC822      = "02 Jan 06 15:04 MST"
export const RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
export const RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
export const RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
export const RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
export const RFC3339     = "2006-01-02T15:04:05Z07:00"
export const Kitchen     = "3:04PM"
export const Stamp      = "Jan _2 15:04:05"
export const StampMilli = "Jan _2 15:04:05.000"
export const DateTime   = "2006-01-02 15:04:05"
export const DateOnly   = "2006-01-02"
export const TimeOnly   = "15:04:05"
