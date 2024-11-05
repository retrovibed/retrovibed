package cmdopts

import "log"

func ReportError(err error) error {
	if err != nil {
		log.Println(err)
	}
	return err
}
