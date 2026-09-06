package jsonx

import (
	stdjson "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
)

func Codec() []json.Options {
	return []json.Options{
		json.DefaultOptionsV2(),
		stdjson.ReportErrorsWithLegacySemantics(true),
		// v2 narrowed omitempty to "encodes to an empty JSON value", so a zero
		// number or false no longer drops out. protobuf generated structs tag
		// every field omitempty, so without v1's meaning every payload gains
		// its zero valued fields back.
		stdjson.OmitEmptyWithLegacySemantics(true),
		// v2 matches object members to struct fields case sensitively, v1 did
		// not. keep v1's behaviour: this package is a drop in replacement for
		// encoding/json, and without it any payload whose casing differs from
		// the field name (or its tag) silently decodes to the zero value.
		json.MatchCaseInsensitiveNames(true),
		json.WithUnmarshalers(json.JoinUnmarshalers(lenientInt64(), lenientUint64(), protoEnumUnmarshal())),
	}
}

func MarshalWrite(out io.Writer, in any) (err error) {
	return json.MarshalWrite(out, in, Codec()...)
}

// UnmarshalRead decodes a single value from in. It goes through
// UnmarshalDecode, not UnmarshalRead: UnmarshalRead is Unmarshal over a
// reader, so it reports an empty stream as a syntax error and rejects
// anything trailing the value. UnmarshalDecode is the analogue of
// encoding/json's Decoder.Decode, which is what callers of this replaced -
// an empty stream comes back as io.EOF, which they branch on to treat an
// absent body as absent rather than malformed.
func UnmarshalRead(in io.Reader, out any) (err error) {
	return json.UnmarshalDecode(jsontext.NewDecoder(in), out, Codec()...)
}

func Marshal(in any) (_ []byte, err error) {
	return json.Marshal(in, Codec()...)
}

func Unmarshal(in []byte, out any) (err error) {
	return json.Unmarshal(in, out, Codec()...)
}
