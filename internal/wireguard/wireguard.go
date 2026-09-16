// Package wireguard manages the WireGuard tunnel via wg-quick.
package wireguard

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// wgInterfaceName is the fixed WireGuard interface name this tool manages.
// wg-quick derives the interface name from the config file's basename.
const wgInterfaceName = "eduvpn"

// WgConfigPath returns the fixed path of the WireGuard config file inside
// the state directory.
func WgConfigPath(dir string) string {
	return filepath.Join(dir, wgInterfaceName+".conf")
}

// WriteWGConfig writes the WireGuard config contents to path with 0600
// permissions, since it embeds a private key.
func WriteWGConfig(path, contents string) error {
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		return err
	}
	// os.WriteFile does not chmod an existing file on O_TRUNC, so make sure
	// permissions are correct even if the file pre-existed with looser ones.
	return os.Chmod(path, 0o600)
}

// runSudo runs a command under sudo with inherited stdio so the user sees
// the interactive password prompt.
func runSudo(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// runSudoOutput is like runSudo but captures stdout instead of streaming it,
// so the caller can inspect the command's output.
func runSudoOutput(args ...string) (string, error) {
	cmd := exec.Command("sudo", args...)
	cmd.Stdin, cmd.Stderr = os.Stdin, os.Stderr
	out, err := cmd.Output()
	return string(out), err
}

// WgQuickUp brings up the WireGuard interface defined by the config at path.
func WgQuickUp(path string) error {
	return runSudo("wg-quick", "up", path)
}

// WgQuickDown tears down the WireGuard interface defined by the config at path.
func WgQuickDown(path string) error {
	return runSudo("wg-quick", "down", path)
}

// realInterface resolves the OS-assigned utun device backing
// wgInterfaceName. macOS can't create a kernel interface with an arbitrary
// name the way Linux can, so wg-quick creates a utunN device and records
// the mapping in /var/run/wireguard/<name>.name; `wg` itself has no notion
// of the friendly name, so this file must be read to find the real device.
func realInterface() (string, error) {
	out, err := runSudoOutput("cat", "/var/run/wireguard/"+wgInterfaceName+".name")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// WgShow returns the `wg show` output for the interface this tool manages.
// If the interface is not up, it returns ("", nil) rather than an error,
// since a leftover config file with no active interface (e.g. after a
// reboot, or a manual `wg-quick down`) is a normal disconnected state, not
// a failure. Reading WireGuard's control socket requires root, so this
// prompts for a password like WgQuickUp/WgQuickDown do.
func WgShow() (string, error) {
	iface, err := realInterface()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", nil
		}
		return "", err
	}

	out, err := runSudoOutput("wg", "show", iface)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", nil
		}
		return "", err
	}
	return out, nil
}
