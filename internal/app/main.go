package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Config struct {
	Scrollback int
	Connection Connection
}

// todo:
//   handle untriggered disconnects and update ui state. use mvc type pattern?

func Main() error {
	a := app.NewWithID("com.tedpearson.openconnectwrapper")
	config := ReadPreferences(a.Preferences())
	vpnChan := make(chan Command)
	uiChan := make(chan State)
	ui := NewUi(a, config, uiChan, vpnChan)
	menu, err := NewMenu(a)
	if err != nil {
		return err
	}
	vpn := NewVpn(a, ui, config.Scrollback, menu, uiChan, vpnChan)
	ui.SetVpn(vpn)
	ui.Run()
	return nil
}

func ReadPreferences(p fyne.Preferences) Config {
	return Config{
		Scrollback: p.Int("scrollback"),
		Connection: Connection{
			server:    p.String("server"),
			user:      p.String("user"),
			search:    p.String("search"),
			group:     p.String("group"),
			extraArgs: p.String("extraArgs"),
		},
	}
}

func WritePreferences(p fyne.Preferences, c Config) {
	p.SetInt("scrollback", c.Scrollback)
	conn := c.Connection
	p.SetString("server", conn.server)
	p.SetString("user", conn.user)
	p.SetString("search", conn.search)
	p.SetString("group", conn.group)
	p.SetString("extraArgs", conn.extraArgs)
}
