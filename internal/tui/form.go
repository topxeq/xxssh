package tui

import (
	"strconv"

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
		a.refreshMainView()
		a.pages.SwitchToPage("main")
	})

	form.AddButton("Cancel", func() {
		a.pages.SwitchToPage("main")
	})

	title := "Add Server"
	if idx >= 0 {
		title = "Edit Server"
	}
	form.SetBorder(true).SetTitle(title)
	a.pages.AddPage("form", form, true, true)
	a.pages.SwitchToPage("form")
}
