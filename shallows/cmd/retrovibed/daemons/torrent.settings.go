package daemons

import "google.golang.org/protobuf/encoding/protojson"

func (t *TorrentSettings) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *TorrentSettings) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}
