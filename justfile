bin_dir := env("AUP_BIN_DIR", home_directory() / ".local/bin")
version := `git describe --tags --always --dirty 2>/dev/null || echo dev`

default:
    @just --list

# build and install aup to $AUP_BIN_DIR (default ~/.local/bin)
install:
    mkdir -p "{{bin_dir}}"
    go build -trimpath -ldflags "-s -w -X main.version={{version}}" -o "{{bin_dir}}/aup" ./cmd/aup
    @echo "installed {{bin_dir}}/aup ({{version}})"

# remove the installed aup binary
uninstall:
    rm -f "{{bin_dir}}/aup"
    @echo "removed {{bin_dir}}/aup"

# pull master and reinstall
update:
    git pull --ff-only
    just install

# vet and run every test
test:
    go vet ./...
    go test -race ./...
