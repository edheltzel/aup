# aup

## Purpose

`aup` updates every installed agent harness (or any tool with an update command) in one run, with a live Bubble Tea TUI on a terminal and plain output everywhere else. Module `github.com/edheltzel/aup`.

## Ownership

- `cmd/aup` - entry point, version via `-ldflags "-X main.version=..."`.
- `internal/config` - TOML load, validation, defaults, config path lookup.
- `internal/runner` - plan and run steps, PATH setup, duplicate detection, version probe. Never imports a UI package.
- `internal/plain` - plain line output and the ANSI-stripping log writer.
- `internal/tui` - Bubble Tea v2 model, PTY bridge, line buffer, styles.
- `internal/cli` - Cobra + Fang flags, mode selection, exit codes.
- Ed's personal step list lives in his dotfiles (`config/.config/aup/harnesses.toml`), not here. `examples/config.toml` documents every field.

## Local Contracts

- Version control is GitButler (`but`). Work on a branch, push it, open a PR; `master` is the target. Never create native git worktrees.
- Agent commits use the Atlas-Key identity (`Atlas`, `296298943+Atlas-Key@users.noreply.github.com`, signing key `/Users/ed/.ssh/atlas_signing`), set repo-local for the commit and unset after.
- Plain mode output is a contract pinned by `internal/cli/cli_test.go`. Change it only on purpose and update those tests.
- The TUI starts only when stdin and stdout are terminals, `CI` is unset, `TERM` is not `dumb`, and neither `--plain` nor `--dry-run` (without `--tui`) is set.
- A failed step never stops later steps. Exit code 1 when any step failed, 2 on usage errors.
- The log at `~/.local/state/aup/last.log` never contains escape codes.
- Releases: tag `v*` on `master`; `.github/workflows/release.yml` runs GoReleaser, which needs the `HOMEBREW_TAP_TOKEN` secret to write `Casks/aup.rb` to `edheltzel/homebrew-tap`.
- Tasks live in `todo.txt` (tuxedo CLI).

## Work Guidance

- After a deliberate TUI change, refresh snapshots with `go test ./internal/tui -run TestViewGolden -update` and read the diff.
- `go install ./cmd/aup` or `go build -o ~/.local/bin/aup ./cmd/aup` for a local build.

## Verification

- `go vet ./... && go test -race ./...` (includes a PTY end-to-end run of the built binary).
- `goreleaser check` and `goreleaser release --snapshot --clean --skip=publish` before changing `.goreleaser.yaml`.

## Child DOX Index

None.
