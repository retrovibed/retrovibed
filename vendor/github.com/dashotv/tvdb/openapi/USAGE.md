<!-- Start SDK Example Usage [usage] -->
```go
package main

import (
	"context"
	"github.com/dashotv/tvdb/openapi"
	"github.com/dashotv/tvdb/openapi/models/shared"
	"log"
)

func main() {
	s := openapi.New(
		openapi.WithSecurity(shared.Security{
			BearerAuth: "<YOUR_BEARER_TOKEN_HERE>",
		}),
	)

	ctx := context.Background()
	res, err := s.ArtworkStatuses.GetAllArtworkStatuses(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.Object != nil {
		// handle response
	}
}

```
<!-- End SDK Example Usage [usage] -->