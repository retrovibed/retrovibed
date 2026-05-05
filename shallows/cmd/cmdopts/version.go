package cmdopts

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"
)

func BuildVersion() (_ string, err error) {
	var (
		ok   bool
		info *debug.BuildInfo
		ts   time.Time
		id   string
		// dirty bool
	)

	if info, ok = debug.ReadBuildInfo(); ok {
		for _, v := range info.Settings {
			switch v.Key {
			case "vcs.modified":
				// if dirty, err = strconv.ParseBool(v.Value); err != nil {
				// 	return "", err
				// }
			case "vcs.revision":
				id = v.Value
			case "vcs.time":
				if ts, err = time.Parse(time.RFC3339, v.Value); err != nil {
					return "", err
				}
			}
		}

		return fmt.Sprintf("%s %s %s", info.Main.Path, ts.Format("2006-01-02"), id), nil
	}

	return "", errors.New("unable to detect build information")
}
