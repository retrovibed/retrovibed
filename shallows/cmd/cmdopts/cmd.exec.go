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
	Program  string `kong:"-"`
	Override string `flag:"" name:"program" optional:"" help:"override the program to exec"`
	// Cwd defaults to ${vars_exec_cwd} (APPRUN_RUNTIME) so exec works inside an appimage, the bundled binaries
	// use a relative elf interpreter (lib64/ld-linux-x86-64.so.2) which only resolves from there.
	Cwd  string   `flag:"" name:"cwd" optional:"" default:"${vars_exec_cwd}" help:"directory to change into before exec"`
	Args []string `arg:"" optional:"" passthrough:""`
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

	// apprun's libc hooks normally chdir for exec, but syscall.Exec bypasses them.
	// the hooks of the exec'd program restore APPRUN_CWD.
	if err = os.Chdir(t.Cwd); err != nil {
		return err
	}

	return syscall.Exec(path, append([]string{program}, t.Args...), os.Environ())
}
