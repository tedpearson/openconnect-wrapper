package app

import "flag"

type Config struct {
	Scrollback int
	Connection Connection
}

// todo:
//   handle untriggered disconnects and update ui state. use mvc type pattern?

func Main() error {
	config := ParseFlags()
	ui := NewUi()
	vpn := NewVpn(ui, config.Scrollback, config.Connection)
	ui.SetVpn(vpn)
	ui.Run()
	return nil
}

func ParseFlags() Config {
	server := flag.String("server", "", "VPN server")
	user := flag.String("user", "", "VPN username")
	search := flag.String("search", "", "Item to search 1Password for")
	group := flag.String("group", "", "VPN group")
	scrollback := flag.Int("scrollback", 500, "Number of lines of console to preserve")
	flag.Parse()
	return Config{
		Scrollback: *scrollback,
		Connection: Connection{
			server: *server,
			user:   *user,
			search: *search,
			group:  *group,
		},
	}
}
