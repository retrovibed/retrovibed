package grpcx

import (
	"encoding/json"
	"time"

	"google.golang.org/protobuf/proto"
)

func DecodeTime(formatted string) (p time.Time, err error) {
	return time.Parse(time.RFC3339Nano, formatted)
}

func EncodeTime(ts time.Time) string {
	return ts.Format(time.RFC3339Nano)
}

func JSONEncode[X proto.Message, Y any](from X, to *Y) (err error) {
	var (
		encoded []byte
	)

	if encoded, err = json.Marshal(from); err != nil {
		return err
	}

	return json.Unmarshal(encoded, to)
}

func JSONDecode[X proto.Message, Y any](from Y, to X) (err error) {
	var (
		encoded []byte
	)

	if encoded, err = json.Marshal(from); err != nil {
		return err
	}

	return json.Unmarshal(encoded, to)
}
