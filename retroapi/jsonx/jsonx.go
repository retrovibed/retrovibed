package jsonx

import (
	stdjson "encoding/json"
	"encoding/json/v2"
	"io"
)

func Codec() []json.Options {
	return []json.Options{
		json.DefaultOptionsV2(),
		stdjson.ReportErrorsWithLegacySemantics(true),
		json.WithUnmarshalers(json.JoinUnmarshalers(lenientInt64(), lenientUint64(), protoEnumUnmarshal())),
	}
}

func MarshalWrite(out io.Writer, in any) (err error) {
	return json.MarshalWrite(out, in, Codec()...)
}

func UnmarshalRead(in io.Reader, out any) (err error) {
	return json.UnmarshalRead(in, out, Codec()...)
}

func Marshal(in any) (_ []byte, err error) {
	return json.Marshal(in, Codec()...)
}

func Unmarshal(in []byte, out any) (err error) {
	return json.Unmarshal(in, out, Codec()...)
}
