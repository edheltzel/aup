#!/bin/sh
# Install aup from its GitHub release:
#   curl -fsSL https://raw.githubusercontent.com/edheltzel/aup/master/install.sh | sh
# Uninstall:
#   curl -fsSL https://raw.githubusercontent.com/edheltzel/aup/master/install.sh | sh -s -- --uninstall
#
# AUP_VERSION      release tag to install, like v0.1.0 (default: latest)
# AUP_BIN_DIR      where the binary goes (default: ~/.local/bin)
# AUP_RELEASE_URL  base URL holding the release files (default: the GitHub release for AUP_VERSION)
set -eu

repo="edheltzel/aup"
bin_dir="${AUP_BIN_DIR:-$HOME/.local/bin}"
config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/aup"

say() { printf 'aup: %s\n' "$*"; }
die() { printf 'aup: %s\n' "$*" >&2; exit 1; }

uninstall() {
	if [ -e "$bin_dir/aup" ]; then
		rm -f "$bin_dir/aup"
		say "removed $bin_dir/aup"
	else
		say "nothing to remove at $bin_dir/aup"
	fi
	say "kept your config in $config_dir and the log in ~/.local/state/aup; delete them by hand if you want them gone"
}

platform() {
	case "$(uname -s)" in
	Darwin) os=darwin ;;
	Linux) os=linux ;;
	*) die "unsupported OS $(uname -s); build from source with: go install github.com/$repo/cmd/aup@latest" ;;
	esac
	case "$(uname -m)" in
	x86_64 | amd64) arch=amd64 ;;
	arm64 | aarch64) arch=arm64 ;;
	*) die "unsupported CPU $(uname -m)" ;;
	esac
}

latest() {
	url=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "https://github.com/$repo/releases/latest") || die "cannot reach GitHub"
	tag=${url##*/}
	case "$tag" in
	v*) printf '%s\n' "$tag" ;;
	*) die "no published release yet" ;;
	esac
}

sha256() {
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$1" | cut -d' ' -f1
	else
		shasum -a 256 "$1" | cut -d' ' -f1
	fi
}

seed_config() {
	if [ -n "${AUP_CONFIG:-}" ] || [ -e "$config_dir/config.toml" ] || [ -e "$config_dir/harnesses.toml" ]; then
		say "kept your existing config"
		return
	fi
	mkdir -p "$config_dir"
	cp "$1" "$config_dir/config.toml"
	say "wrote a starter config to $config_dir/config.toml; edit it to list your tools"
}

install() {
	command -v curl >/dev/null 2>&1 || die "curl is required"
	platform
	tag="${AUP_VERSION:-$(latest)}"
	version=${tag#v}
	base="${AUP_RELEASE_URL:-https://github.com/$repo/releases/download/$tag}"
	archive="aup_${version}_${os}_${arch}.tar.gz"
	tmp=$(mktemp -d)
	trap 'rm -rf "$tmp"' EXIT INT TERM
	say "downloading $archive"
	curl -fsSL -o "$tmp/$archive" "$base/$archive" || die "download failed: $base/$archive"
	curl -fsSL -o "$tmp/checksums.txt" "$base/checksums.txt" || die "download failed: $base/checksums.txt"
	want=$(awk -v f="$archive" '$2 == f { print $1 }' "$tmp/checksums.txt")
	[ -n "$want" ] || die "$archive is not listed in checksums.txt"
	[ "$want" = "$(sha256 "$tmp/$archive")" ] || die "checksum mismatch for $archive"
	tar -xzf "$tmp/$archive" -C "$tmp"
	mkdir -p "$bin_dir"
	cp "$tmp/aup" "$bin_dir/aup.new"
	chmod 755 "$bin_dir/aup.new"
	mv -f "$bin_dir/aup.new" "$bin_dir/aup"
	say "installed $bin_dir/aup ($("$bin_dir/aup" --version 2>/dev/null || printf '%s' "$tag"))"
	seed_config "$tmp/examples/config.toml"
	case ":$PATH:" in
	*":$bin_dir:"*) ;;
	*) say "$bin_dir is not on your PATH; add it in your shell config, then open a new terminal" ;;
	esac
	found=$(command -v aup 2>/dev/null || true)
	if [ -n "$found" ] && [ "$found" != "$bin_dir/aup" ]; then
		say "warning: 'aup' currently resolves to $found, not $bin_dir/aup"
	fi
}

case "${1:-}" in
"") install ;;
--uninstall) uninstall ;;
*) die "unknown option $1 (use --uninstall)" ;;
esac
