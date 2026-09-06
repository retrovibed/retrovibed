package publishplugin

import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/internal/md5x"
)

// Identity is the catalog id for the module installed at path: the md5 of its
// contents, folded together with the filename when the filename is something
// other than the content digest itself.
//
// The id is the plugin. Identical bytes under the same identity always hash to
// the same id, so a reinstall, a re-import of the same module from a second
// community, and a duplicate install all land on the row that already exists -
// whatever a community selected stays selected. Different bytes are a different
// plugin and get a row of their own.
//
// The canonical {digest}.wasm form - what the upload endpoint and the community
// importer both write - hashes to exactly its contents, so every install path
// agrees on the id for the same module. Anything installed under a name of its
// own (the CLI, a hand copy) folds that name in, which is what gives each
// symlink to a shared module an identity of its own, matching how EnvPath keys
// configuration by path rather than by contents.
//
// The upload endpoint depends on the canonical case: it digests the upload as
// it streams and writes the result to {id}.wasm, so Identity of what it wrote
// is the id it recorded.
func Identity(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errorsx.Wrapf(err, "unable to open plugin: %s", path)
	}
	defer f.Close()

	digest := md5.New()
	if _, err := io.Copy(digest, f); err != nil {
		return "", errorsx.Wrapf(err, "unable to digest plugin: %s", path)
	}

	content := md5x.FormatUUID(digest)

	stem := strings.TrimSuffix(filepath.Base(path), ".wasm")
	if stem == content {
		return content, nil
	}

	return md5x.FormatUUID(md5x.Digest(content, stem)), nil
}
