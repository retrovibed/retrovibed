package media

import "google.golang.org/protobuf/encoding/protojson"

func (t *RecentRecordRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *RecentRecordRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
}
