package tracking

import (
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
)

func TestDuckdbNowPlusInterval(t *testing.T) {
	q := sqltestx.Metadatabase(t)
	row := q.QueryRowContext(t.Context(), "SELECT NOW() + to_seconds(30)")
	var s string
	if err := row.Scan(&s); err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("result: %s", s)
}
