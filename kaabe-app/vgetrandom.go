//go:build linux
// +build linux

package unix

import (
    "golang.org/x/sys/unix"
)

// VGetRandom reads random bytes into p using the getrandom syscall.
func VGetRandom(p []byte, flags uint32) (int, bool) {
    n, err := unix.Getrandom(p, int(flags))
    return n, err == nil
}
