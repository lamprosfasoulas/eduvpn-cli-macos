package oauth

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"codeberg.org/eduVPN/eduvpn-common/client"
	"codeberg.org/eduVPN/eduvpn-common/i18n"
	discotypes "codeberg.org/eduVPN/eduvpn-common/types/discovery"
	srvtypes "codeberg.org/eduVPN/eduvpn-common/types/server"

	"github.com/pkg/browser"
)

// openBrowser opens the OAuth URL carried in data in the default browser.
// The URL is always printed, since the browser may open in the background
// (or there may be no GUI session at all), leaving the terminal looking
// stuck while it actually just waits for you to complete login.
func OpenBrowser(data any) {
	str, ok := data.(string)
	if !ok {
		return
	}
	fmt.Println("Waiting for you to log in. If a browser did not open, visit this URL:")
	fmt.Println(str)
	go func() {
		if err := browser.OpenURL(str); err != nil {
			fmt.Fprintln(os.Stderr, "failed to open browser automatically:", err)
		}
	}()
}

// getProfileInteractive asks the user to pick a profile from stdin.
func GetProfileInteractive(profiles *srvtypes.Profiles) (string, error) {
	fmt.Printf("Multiple VPN profiles found. Please select a profile by entering e.g. 1")
	var ps strings.Builder
	var options []string
	i := 0
	for k, v := range profiles.Map {
		fmt.Fprintf(&ps, "\n%d - %s", i+1, i18n.GetLanguageMatched(v.DisplayName, "en"))
		options = append(options, k)
		i++
	}
	fmt.Println(ps.String())

	var idx int
	if _, err := fmt.Scanf("%d", &idx); err != nil || idx <= 0 || idx > len(profiles.Map) {
		fmt.Fprintln(os.Stderr, "invalid profile chosen, please retry")
		return GetProfileInteractive(profiles)
	}

	p := options[idx-1]
	fmt.Println("Sending profile ID", p)
	return p, nil
}

// SelectInstitute asks the user to pick one of several matching institute
// access servers from stdin, returning the chosen server's base URL.
func SelectInstitute(servers []discotypes.Server) (string, error) {
	fmt.Println("Multiple institutions matched. Please select one by entering e.g. 1")
	for i, s := range servers {
		fmt.Printf("%d - %s\n", i+1, i18n.GetLanguageMatched(s.DisplayName, "en"))
	}

	var idx int
	if _, err := fmt.Scanf("%d", &idx); err != nil || idx <= 0 || idx > len(servers) {
		fmt.Fprintln(os.Stderr, "invalid institution chosen, please retry")
		return SelectInstitute(servers)
	}
	return servers[idx-1].BaseURL, nil
}

// sendProfile replies to a StateAskProfile transition, using profile if
// given, otherwise falling back to an interactive picker.
func SendProfile(profile string, data any) {
	d, ok := data.(*srvtypes.RequiredAskTransition)
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid data type: %v\n", reflect.TypeOf(data))
		os.Exit(1)
	}
	sps, ok := d.Data.(*srvtypes.Profiles)
	if !ok {
		fmt.Fprintf(os.Stderr, "invalid data type for profiles: %v\n", reflect.TypeOf(d.Data))
		os.Exit(1)
	}

	if profile == "" {
		gprof, err := GetProfileInteractive(sps)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed getting profile interactively: %v\n", err)
			os.Exit(1)
		}
		profile = gprof
	}
	if err := d.C.Send(profile); err != nil {
		fmt.Fprintf(os.Stderr, "failed setting profile with error: %v\n", err)
		os.Exit(1)
	}
}

// stateCallback drives the interactive parts of the FSM: opening a browser
// for OAuth and answering profile/location questions.
//
// country is only used to decide whether we can give the user a useful error
// message if the server unexpectedly asks for a secure internet location;
// SetSecureLocation is called proactively before GetConfig in doConnect, so
// this branch should not normally be reached.
func StateCallback(_, newState client.FSMStateID, data any, profile, country string) {
	switch newState {
	case client.StateOAuthStarted:
		OpenBrowser(data)
	case client.StateAskProfile:
		SendProfile(profile, data)
	case client.StateAskLocation:
		if country == "" {
			fmt.Fprintln(os.Stderr, "the server requires a secure internet location; pass -country <cc>")
			os.Exit(1)
		}
	}
}

// noopCallback is used for subcommands (disconnect, status) that never
// trigger interactive FSM transitions.
func NoopCallback(_, _ client.FSMStateID, _ any) bool {
	return true
}
