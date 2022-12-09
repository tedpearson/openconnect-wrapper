package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

type Vpn struct {
	app        fyne.App
	ui         *Ui
	vpnChan    chan Command
	uiChan     chan State
	scrollback int
	output     []string
	menu       *Menu
}

type Connection struct {
	server    string
	user      string
	search    string
	group     string
	extraArgs string
}

func NewVpn(app fyne.App, ui *Ui, scrollback int, menu *Menu, uiChan chan State, vpnChan chan Command) *Vpn {
	vpn := &Vpn{
		app:        app,
		ui:         ui,
		scrollback: scrollback,
		output:     make([]string, 0, scrollback+1),
		menu:       menu,
		vpnChan:    vpnChan,
		uiChan:     uiChan,
	}
	go vpn.Listen()
	return vpn
}

func (v *Vpn) Connect(config Connection) {
	v.log("Starting openconnect...")
	v.menu.SetConnected(Connecting)
	password, err := PasswordSearch(config.search, v)
	if err != nil {
		v.log(err.Error())
		//log.Fatal(err)
	}
	cmd := exec.Command("sudo", "/opt/homebrew/bin/openconnect", "--server", config.server,
		"-u", config.user, "--authgroup", config.group, "--passwd-on-stdin", "--reconnect-timeout", "10")
	if config.extraArgs != "" {
		cmd.Args = append(cmd.Args, strings.Split(config.extraArgs, " ")...)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		v.error(err)
		//log.Fatal(err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		v.error(err)
		//log.Fatal(err)
	}
	_, err = stdin.Write(([]byte)(password + "\n"))
	if err != nil {
		v.error(err)
		//log.Fatal(err)
	}
	reader := bufio.NewReader(stdout)
	err = cmd.Start()
	if err != nil {
		v.error(err)
	}
	// fixme: how to catch panic in goroutine
	v.menu.SetConnected(Connected)
	go v.consumeLoop(reader)
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
				v.uiChan <- Disconnected
				v.menu.SetConnected(Disconnected)
				return
			}
			panic(err)
		} else {
			v.handleLine(text)
		}
	}
}

var (
	connectedPattern    = regexp.MustCompile("^Session authentication will expire")
	disconnectedPattern = regexp.MustCompile("Disconnected??")
)

func (v *Vpn) handleLine(text string) {
	print("output:" + text)
	v.log(text)
	if connectedPattern.MatchString(text) {
		v.uiChan <- Connected
	} else if disconnectedPattern.MatchString(text) {
		v.uiChan <- Disconnected
	}
}

func (v *Vpn) Disconnect() {
	cmd := exec.Command("sudo", "/usr/bin/pkill", "openconnect")
	if err := cmd.Run(); err != nil {
		v.error(err)
		//log.Fatal(err)
	}
	v.uiChan <- Disconnected
}

func (v *Vpn) waitForUI() {
	for {
		time.Sleep(1 * time.Second)
		if v.ui.grid != nil {
			return
		}
	}
}

func (v *Vpn) Listen() {
	for {
		select {
		case cmd := <-v.vpnChan:
			switch cmd {
			case Connect:
				config := ReadPreferences(v.app.Preferences())
				v.log(fmt.Sprintf("%+v", config))
				println("HELLO")
				v.Connect(config.Connection)
			case Disconnect:
				v.Disconnect()
			case DisconnectQuit:
				v.Disconnect()
				os.Exit(0)
			}
		}
	}
}
