package bytesx

import (
	"fmt"
	"strconv"
	"strings"
)

type Unit int64

func (t Unit) Format(f fmt.State, verb rune) {
	div := int64(1)
	suffix := ""
	switch {
	case t > EiB:
		div = EiB
		suffix = "e"
	case t > PiB:
		div = PiB
		suffix = "p"
	case t > TiB:
		div = TiB
		suffix = "t"
	case t > GiB:
		div = GiB
		suffix = "g"
	case t > MiB:
		div = MiB
		suffix = "m"
	case t > KiB:
		div = KiB
		suffix = "k"
	}

	fmt.Fprintf(f, "%d%s", uint64(float64(t)/float64(div)), suffix)
}

// base 2 byte units
const (
	_   Unit = iota
	KiB      = 1 << (10 * iota)
	MiB
	GiB
	TiB
	PiB
	EiB
)

// unitSuffixes maps the size suffixes rendered by scraped listings (e.g.
// "214.2 MB", "1.2 GB") onto their base-2 byte multiplier.
var unitSuffixes = map[string]Unit{
	"B":  1,
	"KB": KiB,
	"MB": MiB,
	"GB": GiB,
	"TB": TiB,
	"PB": PiB,
}

// Parse converts a rendered size string like "214.2 MB" into bytes,
// returning 0 if s doesn't match the expected "<number> <unit>" shape.
func Parse(s string) uint64 {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	mult, ok := unitSuffixes[strings.ToUpper(parts[1])]
	if !ok {
		return 0
	}

	return uint64(value * float64(mult))
}
