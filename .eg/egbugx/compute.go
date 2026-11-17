// Package egbugx extends egbug with debugging operations specific to this repository.
package egbugx

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

// Compute prints the compute resources available to the workload: CPU, RAM, disk and VRAM.
// cgroup limits are printed alongside the host totals because the limits are what the workload is
// actually held to, and tmpfs mounts are counted against the memory limit. Every command is
// lenient, tools or files that do not exist on the machine (e.g. no gpu) are reported and skipped.
// it has the shape of a debug operation, see egbug.DebugFailure.
func Compute(ctx context.Context, op eg.Op) error {
	debug := shell.Runtime().Lenient(true)
	return shell.Run(
		ctx,
		debug.New("echo '=== cpu'"),
		debug.New("grep -m1 'model name' /proc/cpuinfo; echo \"available: $(nproc) of $(nproc --all)\""),
		debug.New("cat /sys/fs/cgroup/cpu.max 2>/dev/null || echo 'no cgroup cpu limit'"),

		debug.New("echo '=== ram'"),
		debug.New("free -m"),
		debug.New("echo \"cgroup limit: $(cat /sys/fs/cgroup/memory.max 2>/dev/null || echo none) used: $(cat /sys/fs/cgroup/memory.current 2>/dev/null || echo unknown)\""),

		debug.New("echo '=== disk'"),
		debug.New("df -h"),
		debug.Newf("df -ih /tmp %s %s", egenv.DirectoryWorkload(), egenv.DirectoryCache()),

		debug.New("echo '=== vram'"),
		debug.New("nvidia-smi --query-gpu=name,memory.total,memory.used,memory.free --format=csv 2>/dev/null || echo 'no gpu'"),
		debug.New("for d in /sys/class/drm/card*/device; do [ -r $d/mem_info_vram_total ] && echo \"$d total: $(cat $d/mem_info_vram_total) used: $(cat $d/mem_info_vram_used)\"; done; true"),
	)
}
