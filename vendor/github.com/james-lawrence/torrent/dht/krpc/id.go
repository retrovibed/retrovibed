package krpc

import (
	"encoding"
	"encoding/hex"
	"fmt"

	"github.com/james-lawrence/torrent/bencode"
	"github.com/james-lawrence/torrent/dht/int160"
)

func RandomID() (id ID) {
	return ID(int160.Random().Bytes())
}

type ID [20]byte

var (
	_ interface {
		bencode.Unmarshaler
		encoding.TextUnmarshaler
	} = (*ID)(nil)
	_ bencode.Marshaler = ID{}
	_ fmt.Formatter     = ID{}
)

func (h ID) Format(f fmt.State, c rune) {
	// See metainfo.Hash.
	f.Write([]byte(h.String()))
}

func IdFromString(s string) (id ID) {
	if n := copy(id[:], s); n != 20 {
		panic(n)
	}
	return
}

func (id ID) MarshalBencode() ([]byte, error) {
	return []byte("20:" + string(id[:])), nil
}

func (id *ID) UnmarshalBencode(b []byte) error {
	var s string
	if err := bencode.Unmarshal(b, &s); err != nil {
		return err
	}
	if n := copy(id[:], s); n != 20 {
		return fmt.Errorf("string has wrong length: %d", n)
	}
	return nil
}

func (id *ID) UnmarshalText(b []byte) (err error) {
	n, err := hex.Decode(id[:], b)
	if err != nil {
		return
	}
	if n != len(*id) {
		err = fmt.Errorf("expected %v bytes, only got %v", len(*id), n)
	}
	return
}

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

func (id ID) Int160() int160.T {
	return int160.FromByteArray(id)
}

func (id ID) IsZero() bool {
	return id == [20]byte{}
}
