package torrent

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"
	"sync/atomic"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/pkg/errors"
)

func newDigestsFromTorrent(t *torrent) digests {
	return newDigests(
		t.storage,
		t.piece,
		func(idx int, cause error) {
			if t.chunks == nil {
				panic("gorp")
			}

			// log.Printf("hashed %p %d / %d - %v", t.chunks, idx+1, t.numPieces(), cause)
			t.chunks.Hashed(idx, cause)

			t.event.Broadcast()
			t.cln.event.Broadcast() // cause the client to detect completed torrents.
		},
	)
}

func newDigests(iora io.ReaderAt, retrieve func(int) *metainfo.Piece, complete func(int, error)) digests {
	if iora == nil {
		panic("digests require a storage implementation")
	}

	// log.Printf("new digest %T\n", iora)
	return digests{
		ReaderAt: iora,
		retrieve: retrieve,
		complete: complete,
		pending:  newBitQueue(),
	}
}

// digests is responsible correctness of received data.
type digests struct {
	ReaderAt io.ReaderAt
	retrieve func(int) *metainfo.Piece
	complete func(int, error)
	// marks whether digest is actively processing.
	reaping int64
	// cache of the pieces that need to be verified.
	pending *bitQueue
}

// Enqueue a piece to check its completed digest.
func (t *digests) Enqueue(idx int) {
	t.pending.Push(idx)
	t.verify()
}

func (t *digests) verify() {
	if atomic.AddInt64(&t.reaping, 1) > int64(runtime.NumCPU()) {
		atomic.AddInt64(&t.reaping, -1)
		return
	}

	go func() {
		for idx, ok := t.pending.Pop(); ok; idx, ok = t.pending.Pop() {
			t.check(idx)
		}

		atomic.AddInt64(&t.reaping, -1)
	}()
}

func (t *digests) check(idx int) {
	var (
		err    error
		digest metainfo.Hash
		p      *metainfo.Piece
	)

	if p = t.retrieve(idx); p == nil {
		t.complete(idx, fmt.Errorf("piece %d not found during digest", idx))
		return
	}

	if digest, err = t.compute(p); err != nil {
		t.complete(idx, err)
		return
	}

	if digest != p.Hash() {
		t.complete(idx, fmt.Errorf("piece %d digest mismatch %s != %s", idx, hex.EncodeToString(digest[:]), p.Hash().HexString()))
		return
	}

	t.complete(idx, nil)
}

func (t *digests) compute(p *metainfo.Piece) (ret metainfo.Hash, err error) {
	c := sha1.New()
	plen := p.Length()

	n, err := io.Copy(c, io.NewSectionReader(t.ReaderAt, p.Offset(), plen))
	if err != nil {
		return ret, errors.Wrapf(err, "piece %d digest failed", p.Offset())
	}

	if n != plen {
		return ret, fmt.Errorf("piece digest failed short copy %d: %d != %d", p.Offset(), n, plen)
	}

	copy(ret[:], c.Sum(nil))

	return ret, nil
}
