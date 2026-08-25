# 🗂️🪓 Git Branch Grouper

A fast, colorized CLI tool that displays git branches grouped by prefix with hierarchical navigation and powerful filtering.

### Core Principles

* **Native Git Access**: Powered by `go-git` — no external `git` binary required. Works anywhere Go compiles.
* **Hierarchical Branch View**: Branches like `feat/auth/login` render as a navigable tree, not a flat list.
* **Powerful Filtering**: Include or exclude branches at any depth with slash-separated paths (`backup/v2/refactor`).
* **Color-Coded Output**: Every group gets its own color. Customize everything through `config.toml`.


### Features

| Feature | Description |
|---|---|
| **Prefix Grouping** | Branches are automatically grouped by their first path segment (`feat/`, `fix/`, etc.) |
| **Hierarchical Tree** | Deep paths like `feat/auth/login` display as a nested tree structure |
| **Include Filter** | Show only specific groups or sub-paths (`-i feat,backup/v1`) |
| **Exclude Filter** | Hide specific groups or sub-paths (`-e old,backup/v2`) |
| **Active Branch Star** | The current branch is marked with a `*` prefix |
| **Sparse Mode** | Add blank lines between groups for visual separation (`-s`) |
| **Custom Colors** | Full color customization per group via TOML configuration |
| **Default Branch Detection** | `main`, `master`, `develop` are automatically categorized |


### Quick Start

```bash
git clone https://github.com/ANovelForStudy/GitBranchGrouper.git
cd GitBranchGrouper
go build -o git-branch-grouper .
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
| `--help` | `-h` | Show help message |

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


### Branch Groups

| Group | Description |
|---|---|
| `[default]` | Standard branches: `main`, `master`, `develop` |
| `[other]` | Branches without a recognized group prefix |
| `prefix/` | Hierarchical groups (`feat/`, `fix/`, `hotfix/`, etc.) |
| `sub/path` | Nested sub-paths within groups (`backup/v2`, `feat/auth/login`) |


### Configuration

Colors and group assignments are customized via `config.toml` in the working directory:

```toml
[main]
sparse = false

[colors]
default     = "red"
other       = "yellow"
active_star = "green"

[groups]
backup   = "blue"
build    = "magenta"
chore    = "white"
ci       = "cyan"
docs     = "blue"
exp      = "magenta"
feat     = "green"
fix      = "red"
hotfix   = "red"
old      = "yellow"
refactor = "cyan"
release  = "green"
test     = "blue"
```

#### Supported Colors

`red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`

`hi-red`, `hi-green`, `hi-yellow`, `hi-blue`, `hi-magenta`, `hi-cyan`, `hi-white`


### Project Structure

```
.
├── main.go           # Entry point, CLI flag parsing, config loading
├── model.go          # Domain types: Config, Node, BranchData
├── git.go            # Git repository operations and branch collection
├── filter.go         # Branch filtering and tree manipulation logic
├── display.go        # Terminal output and color rendering
├── main_test.go      # Unit tests for filter and utility functions
├── config.toml       # Default color and group configuration
├── go.mod            # Go module definition
└── go.sum            # Dependency checksums
```


### How It Works

1. Opens the git repository in the current working directory
2. Collects all branches and builds a hierarchical tree based on `/` separators
3. Applies include/exclude filters if specified via flags
4. Renders the tree to the terminal with color-coded groups and active branch highlighting


### Running Tests

```bash
go test -v .
```


### Dependencies

| Package | Purpose |
|---|---|
| [`go-git`](https://github.com/go-git/go-git) | Native Go git implementation |
| [`fatih/color`](https://github.com/fatih/color) | Terminal color output |
| [`go-toml`](https://github.com/pelletier/go-toml) | TOML config file parsing |
