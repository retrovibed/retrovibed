package meta

import "google.golang.org/protobuf/encoding/protojson"

func (t *PublishContentRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *PublishContentRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}
