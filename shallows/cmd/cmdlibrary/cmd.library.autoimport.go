package cmdlibrary

import (
	"fmt"
	"io"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
)

// autoimport command examples
// go -C shallows run ./cmd/retrovibe/... library autoimport ls --insecure --endpoint="https://localhost:9998"
// go -C shallows run ./cmd/retrovibe/... library autoimport create /tmp/inbox --debounce=1m --mode=move --insecure --endpoint="https://localhost:9998"
// go -C shallows run ./cmd/retrovibe/... library autoimport show <id> --insecure --endpoint="https://localhost:9998"
// go -C shallows run ./cmd/retrovibe/... library autoimport edit <id> --mode=copy --insecure --endpoint="https://localhost:9998"
// go -C shallows run ./cmd/retrovibe/... library autoimport rm <id> --insecure --endpoint="https://localhost:9998"
type cmdAutoimport struct {
	Ls     cmdAutoimportLs     `cmd:"" help:"list the directories monitored for automatic import"`
	Show   cmdAutoimportShow   `cmd:"" help:"show a monitored directory along with its pending and imported file counts"`
	Create cmdAutoimportCreate `cmd:"" help:"monitor a directory, importing files into the library once they stop being modified"`
	Rm     cmdAutoimportRm     `cmd:"" help:"stop monitoring a directory, previously imported media remains in the library"`
	Edit   cmdAutoimportEdit   `cmd:"" help:"edit a monitored directory (description, debounce, mode, library directory)"`
}

const (
	autoimportModeCopy = "copy"
	autoimportModeMove = "move"
)

func autoimportModeFromString(s string) (uint32, error) {
	switch s {
	case autoimportModeCopy:
		return library.AutoimportModeCopy, nil
	case autoimportModeMove:
		return library.AutoimportModeMove, nil
	default:
		return 0, errorsx.Errorf("unknown mode: %s", s)
	}
}

func autoimportModeString(m uint32) string {
	switch m {
	case library.AutoimportModeCopy:
		return autoimportModeCopy
	case library.AutoimportModeMove:
		return autoimportModeMove
	default:
		return fmt.Sprintf("unknown(%d)", m)
	}
}

func printAutoimportDirectory(w io.Writer, d *media.AutoimportDirectory) error {
	_, err := fmt.Fprintf(w, "id='%s' path='%s' mode=%s debounce=%s library_directory='%s' description='%s' last_scanned='%s'\n", d.Id, d.Path, autoimportModeString(d.Mode), time.Duration(d.Debounce)*time.Second, d.LibraryDirectoryId, d.Description, d.LastScannedAt)
	return err
}
