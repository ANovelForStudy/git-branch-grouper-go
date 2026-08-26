# 🌿🎨 Git Branch Grouper

A fast, colorized CLI tool that displays git branches grouped by prefix with hierarchical navigation and powerful filtering.

### Core Principles

* **Native Git Access**: Powered by `go-git` - no external `git` binary required. Works anywhere Go compiles.
* **Hierarchical Branch View**: Branches like `feat/auth/login` render as a navigable tree, not a flat list.
* **Powerful Filtering**: Include or exclude branches at any depth with slash-separated paths (`backup/v2/refactor`).
* **Color-Coded Output**: Every group gets its own color. Supports named ANSI colors, 256-color, and hex (`#RRGGBB`) via `config.toml`.


### Features

| Feature | Description |
|---|---|
| **Prefix Grouping** | Branches are automatically grouped by their first path segment (`feat/`, `fix/`, etc.) |
| **Hierarchical Tree** | Deep paths like `feat/auth/login` display as a nested tree structure |
| **Include Filter** | Show only specific groups or sub-paths (`-i feat,backup/v1`) |
| **Exclude Filter** | Hide specific groups or sub-paths (`-e old,backup/v2`) |
| **Active Branch Star** | The current branch is marked with a `*` prefix |
| **Sparse Mode** | Add blank lines between groups for visual separation (`-s`) |
| **Custom Colors** | Named ANSI, 256-color, and hex (`#RRGGBB`) per group via TOML |
| **No-Color Mode** | `--no-color` flag or `NO_COLOR` env var disables all color output |
| **Custom Config** | `--config` flag for explicit config file path (CI/CD, alternate configs) |
| **Default Branch Detection** | `main`, `master`, `develop` are automatically categorized |


### Download

