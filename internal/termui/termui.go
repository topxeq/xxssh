package termui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/topxeq/xxssh/internal/config"
	"github.com/topxeq/xxssh/internal/ssh"
	"github.com/topxeq/xxssh/internal/store"
	cssh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// ANSI escape codes
const (
	clearLine   = "\r\033[K"
	clearScreen = "\033[2J\033[H"
	moveToTop   = "\033[H"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
	bold        = "\033[1m"
	reset       = "\033[0m"
	cyan        = "\033[36m"
	green       = "\033[32m"
	yellow      = "\033[33m"
	red         = "\033[31m"
	dim         = "\033[2m"
)

type UI struct {
	store  *store.Store
	reader *bufio.Reader
}

func NewUI(s *store.Store) *UI {
	return &UI{
		store:  s,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (ui *UI) Run() error {
	// Set terminal to raw mode
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	// Handle ctrl+c
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	defer signal.Stop(sigChan)

	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	for {
		select {
		case <-sigChan:
			fmt.Print(clearScreen)
			fmt.Print(moveToTop)
			fmt.Println("Goodbye!")
			return nil
		default:
		}

		if err := ui.drawMainMenu(); err != nil {
			return err
		}

		key, err := ui.readKey()
		if err != nil {
			return err
		}

		switch key {
		case "q", "Q", "ctrl-c":
			fmt.Print(clearScreen)
			fmt.Print(moveToTop)
			fmt.Println("Goodbye!")
			return nil
		case "a", "A":
			if err := ui.addServer(); err != nil {
				ui.showMessage("Error: "+err.Error(), 3)
			}
		case "e", "E":
			ui.editServer()
		case "d", "D":
			ui.deleteServer()
		case "up":
			ui.moveSelection(-1)
		case "down":
			ui.moveSelection(1)
		case "enter":
			ui.connectSelected()
		default:
			// Number keys for quick connect
			if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
				ui.connectByNumber(int(key[0] - '1'))
			}
		}
	}
}

var selectedIndex = 0

func (ui *UI) drawMainMenu() error {
	cfg, err := ui.store.Load()
	if err != nil {
		cfg = &config.StoresConfig{}
	}

	fmt.Print(clearScreen)
	fmt.Print(moveToTop)

	// Draw header
	fmt.Printf("%s╔══════════════════════════════════════════════════╗%s\n", cyan, reset)
	fmt.Printf("%s║%s           xxssh - SSH Client                   %s║%s\n", cyan, reset, cyan, reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════╝%s\n\n", cyan, reset)

	// Lock warning
	if ui.store.IsLocked() && !ui.store.IsUnlocked() {
		fmt.Printf("%s⚠ Master password required. Use --reset to clear.%s\n\n", yellow, reset)
	}

	// Server list
	if len(cfg.Servers) == 0 {
		fmt.Printf("  %sNo servers configured.%s\n", dim, reset)
		fmt.Printf("  Press %sa%s to add a new server.\n\n", green, reset)
	} else {
		fmt.Printf("  %sServers:%s\n", bold, reset)
		for i, srv := range cfg.Servers {
			prefix := "  "
			marker := "  "
			if i == selectedIndex {
				prefix = ">"
				marker = "*"
			}

			nameStr := srv.Name
			if len(nameStr) > 20 {
				nameStr = nameStr[:17] + "..."
			}

			fmt.Printf("%s%s%d.%s %s%-20s%s %s%s:%d%s\n",
				prefix, green, i+1, reset,
				bold, nameStr, reset,
				dim, srv.Host, srv.Port, reset)
			fmt.Printf("%s  %sUser: %s%s\n", marker, dim, srv.Username, reset)
		}
		fmt.Println()
	}

	// Help
	fmt.Printf("%s─── Commands ───%s\n", bold, reset)
	fmt.Printf("  %s↑↓%s   Navigate   %sEnter%s  Connect\n", green, reset, green, reset)
	fmt.Printf("  %sa%s     Add       %se%s     Edit\n", green, reset, green, reset)
	fmt.Printf("  %sd%s     Delete    %sq%s     Quit\n\n", green, reset, green, reset)

	return nil
}

func (ui *UI) readKey() (string, error) {
	// Read 3 bytes for possible escape sequence
	b := make([]byte, 3)
	n, err := os.Stdin.Read(b)
	if err != nil {
		return "", err
	}

	if n == 1 {
		switch b[0] {
		case '\r', '\n':
			return "enter", nil
		case 3: // Ctrl+C
			return "ctrl-c", nil
		case 'q', 'Q', 'a', 'A', 'e', 'E', 'd', 'D':
			return string(b), nil
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return string(b), nil
		}
	}

	if n == 3 && b[0] == 27 && b[1] == '[' {
		switch b[2] {
		case 'A':
			return "up", nil
		case 'B':
			return "down", nil
		}
	}

	return string(b), nil
}

func (ui *UI) moveSelection(delta int) {
	cfg, _ := ui.store.Load()
	if len(cfg.Servers) == 0 {
		return
	}

	selectedIndex += delta
	if selectedIndex < 0 {
		selectedIndex = len(cfg.Servers) - 1
	}
	if selectedIndex >= len(cfg.Servers) {
		selectedIndex = 0
	}
}

func (ui *UI) connectByNumber(num int) {
	cfg, _ := ui.store.Load()
	if num >= 0 && num < len(cfg.Servers) {
		selectedIndex = num
		ui.connectSelected()
	}
}

func (ui *UI) connectSelected() {
	cfg, err := ui.store.Load()
	if err != nil || len(cfg.Servers) == 0 || selectedIndex >= len(cfg.Servers) {
		return
	}

	srv := cfg.Servers[selectedIndex]
	ui.connectToServer(&srv)
}

func (ui *UI) connectToServer(srv *config.ServerConfig) {
	// Clear screen before connection
	fmt.Print(clearScreen)
	fmt.Print(moveToTop)
	fmt.Printf("Connecting to %s%s%s (%s:%d)...\n\n", bold, srv.Name, reset, srv.Host, srv.Port)

	client := ssh.NewSSHClient(srv)

	if err := client.Connect(); err != nil {
		ui.showMessage("Connection failed: "+err.Error(), 3)
		return
	}

	defer client.Close()

	session, err := client.Session()
	if err != nil {
		ui.showMessage("Session failed: "+err.Error(), 3)
		return
	}
	defer session.Close()

	// Set up terminal
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	modes := cssh.TerminalModes{}
	modes[cssh.ECHO] = 1
	session.RequestPty("xterm-256color", width, height, modes)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		ui.showMessage("Shell failed: "+err.Error(), 3)
		return
	}

	session.Wait()
}

func (ui *UI) addServer() error {
	fmt.Print(clearScreen)
	fmt.Print(moveToTop)

	fmt.Printf("%s=== Add New Server ===%s\n\n", bold+cyan, reset)

	name := ui.readLine("  Name: ")
	if name == "" {
		return fmt.Errorf("name is required")
	}

	host := ui.readLine("  Host: ")
	if host == "" {
		return fmt.Errorf("host is required")
	}

	portStr := ui.readLine("  Port (default 22): ")
	port := 22
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	username := ui.readLine("  Username: ")
	if username == "" {
		return fmt.Errorf("username is required")
	}

	password := ui.readPassword("  Password: ")

	srv := config.ServerConfig{
		Name:     name,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}

	cfg, err := ui.store.Load()
	if err != nil {
		return err
	}

	cfg.Servers = append(cfg.Servers, srv)
	if err := ui.store.Save(cfg); err != nil {
		return err
	}

	ui.showMessage("Server added successfully!", 2)
	return nil
}

func (ui *UI) editServer() {
	cfg, err := ui.store.Load()
	if err != nil || len(cfg.Servers) == 0 || selectedIndex >= len(cfg.Servers) {
		return
	}

	srv := &cfg.Servers[selectedIndex]

	fmt.Print(clearScreen)
	fmt.Print(moveToTop)

	fmt.Printf("%s=== Edit Server ===%s\n\n", bold+cyan, reset)

	name := ui.readLineWithDefault("  Name: ", srv.Name)
	host := ui.readLineWithDefault("  Host: ", srv.Host)
	portStr := ui.readLineWithDefault("  Port: ", strconv.Itoa(srv.Port))
	username := ui.readLineWithDefault("  Username: ", srv.Username)
	password := srv.Password
	newPass := ui.readLineWithDefault("  Password (leave empty to keep): ", "")
	if newPass != "" {
		password = newPass
	}

	srv.Name = name
	srv.Host = host
	srv.Port, _ = strconv.Atoi(portStr)
	srv.Username = username
	srv.Password = password

	if err := ui.store.Save(cfg); err != nil {
		ui.showMessage("Failed to save: "+err.Error(), 3)
		return
	}

	ui.showMessage("Server updated!", 2)
}

func (ui *UI) deleteServer() {
	cfg, err := ui.store.Load()
	if err != nil || len(cfg.Servers) == 0 || selectedIndex >= len(cfg.Servers) {
		return
	}

	srv := cfg.Servers[selectedIndex]

	fmt.Print(clearScreen)
	fmt.Print(moveToTop)

	fmt.Printf("%s=== Delete Server ===%s\n\n", bold+red, reset)
	fmt.Printf("  Delete %s%s%s (%s)?\n\n", bold, srv.Name, reset, srv.Host)
	fmt.Printf("  %sy%s = Yes, %sn%s = No\n\n", green, reset, green, reset)

	key, _ := ui.readKey()
	if key == "y" || key == "Y" {
		cfg.Servers = append(cfg.Servers[:selectedIndex], cfg.Servers[selectedIndex+1:]...)
		ui.store.Save(cfg)
		if selectedIndex > 0 {
			selectedIndex--
		}
		ui.showMessage("Server deleted.", 2)
	}
}

func (ui *UI) readLine(prompt string) string {
	fmt.Printf("%s%s%s", bold, prompt, reset)
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func (ui *UI) readLineWithDefault(prompt, defaultVal string) string {
	fmt.Printf("%s%s [%s]: %s", bold, prompt, defaultVal, reset)
	line, _ := ui.reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

func (ui *UI) readPassword(prompt string) string {
	fmt.Printf("%s%s%s", bold, prompt, reset)
	// Simple read - password will be visible
	// For production, use term.ReadPassword
	line, _ := ui.reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func (ui *UI) showMessage(msg string, seconds int) {
	fmt.Printf("\n%s%s%s\n", bold+green, msg, reset)
	fmt.Printf("Press Enter to continue...")
	ui.reader.ReadString('\n')
}
