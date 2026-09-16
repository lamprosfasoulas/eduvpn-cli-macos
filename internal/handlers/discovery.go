package handlers

import (
	"fmt"

	"codeberg.org/eduVPN/eduvpn-common/i18n"
	discotypes "codeberg.org/eduVPN/eduvpn-common/types/discovery"
	"github.com/lamprosfasoulas/eduvpn-cli-macos/internal/oauth"
)

// instituteTypeName is the discovery API's type tag for Institute Access
// servers, as opposed to Secure Internet.
const instituteTypeName = "institute_access"

// findInstitutes searches discovered Institute Access servers matching
// query, or returns all of them if query is empty.
func (c *Client) findInstitutes(query string) ([]discotypes.Server, error) {
	servers, err := c.eduClient.DiscoServers(c.cookie, false, query)
	if err != nil {
		return nil, fmt.Errorf("failed searching institute access servers: %w", err)
	}
	var out []discotypes.Server
	for _, s := range servers.List {
		if s.Type == instituteTypeName {
			out = append(out, s)
		}
	}
	return out, nil
}

// resolveInstitute finds the base URL identifier of the institute access
// server matching query: a single match is used directly, and multiple
// matches are resolved by asking the user to pick one interactively.
func (c *Client) resolveInstitute(query string) (string, error) {
	servers, err := c.findInstitutes(query)
	if err != nil {
		return "", err
	}
	switch len(servers) {
	case 0:
		return "", fmt.Errorf("no institute access server found matching %q (use --custom for a self-hosted server)", query)
	case 1:
		return servers[0].BaseURL, nil
	default:
		return oauth.SelectInstitute(servers)
	}
}

// Search prints the institute access servers matching query, or all of them
// if query is empty.
func (c *Client) Search(query string) error {
	servers, err := c.findInstitutes(query)
	if err != nil {
		return err
	}
	if len(servers) == 0 {
		fmt.Println("no institute access servers found")
		return nil
	}
	for _, s := range servers {
		fmt.Printf("%s  %s\n", s.BaseURL, i18n.GetLanguageMatched(s.DisplayName, "en"))
	}
	return nil
}
