package retrovibedbindx

import "github.com/retrovibed/retrovibed/shallows/mediaapi"

// RemoteControlListenToken mints a bearer authorized to connect to this
// process's local /rc/listen endpoint. Never valid outside this process.
func RemoteControlListenToken() (string, error) {
	return mediaapi.RemoteControlListenToken()
}
