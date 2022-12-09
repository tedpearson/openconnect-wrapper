package app

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Command int

const (
	Connect Command = iota
	Disconnect
	DisconnectQuit
)

type State int

const (
	Disconnected State = iota
	Connected
	Connecting
	Disconnecting
)

type Ui struct {
	win     fyne.Window
	scroll  *container.Scroll
	grid    *widget.TextGrid
	connect *widget.Button
	vpn     *Vpn
	vpnChan chan Command
	uiChan  chan State
}

func NewUi(app fyne.App, config Config, uiChan chan State, vpnChan chan Command) *Ui {
	win := app.NewWindow("OpenConnect Wrapper")
	entry := widget.NewTextGrid()
	ui := Ui{win: win, grid: entry, uiChan: uiChan, vpnChan: vpnChan}
	ui.connect = widget.NewButton("Connect", func() { go ui.Connect() })
	prefWin := app.NewWindow("Preferences")
	prefWin.SetCloseIntercept(func() { prefWin.Hide() })
	prefs := widget.NewButton("Settings", func() { prefWin.Show() })
	buttonLayout := container.NewHBox(ui.connect, prefs)
	ui.scroll = container.NewVScroll(entry)
	ui.scroll.SetMinSize(fyne.NewSize(650, 500))
	border := container.NewBorder(nil, buttonLayout, nil, nil, ui.scroll)
	win.SetContent(border)
	// fixme: in case this breaks killing the spawned process
	win.SetOnClosed(func() { vpnChan <- DisconnectQuit })
	scrollback := widget.NewEntry()
	scrollback.SetText(strconv.Itoa(config.Scrollback))
	scrollback.SetPlaceHolder("99999")
	server := widget.NewEntry()
	server.SetPlaceHolder("vpn.example.com")
	conn := config.Connection
	server.SetText(conn.server)
	user := widget.NewEntry()
	user.SetText(conn.user)
	user.SetPlaceHolder("vpnuser")
	search := widget.NewEntry()
	search.SetPlaceHolder("Vpn Item In 1Password")
	search.SetText(conn.search)
	group := widget.NewEntry()
	group.SetText(conn.group)
	extraArgs := widget.NewEntry()
	extraArgs.SetText(conn.extraArgs)
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Scrollback Lines:", Widget: scrollback},
			{Text: "Server", Widget: server},
			{Text: "User:", Widget: user},
			{Text: "Search:", Widget: search},
			{Text: "Group:", Widget: group},
			{Text: "Extra Args:", Widget: extraArgs},
		},
		OnCancel: func() {
			prefWin.Hide()
		},
		OnSubmit: func() {
			sb, err := strconv.Atoi(scrollback.Text)
			if err != nil {
				sb = 500
			}
			WritePreferences(app.Preferences(), Config{
				Scrollback: sb,
				Connection: Connection{
					server:    server.Text,
					user:      user.Text,
					search:    search.Text,
					group:     group.Text,
					extraArgs: extraArgs.Text,
				},
			})
			prefWin.Hide()
		},
		SubmitText: "Save",
	}
	form.MinSize()
	prefWin.SetContent(form)
	prefWin.Resize(fyne.NewSize(450, 300))
	return &ui
}

func (ui *Ui) SetVpn(vpn *Vpn) {
	ui.vpn = vpn
}

func (ui *Ui) Run() {
	go ui.Listen()
	ui.win.ShowAndRun()
}

func (ui *Ui) Connect() {
	button := ui.connect
	if button.Text == "Connect" {
		ui.vpnChan <- Connect
		// fixme: cancel feature.
		button.SetText("Connecting...")
		button.Disable()
		//go ui.vpn.Connect()
	} else if button.Text == "Disconnect" {
		// disconnect
		ui.vpnChan <- Disconnect
		button.SetText("Disconnecting...")
		button.Disable()
		//go ui.vpn.Disconnect()
	}
	button.Refresh()
}

func (ui *Ui) Listen() {
	for {
		select {
		case state := <-ui.uiChan:
			button := ui.connect
			switch state {
			case Connected:
				button.SetText("Disconnect")
				button.Enable()
			case Disconnected:
				button.SetText("Connect")
				button.Enable()
			}
		}
	}
}

func (ui *Ui) SetText(text string) {
	ui.grid.SetText(text)
	ui.scroll.ScrollToBottom()
}
