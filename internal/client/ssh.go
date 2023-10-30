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

	conn        *ssh.Client
	connHasRoot bool
}

func (c *SSHClient) newSSHSession() (*ssh.Session, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	return c.conn.NewSession()
}

func (c *SSHClient) hasRootPrivileges() (bool, error) {
	session, err := c.newSSHSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	output, err := session.Output("whoami")
	if err != nil {
		return false, fmt.Errorf("failed to determine if root: (%s, %s)", err.Error(), string(output))
	}

	c.connHasRoot = (string(output) == "root")
	return c.connHasRoot, nil
}

func (c *SSHClient) isServiceInstalled(service string) (bool, error) {
	session, err := c.newSSHSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	output, err := session.Output(fmt.Sprintf("systemctl list-unit-files %s", service))
	if err != nil {
		if serr, ok := err.(*ssh.ExitError); ok && serr.ExitStatus() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("failed to list unit files: %s", string(output))
	}

	return strings.Contains(string(output), service), nil
}

func (c *SSHClient) determineOS() (string, error) {
	session, err := c.newSSHSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output("cat /etc/os-release | grep -w ID=")
	if err != nil {
		return "", fmt.Errorf("failed to determine os: %s", string(output))
	}

	os := strings.ReplaceAll(string(output), "ID=", "")
	os = strings.ReplaceAll(os, "\"", "")
	return strings.TrimSpace(os), nil
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

	_, err = c.hasRootPrivileges()
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

func (c *SSHClient) InstallBlueChi(installCtrl bool, installAgent bool) error {
	needsInstallCtrl := false
	needsInstallAgent := false

	if installCtrl {
		isInstalled, err := c.isServiceInstalled("bluechi-controller.service")
		if err != nil {
			return err
		}
		needsInstallCtrl = !isInstalled
	}
	if installAgent {
		isInstalled, err := c.isServiceInstalled("bluechi-agent.service")
		if err != nil {
			return err
		}
		needsInstallAgent = !isInstalled
	}

	if !needsInstallCtrl && !needsInstallAgent {
		return nil
	}

	os, err := c.determineOS()
	if err != nil {
		return err
	}

	if os == "autosd" || os == "centos" {
		packagesToInstall := ""
		if needsInstallCtrl {
			packagesToInstall += " bluechi-controller bluechi-ctl "
		}
		if needsInstallAgent {
			packagesToInstall += " bluechi-agent"
		}

		session, err := c.newSSHSession()
		if err != nil {
			return err
		}
		defer session.Close()

		sudoPrefix := ""
		if c.connHasRoot {
			sudoPrefix = "sudo"
		}
		output, err := session.Output(fmt.Sprintf("%s dnf install -y %s", sudoPrefix, packagesToInstall))
		if err != nil {
			return fmt.Errorf("failed to install packages '%s': %s", packagesToInstall, output)
		}
	}

	return nil
}

func (c *SSHClient) CreateControllerConfig(file string, cfg BlueChiControllerConfig) error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	cmd := fmt.Sprintf("%s bash -c 'echo \"%s\" > %s'", sudoPrefix, cfg.Serialize(), BlueChiControllerConfdDirectory+file)
	output, err := session.Output(cmd)
	if err != nil {
		return fmt.Errorf("failed to create controller config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RemoveControllerConfig(file string) error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	cmd := fmt.Sprintf("%s rm %s", sudoPrefix, BlueChiControllerConfdDirectory+file)
	output, err := session.Output(cmd)
	if err != nil {
		return fmt.Errorf("failed to remove controller config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RestartBlueChiController() error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	output, err := session.Output(fmt.Sprintf("%ssystemctl start bluechi-controller", sudoPrefix))
	if err != nil {
		return fmt.Errorf("failed to restart controller service: %s", string(output))
	}

	return nil
}

func (c *SSHClient) StopBlueChiController() error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	output, err := session.Output(fmt.Sprintf("%s systemctl stop bluechi-controller", sudoPrefix))
	if err != nil {
		return fmt.Errorf("failed to stop controller service: %s", string(output))
	}

	return nil
}

func (c *SSHClient) CreateAgentConfig(file string, cfg BlueChiAgentConfig) error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	cmd := fmt.Sprintf("%s bash -c 'echo \"%s\" > %s'", sudoPrefix, cfg.Serialize(), BlueChiAgentConfdDirectory+file)
	output, err := session.Output(cmd)
	if err != nil {
		return fmt.Errorf("failed to create agent config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RemoveAgentConfig(file string) error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	cmd := fmt.Sprintf("%s rm %s", sudoPrefix, BlueChiAgentConfdDirectory+file)
	output, err := session.Output(cmd)
	if err != nil {
		return fmt.Errorf("failed to remove agent config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RestartBlueChiAgent() error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	output, err := session.Output(fmt.Sprintf("%s systemctl start bluechi-agent", sudoPrefix))
	if err != nil {
		return fmt.Errorf("failed to restart agent service: %s", string(output))
	}

	return nil
}

func (c *SSHClient) StopBlueChiAgent() error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	sudoPrefix := ""
	if c.connHasRoot {
		sudoPrefix = "sudo"
	}
	output, err := session.Output(fmt.Sprintf("%s systemctl stop bluechi-agent", sudoPrefix))
	if err != nil {
		return fmt.Errorf("failed to stop agent service: %s", string(output))
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

func (c *SSHClientMock) InstallBlueChi(installCtrl bool, installAgent bool) error {
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
