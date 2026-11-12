package authn

import (
	"fmt"
	"io"
	"strings"
)

// PrintPermissions writes the permissions granted by the token as a single line.
func PrintPermissions(w io.Writer, t *Token) error {
	var granted []string
	flag := func(name string, ok bool) {
		if ok {
			granted = append(granted, name)
		}
	}

	flag("usermanagement", t.Usermanagement)
	flag("billing_read", t.BillingRead)
	flag("billing_modify", t.BillingModify)
	flag("community_modify", t.CommunityModify)
	flag("community_sync", t.CommunitySync)
	flag("library_read", t.LibraryRead)
	flag("library_modify", t.LibraryModify)
	flag("archive_sync", t.ArchiveSync)
	flag("device_backup", t.DeviceBackup)
	flag("remote_control", t.RemoteControl)
	flag("local_only", t.LocalOnly)
	flag(fmt.Sprintf("archive_upload=%d", t.ArchiveUpload), t.ArchiveUpload > 0)
	flag(fmt.Sprintf("archive_download=%d", t.ArchiveDownload), t.ArchiveDownload > 0)

	_, err := fmt.Fprintln(w, "permissions", strings.Join(granted, ","))
	return err
}
