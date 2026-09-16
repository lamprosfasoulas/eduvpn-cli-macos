package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// aliasesFileName holds servers the user has saved under a short alias, so
// they can be connected to later without retyping a URL or search query.
const aliasesFileName = "servers.json"

// serverAlias is a saved server: its resolved identifier (a base URL) and
// whether it's a custom/self-hosted server or an institute access one.
type serverAlias struct {
	Server string `json:"server"`
	Custom bool   `json:"custom"`
}

func aliasesPath(dir string) string {
	return filepath.Join(dir, aliasesFileName)
}

// loadAliases reads the saved server aliases from dir. A missing file just
// means no server has been saved yet.
func loadAliases(dir string) (map[string]serverAlias, error) {
	data, err := os.ReadFile(aliasesPath(dir))
	if os.IsNotExist(err) {
		return map[string]serverAlias{}, nil
	}
	if err != nil {
		return nil, err
	}
	aliases := map[string]serverAlias{}
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, err
	}
	return aliases, nil
}

// saveAliases writes aliases to dir.
func saveAliases(dir string, aliases map[string]serverAlias) error {
	data, err := json.Marshal(aliases)
	if err != nil {
		return err
	}
	return os.WriteFile(aliasesPath(dir), data, 0o600)
}

func aliasKind(custom bool) string {
	if custom {
		return "custom"
	}
	return "institute access"
}

// AddAlias saves server under alias. Unless custom is set, server is
// resolved to a canonical institute access base URL first (via the same
// discovery search Connect uses), so connecting via the alias later never
// needs to search again.
func (c *Client) AddAlias(alias, server string, custom bool) error {
	identifier := server
	if !custom {
		resolved, err := c.resolveInstitute(server)
		if err != nil {
			return err
		}
		identifier = resolved
	}

	aliases, err := loadAliases(c.stateDir)
	if err != nil {
		return err
	}
	aliases[alias] = serverAlias{Server: identifier, Custom: custom}
	if err := saveAliases(c.stateDir, aliases); err != nil {
		return err
	}

	fmt.Printf("saved %s as %q (%s)\n", identifier, alias, aliasKind(custom))
	return nil
}

// RemoveAlias deletes a saved server alias.
func (c *Client) RemoveAlias(alias string) error {
	aliases, err := loadAliases(c.stateDir)
	if err != nil {
		return err
	}
	if _, ok := aliases[alias]; !ok {
		return fmt.Errorf("no saved server named %q", alias)
	}
	delete(aliases, alias)
	return saveAliases(c.stateDir, aliases)
}

// ListAliases prints every saved server alias.
func (c *Client) ListAliases() error {
	aliases, err := loadAliases(c.stateDir)
	if err != nil {
		return err
	}
	if len(aliases) == 0 {
		fmt.Println("no saved servers")
		return nil
	}
	for alias, a := range aliases {
		fmt.Printf("%s  %s  (%s)\n", alias, a.Server, aliasKind(a.Custom))
	}
	return nil
}
