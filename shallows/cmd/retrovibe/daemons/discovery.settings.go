package daemons

import "google.golang.org/protobuf/encoding/protojson"

func (t *DiscoverySettings) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *DiscoverySettings) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}
