package mediaapi

import "google.golang.org/protobuf/encoding/protojson"

// MarshalJSON ensures encoding/json (and anything that defers to it, such as
// json.Marshal on a struct embedding Stream) produces protojson output
// instead of reflecting over the generated struct, which would mishandle the
// Command oneof and internal protoimpl state.
func (x *Stream) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

// UnmarshalJSON is the counterpart to MarshalJSON, routing decoding through
// protojson for the same reason.
func (x *Stream) UnmarshalJSON(data []byte) error {
	return protojson.Unmarshal(data, x)
}
