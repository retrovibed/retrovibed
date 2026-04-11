package duckdbx

import (
	"fmt"

	"github.com/grindlemire/go-lucene/pkg/driver"
	"github.com/grindlemire/go-lucene/pkg/lucene/expr"
)

type Lucene struct {
	driver.Base
}

func NewLucene() Lucene {
	fns := map[expr.Operator]driver.RenderFN{
		expr.Equals: func(left, right string) (string, error) {
			return fmt.Sprintf("%s ILIKE '%%' || %s || '%%'", left, right), nil
		},
	}

	// iterate over the existing base render functions and swap out any that you want to
	for op, sharedFN := range driver.Shared {
		_, found := fns[op]
		if !found {
			fns[op] = sharedFN
		}
	}

	// return the new driver ready to use
	return Lucene{
		driver.Base{
			RenderFNs: fns,
		},
	}
}
