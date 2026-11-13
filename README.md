# aup

Update every installed agent harness (or any tool with an update command) in one run.

- Steps run in order from a TOML config. A step whose command is not installed is skipped.
- A failed step never stops the next one. The exit code is 1 when any step failed.
- On a terminal you get a live TUI: a step list with spinners, the running step's output, a progress bar, and a summary table. Off a terminal, under `CI`, with `TERM=dumb`, or with `--plain`, you get plain line output.
- Every run is logged, without colour codes, to `~/.local/state/aup/last.log`.

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/edheltzel/aup/master/install.sh | sh
# or
brew install --cask edheltzel/tap/aup
# or
go install github.com/edheltzel/aup/cmd/aup@latest
```

The script installs the latest release for macOS or Linux (amd64, arm64) to `~/.local/bin`, checks it against the release checksums, and writes a starter config to `~/.config/aup/config.toml` when you have none. An existing `config.toml`, legacy `harnesses.toml`, or `$AUP_CONFIG` is never touched. Set `AUP_VERSION=v0.1.0` to pin a release or `AUP_BIN_DIR` to install elsewhere. Run it again to update.

Uninstall (keeps your config and log):

```sh
curl -fsSL https://raw.githubusercontent.com/edheltzel/aup/master/install.sh | sh -s -- --uninstall
```

From a clone, with [just](https://github.com/casey/just):

```sh
just install     # build to ~/.local/bin/aup (override with AUP_BIN_DIR)
just update      # git pull, then reinstall
just uninstall   # remove the binary
```

The config in `~/.config/aup/` and the log in `~/.local/state/aup/` are left in place by `uninstall`.

## Use

```sh
aup                 # update everything
aup --plain         # plain output
aup --only codex    # one step, by label, any case
aup --dry-run       # show what would run
aup --config path   # another config file
```

TUI keys: `j`/`k` select a step, `f` follow the running step, `ctrl+c` cancel the running step, `q` stop and quit, `enter` close the summary.

## Config

Looked up in this order: `--config`, `$AUP_CONFIG`, `$XDG_CONFIG_HOME/aup/config.toml` (default `~/.config/aup/config.toml`), then the legacy `~/.config/aup/harnesses.toml`. See [examples/config.toml](examples/config.toml) for every field.

| Field | Required | Meaning |
| --- | --- | --- |
| `label` | yes | Name shown in the output |
| `color` | yes | Hex colour for the label |
| `binary` | yes | Command to run |
| `args` | no | String or list |
| `installed` | no | Command that proves the tool is installed; defaults to `binary` |
| `version` | no | Command probed with `--version` before and after; `-` skips |
| `interactive` | no | Give the step the real terminal (prompts, sudo) |
| `env` | no | Extra environment variables |
| `timeout` | no | Stop the step after a duration, like `"10m"` |

Top-level settings for the TUI summary:

| Setting | Default | Meaning |
| --- | --- | --- |
| `auto_close` | `false` | Close the TUI as soon as the run ends, failures included. The summary still prints to the terminal. |
| `close_after` | `"10s"` | After a run with no failures, close the summary on its own after this long. `"0s"` always waits for Enter. Any key stops the countdown. |

## Develop

```sh
just test                                        # vet and all tests, including a PTY end-to-end run
go test ./internal/tui -run TestViewGolden -update  # refresh TUI snapshots after a deliberate UI change
```

Layout: `internal/config` loads TOML, `internal/runner` plans and runs steps and never touches the screen, `internal/plain` and `internal/tui` render runner events, `internal/cli` picks the mode, `cmd/aup` is the entry point. Releases are cut by tagging `v*`; GoReleaser publishes binaries and the Homebrew cask.
