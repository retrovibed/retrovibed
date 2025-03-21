package blockcache

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"

	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/iox"
)

func NewTorrentFromVirtualFS(v fsx.Virtual) *TorrentCacheStorage {
	return &TorrentCacheStorage{v: v}
}

var _ storage.ClientImpl = &TorrentCacheStorage{}

type TorrentCacheStorage struct {
	v fsx.Virtual
}

func (t *TorrentCacheStorage) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (storage.TorrentImpl, error) {
	path := filepath.Join("tmp", infoHash.HexString())
	if err := t.v.MkDirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	osio, err := t.v.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// initialize the file to all zeros
	n, err := io.Copy(osio, iox.ZeroReader{})
	if err != nil {
		return nil, err
	}

	if n < info.TotalLength() {
		return nil, io.ErrShortWrite
	}

	return &fileTorrentImpl{
		dir:        path,
		info:       info,
		infoHash:   infoHash,
		completion: storage.NewMapPieceCompletion(),
	}, nil
}

func (t *TorrentCacheStorage) Close() error {
	return nil
}

type fileTorrentImpl struct {
	dst        *os.File
	dir        string
	info       *metainfo.Info
	infoHash   metainfo.Hash
	completion storage.PieceCompletion
}

func (t *fileTorrentImpl) Close() error {
	return t.dst.Close()
}

func (t *fileTorrentImpl) Piece(p metainfo.Piece) storage.PieceImpl {
	return &nativePieceImpl{
		pieceManagementImpl: &piecemanagement{
			p:          p,
			completion: t.completion,
		},
		WriterAt: iox.NewSectionWriter(t.dst, p.Offset(), p.Length()),
		ReaderAt: io.NewSectionReader(t.dst, p.Offset(), p.Length()),
	}
}

type pieceManagementImpl interface {
	MarkComplete() error
	MarkNotComplete() error
	Completion() storage.Completion
}

type piecemanagement struct {
	p          metainfo.Piece
	completion storage.PieceCompletion
}

func (me *piecemanagement) pieceKey() metainfo.PieceKey {
	return metainfo.PieceKey{InfoHash: me.p.Hash(), Index: me.p.Index()}
}

func (fs *piecemanagement) Completion() storage.Completion {
	c, err := fs.completion.Get(fs.pieceKey())
	if err != nil {
		log.Printf("error getting piece completion: %s", err)
		c.Ok = false
		return c
	}

	return c
}

func (fs *piecemanagement) MarkComplete() error {
	return fs.completion.Set(fs.pieceKey(), true)
}

func (fs *piecemanagement) MarkNotComplete() error {
	return fs.completion.Set(fs.pieceKey(), false)
}

type nativePieceImpl struct {
	pieceManagementImpl
	io.ReaderAt
	io.WriterAt
}

var _ storage.PieceImpl = nativePieceImpl{}
