package media

import (
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type AutoimportDirectoryOption func(*AutoimportDirectory)

func AutoimportDirectoryOptionFromLibrary(cc library.AutoimportDirectory) AutoimportDirectoryOption {
	return func(c *AutoimportDirectory) {
		c.Id = cc.ID
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.Path = cc.Path
		c.Description = cc.Description
		c.Debounce = uint64(cc.Debounce / time.Second)
		c.Mode = cc.Mode
		c.LibraryDirectoryId = cc.LibraryDirectoryID
		c.LastScannedAt = grpcx.EncodeTime(cc.LastScannedAt)
	}
}

func AutoimportDirectoryOptionCounts(pending, imported uint64) AutoimportDirectoryOption {
	return func(c *AutoimportDirectory) {
		c.Pending = pending
		c.Imported = imported
	}
}

func NewLibraryAutoimportDirectory(cc *AutoimportDirectory) library.AutoimportDirectory {
	return library.AutoimportDirectory{
		ID:                 cc.Id,
		Path:               cc.Path,
		Description:        cc.Description,
		Debounce:           time.Duration(cc.Debounce) * time.Second,
		Mode:               cc.Mode,
		LibraryDirectoryID: cc.LibraryDirectoryId,
	}
}
