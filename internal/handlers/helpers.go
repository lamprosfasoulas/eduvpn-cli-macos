package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"codeberg.org/eduVPN/eduvpn-common/client"
	"codeberg.org/eduVPN/eduvpn-common/types/cookie"
	srvtypes "codeberg.org/eduVPN/eduvpn-common/types/server"
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/oauth"
)

const appName = "eduvpn-cli-macos"
const clientID = "org.eduvpn.app.linux"
const commonver = "5.0.3"

type Client struct {
	eduClient *client.Client
	cookie    *cookie.Cookie
	stateDir  string
}

func NewClient() (*Client, error) {
	// State Directory
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(base, appName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	// Client Registration
	c, err := client.New(clientID, commonver, dir,
		func(_, newState client.FSMStateID, data any) bool {
			oauth.StateCallback(client.StateMain, newState, data, "", "")
			return true
		}, nil)
	if err != nil {
		return nil, err
	}

	// OAuth tokens: eduvpn-common keeps them in memory only, so cache them
	// to disk ourselves to avoid a fresh browser login on every invocation.
	tokens, err := loadTokens(dir)
	if err != nil {
		return nil, err
	}
	c.TokenGetter = func(sid string, t srvtypes.Type) *srvtypes.Tokens {
		tok, ok := tokens[tokenKey(sid, t)]
		if !ok {
			return nil
		}
		return &tok
	}
	c.TokenSetter = func(sid string, t srvtypes.Type, tok srvtypes.Tokens) {
		tokens[tokenKey(sid, t)] = tok
		if err := saveTokens(dir, tokens); err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed saving tokens:", err)
		}
	}

	if err := c.Register(); err != nil {
		return nil, err
	}

	// Cookie
	cookie := cookie.NewWithContext(context.Background())

	return &Client{
		stateDir:  dir,
		eduClient: c,
		cookie:    cookie,
	}, nil
}

func (c *Client) Cancel() error {
	return c.cookie.Cancel()
}
