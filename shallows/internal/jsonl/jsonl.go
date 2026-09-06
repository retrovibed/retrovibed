package jsonl

import (
	"bufio"
	"context"
	"io"
	"iter"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

type Encoder struct {
	out io.Writer // The stream each encoded value is written to
}

// NewEncoder returns a new Encoder that writes to w.
// jsonx.MarshalWrite does not terminate a value with a newline, so Encode
// appends the delimiter itself to make the output suitable for JSONL.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		out: w,
	}
}

// Encode writes the JSON encoding of each v to the stream, each followed by a newline.
// Any error during JSON marshaling or writing to the underlying writer is returned.
func (e *Encoder) Encode(vs ...any) error {
	for _, v := range vs {
		if err := jsonx.MarshalWrite(e.out, v); err != nil {
			return err
		}

		if _, err := io.WriteString(e.out, "\n"); err != nil {
			return err
		}
	}
	return nil
}

// Decoder reads JSON lines from an underlying io.Reader.
// Each call to Decode reads one line, decodes it as a JSON object,
// and stores it in the provided interface.
type Decoder struct {
	scanner *bufio.Scanner // Used to read the input stream line by line
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return NewDecoderWithMax(16*bytesx.MiB, r)
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoderWithMax(max int, r io.Reader) *Decoder {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, bytesx.MiB), langx.FirstNonZero(max, 16*bytesx.MiB))
	return &Decoder{
		scanner: s,
	}
}

// Iter wraps a Decoder as an iterx.Seq[T], decoding one value per iteration.
// Iteration stops at EOF or on error; call Err() after ranging to check for errors.
func Iter[T any](d *Decoder) *DecoderSeq[T] {
	return &DecoderSeq[T]{d: d}
}

// DecoderSeq implements iterx.Seq[T] for a jsonl.Decoder.
type DecoderSeq[T any] struct {
	d   *Decoder
	err error
}

func (s *DecoderSeq[T]) Each(_ context.Context) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			var v T
			if err := s.d.Decode(&v); err == io.EOF {
				return
			} else if err != nil {
				s.err = err
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

func (s *DecoderSeq[T]) Err() error {
	return s.err
}

// Decode reads the next JSON line from the stream and stores the result
// in the value pointed to by v.
//
// It returns io.EOF if no more lines are available in the stream.
// It returns an error if scanning a line fails or if the line cannot be
// decoded as a JSON object into v.
func (d *Decoder) Decode(v any) error {
	if !d.scanner.Scan() {
		if err := d.scanner.Err(); err != nil {
			return errorsx.Wrap(err, "failed to scan a line")
		}
		return io.EOF
	}

	line := d.scanner.Bytes()

	// If the line is empty, it's not a valid JSON object.
	// You might choose to skip empty lines or return an error based on strictness.
	// Here, we let json.NewDecoder handle it, which will likely return an error.
	if len(line) == 0 {
		return errorsx.Errorf("encountered empty line, expected JSON object")
	}

	if err := jsonx.Unmarshal(line, v); err != nil {
		return errorsx.Wrapf(err, "unable to decode '%s'", string(line))
	}

	return nil
}
