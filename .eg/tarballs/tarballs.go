package tarballs

import "github.com/egdaemon/eg/runtime/x/wasi/egtarball"

func Retrovibed() string {
	return egtarball.GitPattern("retrovibed")
}
