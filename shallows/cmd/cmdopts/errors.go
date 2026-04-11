package cmdopts

import "log"

func Fatal(err error) {
	if err == nil {
		return
	}
	log.Fatalln(err)
}
