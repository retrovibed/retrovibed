package retrovibed

import (
	"database/sql/driver"
	"fmt"
)

// Float128 bridges duckdb FLOAT[128] columns to []float32. The duckdb driver
// returns ARRAY values as []any with each element typed float32; Scan converts
// that representation. Value passes []float32 straight through — the driver
// accepts it natively as a FLOAT[] bind parameter.
type Float128 []float32

// Scan implements the database/sql Scanner interface.
func (dst *Float128) Scan(src any) error {
	if src == nil {
		*dst = nil
		return nil
	}

	switch v := src.(type) {
	case []float32:
		out := make(Float128, len(v))
		copy(out, v)
		*dst = out
		return nil
	case []any:
		out := make(Float128, len(v))
		for i, el := range v {
			f, ok := el.(float32)
			if !ok {
				return fmt.Errorf("retrovibed: Float128 scan: element %d not float32 (got %T)", i, el)
			}
			out[i] = f
		}
		*dst = out
		return nil
	}

	return fmt.Errorf("retrovibed: Float128 scan: cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Float128) Value() (driver.Value, error) {
	if src == nil {
		return nil, nil
	}
	return []float32(src), nil
}
