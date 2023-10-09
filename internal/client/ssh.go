package client

import (
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func ignoreHostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

type SSHClient struct {
	Host                  string
	User                  string
	Password              string
	PKPath                string
	InsecureIgnoreHostKey bool

	conn *ssh.Client
}

func (c *SSHClient) Connect() error {
	var err error
	var authMethods []ssh.AuthMethod
	var hostkeyCallback ssh.HostKeyCallback

	if c.PKPath != "" {
		pkPath := c.PKPath
		// resolve home directory
		if strings.HasPrefix(pkPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			pkPath = strings.Replace(pkPath, "~/", homeDir+"/", 1)
		}
		pKey, err := os.ReadFile(pkPath)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey(pKey)
		if err != nil {
			return err
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if c.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.Password))
	}

	hostkeyCallback = ignoreHostKeyCallback
	if !c.InsecureIgnoreHostKey {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		knownHostsPath := fmt.Sprintf("%s/.ssh/known_hosts", homeDir)
		hostkeyCallback, err = knownhosts.New(knownHostsPath)
		if err != nil {
			return err
		}
	}

	conf := &ssh.ClientConfig{
		User:            c.User,
		HostKeyCallback: hostkeyCallback,
		Auth:            authMethods,
	}

	c.conn, err = ssh.Dial("tcp", c.Host, conf)
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) Disconnect() error {
	if c == nil {
		return nil
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (c *SSHClient) CreateControllerConfig(file string, cfg BlueChiControllerConfig) error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("echo \"%s\" > %s", cfg.Serialize(), BlueChiControllerConfdDirectory+file)
	_, err = session.Output(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) RemoveControllerConfig(file string) error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("rm %s", BlueChiControllerConfdDirectory+file)
	_, err = session.Output(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) RestartBlueChiController() error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.Output("systemctl start bluechi-controller")
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) StopBlueChiController() error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.Output("systemctl stop bluechi-controller")
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) CreateAgentConfig(file string, cfg BlueChiAgentConfig) error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("echo \"%s\" > %s", cfg.Serialize(), BlueChiAgentConfdDirectory+file)
	_, err = session.Output(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) RemoveAgentConfig(file string) error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("rm %s", BlueChiAgentConfdDirectory+file)
	_, err = session.Output(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) RestartBlueChiAgent() error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.Output("systemctl start bluechi-agent")
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) StopBlueChiAgent() error {
	if c == nil || c.conn == nil {
		return fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	_, err = session.Output("systemctl stop bluechi-agent")
	if err != nil {
		return err
	}

	return nil
}

func NewSSHClient(host string, user string, password string, pkPath string, insecureIgnoreHostKey bool) Client {
	return &SSHClient{
		Host:                  host,
		User:                  user,
		Password:              password,
		PKPath:                pkPath,
		InsecureIgnoreHostKey: insecureIgnoreHostKey,
	}
}

/*

 */

type SSHClientMock struct{}

func (c *SSHClientMock) Connect() error {
	return nil
}

func (c *SSHClientMock) Disconnect() error {
	return nil
}

func (c *SSHClientMock) CreateControllerConfig(file string, cfg BlueChiControllerConfig) error {
	return nil
}

func (c *SSHClientMock) RemoveControllerConfig(string) error {
	return nil
}

func (c *SSHClientMock) RestartBlueChiController() error {
	return nil
}

func (c *SSHClientMock) StopBlueChiController() error {
	return nil
}

func (c *SSHClientMock) CreateAgentConfig(file string, cfg BlueChiAgentConfig) error {
	return nil
}

func (c *SSHClientMock) RemoveAgentConfig(string) error {
	return nil
}

func (c *SSHClientMock) RestartBlueChiAgent() error {
	return nil
}

func (c *SSHClientMock) StopBlueChiAgent() error {
	return nil
}

func NewSSHClientMock() Client {
	return &SSHClientMock{}
}
