package app

import (
	"io"
	"os/exec"
	"strings"
)

func PasswordSearch(search string, v *Vpn) (string, error) {
	v.log("Getting password from 1Password...")
	cmd := exec.Command("/usr/local/bin/op", "item", "get", "--fields", "label=password", search)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err = cmd.Start(); err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	v.log("Password retrieved...")
	return strings.Trim(string(bytes), " \n\t"), nil
}
