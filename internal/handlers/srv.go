package handlers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	i18nerr "codeberg.org/eduVPN/eduvpn-common/i18n/err"
	"codeberg.org/eduVPN/eduvpn-common/types/protocol"
	srvtypes "codeberg.org/eduVPN/eduvpn-common/types/server"
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/wireguard"
)

func (c *Client) Connect(server string, custom bool) error {
	identifier := server

	aliases, err := loadAliases(c.stateDir)
	if err != nil {
		return err
	}
	if a, ok := aliases[server]; ok {
		identifier, custom = a.Server, a.Custom
	} else if !custom {
		resolved, err := c.resolveInstitute(server)
		if err != nil {
			return err
		}
		identifier = resolved
	}

	srvType := srvtypes.TypeInstituteAccess
	if custom {
		srvType = srvtypes.TypeCustom
	}
	if err := c.addServer(identifier, srvType); err != nil {
		return err
	}
	cfg, err := c.eduClient.GetConfig(c.cookie, identifier, srvType, false, false)
	if err != nil {
		return fmt.Errorf("failed getting VPN config: %w", err)
	}
	c.eduClient.Deregister()
	if cfg.Protocol != protocol.WireGuard {
		return fmt.Errorf("unsupported protocol %v (this tool only supports WireGuard, not OpenVPN or WireGuard-over-proxy)", cfg.Protocol)
	}
	confPath := wireguard.WgConfigPath(c.stateDir)
	if err := wireguard.WriteWGConfig(confPath, cfg.VPNConfig); err != nil {
		return fmt.Errorf("failed writing WireGuard config: %w", err)
	}

	if err := wireguard.WgQuickUp(confPath); err != nil {
		return fmt.Errorf("wg-quick up failed: %w", err)
	}

	fmt.Printf("connected to %s\n", identifier)
	return nil
}

func (c *Client) Disconnect() error {
	wgConfigPath := wireguard.WgConfigPath(c.stateDir)
	if _, err := os.Stat(wgConfigPath); err != nil {
		return errors.New("no active eduvpn-headless connection found")
	}
	if err := wireguard.WgQuickDown(wgConfigPath); err != nil {
		fmt.Fprintf(os.Stderr, "warning: wg-quick down failed: %v\n", err)
	}

	if err := c.eduClient.Cleanup(c.cookie); err != nil {
		fmt.Fprintf(os.Stderr, "warning: server-side disconnect failed: %v\n", err)
	}
	c.eduClient.Deregister()

	_ = os.Remove(wgConfigPath) //nolint:errcheck

	fmt.Println("disconnected")
	return nil

}

func (c *Client) Status() error {
	confPath := wireguard.WgConfigPath(c.stateDir)
	if _, err := os.Stat(confPath); err != nil {
		fmt.Println("disconnected")
		return nil
	}

	out, err := wireguard.WgShow()
	if err != nil {
		return fmt.Errorf("failed getting WireGuard status: %w", err)
	}
	if out == "" {
		fmt.Println("disconnected")
		return nil
	}
	fmt.Print(out)
	return nil
}

func (c *Client) addServer(server string, srvType srvtypes.Type) error {
	if err := c.eduClient.AddServer(c.cookie, server, srvType, nil); err != nil {
		var target *i18nerr.Error
		if errors.As(err, &target) && strings.Contains(target.Error(), "already added") {
			return nil
		}
		return err
	}
	return nil
}
