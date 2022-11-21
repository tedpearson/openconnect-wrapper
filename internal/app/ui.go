package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Button struct {
	values       []string
	currentValue int
	button       *widget.Button
}

type Ui struct {
	win     fyne.Window
	scroll  *container.Scroll
	grid    *widget.TextGrid
	connect *Button
	vpn     *Vpn
}

func NewUi() *Ui {
	a := app.New()
	win := a.NewWindow("OpenConnect Wrapper")
	entry := widget.NewTextGrid()
	ui := Ui{win: win, grid: entry}
	connect := widget.NewButton("Connect", ui.Connect)
	ui.connect = &Button{
		values: []string{"Connect", "Disconnect"},
		button: connect,
	}
	buttonLayout := container.NewHBox(connect)
	ui.scroll = container.NewVScroll(entry)
	ui.scroll.SetMinSize(fyne.NewSize(650, 500))
	border := container.NewBorder(nil, buttonLayout, nil, nil, ui.scroll)
	win.SetContent(border)
	return &ui
}

func (ui *Ui) SetVpn(vpn *Vpn) {
	ui.vpn = vpn
}

func (ui *Ui) Run() {
	ui.win.ShowAndRun()
}

func (ui *Ui) Connect() {
	println("connect")
	if ui.connect.currentValue == 0 {
		// connect
		go ui.vpn.Connect()
	} else {
		// disconnect
		go ui.vpn.Disconnect()
	}
	b := ui.connect
	b.currentValue = (b.currentValue + 1) % len(b.values)
	b.button.SetText(b.values[b.currentValue])
	b.button.Refresh()
}

func (ui *Ui) SetText(text string) {
	ui.grid.SetText(text)
	ui.scroll.ScrollToBottom()
}
