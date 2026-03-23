package ssh

import (
	"fmt"
	"net"
	"time"

	"github.com/topxeq/xxssh/internal/config"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	config  *config.ServerConfig
	client  *ssh.Client
	session *ssh.Session
	done    chan struct{}
}

func NewSSHClient(cfg *config.ServerConfig) *SSHClient {
	return &SSHClient{config: cfg, done: make(chan struct{})}
}

func (c *SSHClient) Connect() error {
	addr := net.JoinHostPort(c.config.Host, fmt.Sprintf("%d", c.config.Port))

	auth := ssh.Password(c.config.Password)

	cfg := &ssh.ClientConfig{
		User:            c.config.Username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Enable TCP keepalive
	netConn, err := (&net.Dialer{
		Timeout:   cfg.Timeout,
		KeepAlive: 30 * time.Second,
	}).Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}

	sc, ch, req, err := ssh.NewClientConn(netConn, addr, cfg)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %w", err)
	}

	c.client = ssh.NewClient(sc, ch, req)

	// Start SSH heartbeat goroutine
	go c.heartbeat()

	return nil
}

// heartbeat sends SSH keepalive requests every 30 seconds
func (c *SSHClient) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if c.client != nil {
				// Send global keepalive request
				_, _, err := c.client.SendRequest("keepalive@topxeq/xxssh", true, nil)
				if err != nil {
					// Connection may be dead
					select {
					case <-c.done:
						return
					default:
					}
				}
			}
		case <-c.done:
			return
		}
	}
}

func (c *SSHClient) Session() (*ssh.Session, error) {
	if c.client == nil {
		return nil, fmt.Errorf("not connected")
	}
	return c.client.NewSession()
}

func (c *SSHClient) Close() {
	close(c.done)
	if c.session != nil {
		c.session.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
}
