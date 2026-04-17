package authn

import (
	"google.golang.org/protobuf/encoding/protojson"
)

func (t *Token) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *Token) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}
