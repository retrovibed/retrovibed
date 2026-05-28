package cmdopts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
)

// CmdExec is a reusable kong command that execs a pre-configured program,
// replacing the current process. All CLI tokens following the command name
// are passed to the target program unchanged (flags included).
//
// Register in a parent struct like:
//
//	Play CmdExec `cmd:"" help:"play media with mpv"`
//
// and construct with:
//
//	Play: cmdopts.NewCmdExec("mpv"),
type CmdExec struct {
	Program string   `flag:"" name:"program" optional:"" help:"override the program to exec"`
	Args    []string `arg:"" optional:"" passthrough:""`
}

// NewCmdExec returns a CmdExec that will exec the named program.
func NewCmdExec(program string) CmdExec {
	return CmdExec{Program: program}
}

func (t CmdExec) Run() (err error) {
	relative := func() string {
		self := errorsx.Zero(os.Executable())
		if self == "" {
			return ""
		}
		candidate := filepath.Join(filepath.Dir(self), t.Program)
		if _, err := os.Stat(candidate); err != nil {
			return ""
		}
		return candidate
	}()

	global := errorsx.Zero(exec.LookPath(t.Program))
	path := langx.FirstNonZero(relative, global)
	if path == "" {
		return fmt.Errorf("unable to locate program: %s", t.Program)
	}

	return syscall.Exec(path, append([]string{t.Program}, t.Args...), os.Environ())
}
