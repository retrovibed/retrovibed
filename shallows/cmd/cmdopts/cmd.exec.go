package cmdopts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
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
	// Program is set by NewCmdExec. it is hidden from kong because kong resets
	// parsed fields to their zero value, which would discard it.
	Program  string   `kong:"-"`
	Override string   `flag:"" name:"program" optional:"" help:"override the program to exec"`
	Args     []string `arg:"" optional:"" passthrough:""`
}

// NewCmdExec returns a CmdExec that will exec the named program.
func NewCmdExec(program string) CmdExec {
	return CmdExec{Program: program}
}

func (t CmdExec) Run() (err error) {
	program := langx.FirstNonZero(t.Override, t.Program)
	if program == "" {
		return fmt.Errorf("no program configured to exec")
	}

	relative := func() string {
		self := errorsx.Zero(os.Executable())
		if self == "" {
			return ""
		}
		candidate := filepath.Join(filepath.Dir(self), program)
		if _, err := os.Stat(candidate); err != nil {
			return ""
		}
		return candidate
	}()

	global := errorsx.Zero(exec.LookPath(program))
	path := langx.FirstNonZero(relative, global)
	if path == "" {
		return fmt.Errorf("unable to locate program: %s", program)
	}

	return syscall.Exec(path, append([]string{program}, t.Args...), os.Environ())
}
