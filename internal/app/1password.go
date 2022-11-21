package app

import (
	"io"
	"os/exec"
	"strings"
)

func PasswordSearch(search string) (string, error) {
	cmd := exec.Command("op", "item", "get", "--fields", "label=password", search)
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
	return strings.Trim(string(bytes), " \n\t"), nil
}
