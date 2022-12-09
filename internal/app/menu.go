package app

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

type Menu struct {
	menu    *fyne.Menu
	desktop desktop.App
}

func NewMenu(app fyne.App) (*Menu, error) {
	if desk, ok := app.(desktop.App); ok {
		menu := fyne.NewMenu("VPN")
		m := &Menu{menu: menu, desktop: desk}
		m.SetConnected(Disconnected)
		desk.SetSystemTrayMenu(menu)
		return m, nil
	}
	return nil, errors.New("not a desktop app")
}

func (m *Menu) SetConnected(s State) {
	var icon fyne.Resource
	switch s {
	case Disconnected:
		icon = resourceRedCirclePng
	case Connected:
		icon = resourceGreenCirclePng
	default:
		icon = resourceYellowCirclePng
	}
	m.desktop.SetSystemTrayIcon(icon)
	m.menu.Refresh()
}
