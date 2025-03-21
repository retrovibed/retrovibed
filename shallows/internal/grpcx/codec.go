package grpcx

import "time"

func DecodeTime(formatted string) (p time.Time, err error) {
	return time.Parse(time.RFC3339Nano, formatted)
}

func EncodeTime(ts time.Time) string {
	return ts.Format(time.RFC3339Nano)
}
