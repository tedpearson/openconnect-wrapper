package app

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Vpn struct {
	ui         *Ui
	scrollback int
	output     []string
	conn       Connection
}

type Connection struct {
	server string
	user   string
	search string
	group  string
}

func NewVpn(ui *Ui, scrollback int, conn Connection) *Vpn {
	return &Vpn{
		ui:         ui,
		scrollback: scrollback,
		output:     make([]string, 0, scrollback+1),
		conn:       conn,
	}
}

func (v *Vpn) Connect() error {
	v.log("Starting openconnect...")
	password, err := PasswordSearch(v.conn.search)
	if err != nil {
		return err
	}
	c := v.conn
	cmd := exec.Command("sudo", "/opt/homebrew/bin/openconnect", "--server", c.server, "-u", c.user, "--authgroup", c.group, "--passwd-on-stdin")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	_, err = stdin.Write(([]byte)(password))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	err = cmd.Start()
	if err != nil {
		v.error(err)
		return err
	}
	// fixme: how to catch panic in goroutine
	go v.consumeLoop(reader)
	return nil
}

func (v *Vpn) log(text string) {
	v.output = append(v.output, strings.TrimRight(text, "\n"))
	if len(v.output) > v.scrollback {
		v.output = v.output[len(v.output)-v.scrollback : len(v.output)]
	}
	// fixme: just add another row and remove from the top?
	v.ui.SetText(strings.Join(v.output, "\n"))
}

func (v *Vpn) error(err error) {
	v.log("Error: " + err.Error())
}

func (v *Vpn) consumeLoop(r *bufio.Reader) {
	for {
		text, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				// fixme: show reason
				v.log("openconnect exited.")
				return
			}
			panic(err)
		}
		print("output:" + text)
		v.log(text)
	}
}

func (v *Vpn) Disconnect() {
	cmd := exec.Command("sudo", "/usr/bin/pkill", "openconnect")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}

func (v *Vpn) waitForUI() {
	for {
		time.Sleep(1 * time.Second)
		if v.ui.grid != nil {
			return
		}
	}
}
