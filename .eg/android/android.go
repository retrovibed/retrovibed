package android

import (
	"context"
	"eg/compute/envx"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_eg "github.com/egdaemon/eg"
	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egfs"
	"github.com/egdaemon/eg/runtime/x/wasi/egsecrets"
	"github.com/egdaemon/wasinet/wasinet"
	"golang.org/x/net/http2"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

const (
	ID        = "space.retrovibe.retrovibed"
	NDKPath   = "/opt/android-sdk/ndk/27.0.12077973"
	Platform  = "android-31"
	Container = "retrovibe.ubuntu.android"
)

// SigningKey generates an Android upload signing key using keytool. The resulting
// keystore is used to sign APKs/AABs before uploading to the Google Play Store.
//
// dstpath is the absolute path where the .keystore file will be written; the
// parent directory is used as the working directory for keytool.
//
// alias is the name used to identify the key entry inside the keystore. The same
// alias must be provided when signing artifacts with jarsigner or apksigner.
//
//	android.SigningKey(egenv.CacheDirectory("release.keystore"), "upload")
func SigningKey(dstpath, alias string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		if _, err := os.Stat(dstpath); err == nil {
			return nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(dstpath), 0755); err != nil {
			return err
		}

		env := shell.Env().
			Environ("DESTINATION", dstpath).
			Environ("ALIAS", alias).
			Directory(filepath.Dir(dstpath))
		return shell.Run(
			ctx,
			env.New("keytool -genkey -v -keystore ${DESTINATION} -alias ${ALIAS} -keyalg RSA -keysize 2048 -validity 10000"),
		)
	}
}

func Keystore(secreturi, dstpath string) eg.OpFn {
	return eg.WhenFn(egfs.FileNotExistsFn(dstpath), egsecrets.CopyIntoFileOp(dstpath, secreturi))
}

func Draft(id, bundlePath, trackName string) eg.OpFn {
	return release(id, bundlePath, trackName, "draft")
}

func Deploy(id, bundlePath, trackName string) eg.OpFn {
	return release(id, bundlePath, trackName, "completed")
}

func release(id, bundlePath, trackName, status string) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) (err error) {
		gcpcreds := option.WithCredentialsJSON(envx.Base64([]byte("{}"), _eg.EnvUnsafeGcloudADCB64))
		scopes := option.WithScopes(androidpublisher.AndroidpublisherScope)

		bt := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return wasinet.DialContext(ctx, network, addr)
			},
		}

		// Configure HTTP/2 on the base transport to control ReadIdleTimeout.
		// cloud.google.com/go/auth/httptransport.defaultBaseTransport sets
		// ReadIdleTimeout=31s, but that function is only called when no
		// BaseRoundTripper is provided. By supplying bt directly, we bypass it.
		// http2.ConfigureTransports is required to get a handle to the embedded
		// *http2.Transport so ReadIdleTimeout can be set at all.
		h2t, err := http2.ConfigureTransports(bt)
		if err != nil {
			return err
		}
		h2t.ReadIdleTimeout = 15 * time.Second

		rt, err := htransport.NewTransport(ctx, bt, gcpcreds, scopes)
		if err != nil {
			return err
		}

		service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(&http.Client{Transport: rt}))
		if err != nil {
			return fmt.Errorf("play store init failed: %w", err)
		}

		edit, err := service.Edits.Insert(id, &androidpublisher.AppEdit{}).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed to create edit: %w", err)
		}

		defer func() {
			if err == nil {
				return
			}

			log.Println(fmt.Errorf("deploy failed: %w", err))
			failed := service.Edits.Delete(id, edit.Id).Context(ctx).Do()
			if failed != nil {
				log.Println(fmt.Errorf("cleanup failed: %w", failed))
			}
		}()

		file, err := os.Open(bundlePath)
		if err != nil {
			return err
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			return err
		}
		reader := &progressReader{
			r:     file,
			total: fi.Size(),
		}

		bundle, err := service.Edits.Bundles.Upload(id, edit.Id).Media(reader, googleapi.ContentType("application/octet-stream")).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("bundle upload failed: %w", err)
		}

		track := &androidpublisher.Track{
			Track: trackName,
			Releases: []*androidpublisher.TrackRelease{
				{
					VersionCodes: []int64{bundle.VersionCode},
					Status:       status,
				},
			},
		}

		_, err = service.Edits.Tracks.Update(id, edit.Id, trackName, track).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("track update failed: %w", err)
		}

		_, err = service.Edits.Commit(id, edit.Id).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("commit failed: %w", err)
		}

		return nil
	}
}

type progressReader struct {
	r       io.Reader
	total   int64
	current int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.current += int64(n)

	// Log every 5MB (5 * 1024 * 1024 bytes)
	if n > 0 && pr.current%(5*1024*1024) < int64(n) {
		log.Printf("Upload Progress: %.2f MB / %.2f MB",
			float64(pr.current)/1024/1024,
			float64(pr.total)/1024/1024)
	}
	return n, err
}