Pre-built binaries are available for **Windows**, **Linux**, and **macOS** on the [Releases](https://github.com/ANovelForStudy/GitBranchGrouper/releases) page.

| Platform | Architecture | Archive |
|---|---|---|
| Windows | x64 | `.zip` |
| Windows | ARM64 | `.zip` |
| Linux | x64 | `.tar.gz` |
| Linux | ARM64 | `.tar.gz` |
| macOS | Intel | `.tar.gz` |
| macOS | Apple Silicon | `.tar.gz` |

Download the archive for your system, extract the binary, and place it in your `PATH`:

```bash
# Linux / macOS
tar -xzf git-branch-grouper_*_linux_amd64.tar.gz
sudo mv git-branch-grouper /usr/local/bin/

# Windows (PowerShell)
Expand-Archive git-branch-grouper_*_windows_amd64.zip -DestinationPath .
Move-Item git-branch-grouper.exe C:\Windows\System32\
```


### Build from Source

Requires [Go 1.22+](https://go.dev/dl/) installed.

```bash
git clone https://github.com/ANovelForStudy/GitBranchGrouper.git
cd GitBranchGrouper
go build -o git-branch-grouper ./cmd/git-branch-grouper
./git-branch-grouper
```


### Usage

```
git-branch-grouper [flags]
```

#### Flags

| Flag | Short | Description |
|---|---|---|
| `--include` | `-i` | Comma-separated list of groups or sub-paths to show |
| `--exclude` | `-e` | Comma-separated list of groups or sub-paths to hide |
| `--sparse` | `-s` | Add blank line between groups |
| `--no-color` | `-n` | Disable colored output |
| `--config` | | Path to config file |
| `--version` | `-v` | Print version and exit |
| `--help` | `-h` | Show help message |

#### Environment Variables

| Variable | Description |
|---|---|
| `NO_COLOR` | Set to any non-empty value to disable color output ([no-color.org](https://no-color.org/)) |

#### Examples

Show only `feat` and `fix` groups:
```bash
git-branch-grouper --include feat,fix
```

Hide the `v2` subtree inside `backup`:
```bash
git-branch-grouper -e backup/v2
```

Show `feat` group and only `v1` inside `backup`:
```bash
git-branch-grouper -i backup/v1,feat
```

Display with blank lines between groups:
```bash
git-branch-grouper -s
```

Use a custom config file:
```bash
git-branch-grouper --config /path/to/config.toml
```

Disable colors:
```bash
git-branch-grouper --no-color
# or
NO_COLOR=1 git-branch-grouper
```


### Branch Groups

| Group | Description |
|---|---|
| `[default]` | Standard branches: `main`, `master`, `develop` |
| `[other]` | Branches without a recognized group prefix (hidden when empty) |
| `prefix/` | Hierarchical groups (`feat/`, `fix/`, `hotfix/`, etc.) |
| `sub/path` | Nested sub-paths within groups (`backup/v2`, `feat/auth/login`) |


### Configuration

Config is loaded from the first path found, or explicitly via `--config`:

1. `./config.toml` (current directory)
2. `$XDG_CONFIG_HOME/git-branch-grouper/config.toml`
3. `~/.config/git-branch-grouper/config.toml`
4. `<binary-dir>/config.toml`

#### Config File Reference

```toml
[main]
sparse = false

[format]
group_prefix = "[{group}]"
indent       = "  "
branch_marker = "*"

[colors]
default     = "#FF3333"
other       = "#FBBF24"
active_star = "#b7ff32"

[groups]
# STABLE
backup   = "#64748B"
default  = "#1E293B"

# DEVELOPMENT
feat     = "#10B981"
refactor = "#34D399"

# CRITICAL
fix      = "#EF4444"
hotfix   = "#DC2626"

# INFRASTRUCTURE
build    = "#3B82F6"
ci       = "#22D3EE"
release  = "#0EA5E9"

# DOCUMENTATION & TESTS
docs     = "#818CF8"
test     = "#C084FC"

# UTILITY
chore    = "#9CA3AF"

# LEGACY & EXPERIMENTS
old      = "#F59E0B"
exp      = "#F97316"
```

#### Color Formats

Colors accept three formats:

| Format | Example | Description |
|---|---|---|
| Named ANSI | `red`, `hi-green`, `cyan` | 16 standard terminal colors |
| 256-color | `196`, `226` | ANSI 256-color palette (0-255) |
| Hex RGB | `#FF3333`, `#f00` | TrueColor hex (6-digit or 3-digit shorthand) |

Named ANSI colors: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white` and their `hi-` variants.


### Project Structure

```
.
├── cmd/
│   └── git-branch-grouper/
│       └── main.go                  # Entry point, CLI flags, usage text
├── internal/
│   ├── config/
│   │   ├── config.go                # Config loading and TOML parsing
│   │   └── config_test.go           # Config unit tests
│   ├── model/
│   │   └── model.go                 # Domain types: Node, BranchData
│   ├── git/
│   │   └── git.go                   # Repository operations, branch collection
│   ├── filter/
│   │   ├── filter.go                # Branch filtering and tree manipulation
│   │   └── filter_test.go           # Filter unit tests
│   └── display/
│       ├── display.go               # Terminal output and color rendering
│       └── display_test.go          # Display unit tests
├── config.toml                      # Default color and group configuration
├── justfile                         # Development recipes (build, test, lint, release)
├── go.mod
└── go.sum
```


### Development

The project uses [just](https://github.com/casey/just) as a task runner.

```bash
just --list              # Show all available recipes
```

#### Common Recipes

| Recipe | Description |
|---|---|
| `just build` | Build binary for current OS |
| `just run` | Build and run |
| `just test` | Run all tests |
| `just check` | Run format, vet, lint, and test |
| `just fmt` | Format with gofumpt |
| `just lint` | Run golangci-lint |
| `just patch` | Bump patch version, create tag |
| `just minor` | Bump minor version, create tag |
| `just major` | Bump major version, create tag |
| `just push-tags` | Push all tags to remote |
| `just release-snapshot` | GoReleaser snapshot (local) |
| `just tidy` | Update Go module dependencies |


### How It Works

1. Opens the git repository in the current working directory
2. Collects all branches and builds a hierarchical tree based on `/` separators
3. Applies include/exclude filters if specified via flags
4. Renders the tree to the terminal with color-coded groups and active branch highlighting


### Running Tests

```bash
go test ./...
```


### Dependencies

| Package | Purpose |
|---|---|
| [`go-git`](https://github.com/go-git/go-git) | Native Go git implementation |
| [`fatih/color`](https://github.com/fatih/color) | Terminal color output (ANSI, 256-color, TrueColor) |
| [`go-toml`](https://github.com/pelletier/go-toml) | TOML config file parsing |
