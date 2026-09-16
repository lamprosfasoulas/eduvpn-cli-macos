# eduvpn-cli-macos

A headless [eduVPN](https://eduvpn.org/) / Let's Connect! client for macOS. It
handles server discovery and OAuth authentication, then manages the resulting
WireGuard tunnel via `wg-quick`, without needing the official GUI client.

## Requirements

- macOS
- [WireGuard tools](https://www.wireguard.com/install/) (`wg`, `wg-quick`) installed and on `PATH`
- `sudo` access, since bringing the tunnel up/down and reading its status require root

## Install

```sh
make install
```

This builds the `eduvpn-cli` binary and installs it to `/usr/local/bin`. Override
the destination with `PREFIX`:

```sh
make install PREFIX=$HOME/.local/bin
```

To just build the binary locally, run `make build`. Run `make uninstall` to
remove an installed binary, or `make clean` to remove the local build output.

## Usage

Connect directly to an institute access server by name, or a self-hosted
server with `--custom`:

```sh
eduvpn-cli connect "My University"
eduvpn-cli connect https://vpn.example.org --custom
```

Save a server under a short alias so you don't have to search or retype a URL
every time:

```sh
eduvpn-cli add uni "My University"
eduvpn-cli add home https://vpn.example.org --custom
eduvpn-cli connect uni
```

Other commands:

```sh
eduvpn-cli search [query]   # list institute access servers, optionally filtered by name
eduvpn-cli servers          # list saved server aliases
eduvpn-cli remove <alias>   # remove a saved server alias
eduvpn-cli status           # show the status of the eduVPN tunnel
eduvpn-cli disconnect       # disconnect the active eduVPN tunnel
```

## Status

Early / work in progress. See [todo.md](todo.md) for planned work such as
release automation, code signing, and a Homebrew tap.

## License

[MIT](LICENSE)
