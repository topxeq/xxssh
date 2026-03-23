package tui

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/topxeq/xxssh/internal/config"
)

// showEditForm handles both add (-1 index) and edit (>=0 index) modes
func (a *App) showEditForm(idx int) {
	form := tview.NewForm()

	var name, host, username, password string
	var port int = 22

	// Pre-fill if editing existing server
	if idx >= 0 {
		cfg, err := a.store.Load()
		if err == nil && idx < len(cfg.Servers) {
			srv := cfg.Servers[idx]
			name = srv.Name
			host = srv.Host
			port = srv.Port
			username = srv.Username
			password = srv.Password
		}
	}

	form.AddInputField("Name", name, 40, nil, nil)
	form.AddInputField("Host", host, 40, nil, nil)
	form.AddInputField("Port", strconv.Itoa(port), 40, nil, nil)
	form.AddInputField("Username", username, 40, nil, nil)
	form.AddInputField("Password", password, 40, nil, nil)

	form.AddButton("Save", func() {
		name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
		host := form.GetFormItemByLabel("Host").(*tview.InputField).GetText()
		portStr := form.GetFormItemByLabel("Port").(*tview.InputField).GetText()
		port, _ := strconv.Atoi(portStr)
		if port == 0 {
			port = 22
		}
		username := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
		password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()

		srv := config.ServerConfig{
			Name: name, Host: host, Port: port,
			Username: username, Password: password,
		}

		cfg, err := a.store.Load()
		if err != nil {
			return
		}

		if idx < 0 {
			// Add new server
			cfg.Servers = append(cfg.Servers, srv)
		} else {
			// Update existing server
			cfg.Servers[idx] = srv
		}
		a.store.Save(cfg)
		a.pages.RemovePage("form")
		a.setupMainView()
	})

	form.AddButton("Cancel", func() {
		a.pages.RemovePage("form")
	})

	title := "Add Server"
	if idx >= 0 {
		title = "Edit Server"
	}
	form.SetBorder(true).SetTitle(title)

	// Wrap in Flex with black background to cover previous content
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)
	flex.SetBackgroundColor(tcell.ColorBlack)

	a.pages.AddPage("form", flex, true, true)
	a.app.SetRoot(a.pages, true)
	a.app.SetFocus(form)
}
