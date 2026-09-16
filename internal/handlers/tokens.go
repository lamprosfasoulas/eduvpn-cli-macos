package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	srvtypes "codeberg.org/eduVPN/eduvpn-common/types/server"
)

// tokensFileName holds OAuth tokens cached across CLI invocations.
// eduvpn-common only keeps tokens in memory itself; it expects the host
// application to persist them via Client.TokenGetter/TokenSetter, which
// NewClient wires up to this file.
const tokensFileName = "tokens.json"

func tokensPath(dir string) string {
	return filepath.Join(dir, tokensFileName)
}

func tokenKey(id string, t srvtypes.Type) string {
	return fmt.Sprintf("%d:%s", t, id)
}

// loadTokens reads the cached tokens for every server from dir. A missing
// file just means no server has been authorized yet.
func loadTokens(dir string) (map[string]srvtypes.Tokens, error) {
	data, err := os.ReadFile(tokensPath(dir))
	if os.IsNotExist(err) {
		return map[string]srvtypes.Tokens{}, nil
	}
	if err != nil {
		return nil, err
	}
	toks := map[string]srvtypes.Tokens{}
	if err := json.Unmarshal(data, &toks); err != nil {
		return nil, err
	}
	return toks, nil
}

// saveTokens writes toks to dir with 0600 permissions, since it contains
// OAuth access and refresh tokens.
func saveTokens(dir string, toks map[string]srvtypes.Tokens) error {
	data, err := json.Marshal(toks)
	if err != nil {
		return err
	}
	return os.WriteFile(tokensPath(dir), data, 0o600)
}
