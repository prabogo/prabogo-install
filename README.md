# Prabogo Installer

This is a command-line tool for creating new projects based on the Prabogo framework.

## Installation

```sh
# Install the latest version
go install github.com/prabogo/prabogo-install@latest

# Or install a specific version
go install github.com/prabogo/prabogo-install@v0.1.0
```

If you encounter a module path error, you can try these solutions:

```sh
# Option 1: Bypass Go proxy (if the tag is very new)
GOPROXY=direct go install github.com/prabogo/prabogo-install@latest

# Option 2: Install from local source
git clone https://github.com/prabogo/prabogo-install.git
cd prabogo-install
go install .
```

## Usage

```sh
prabogo-install my-project-name
```

## Customization Options

Future versions may include flags for customizing your new project:

```sh
prabogo-install my-project-name --skip-compose --minimal
```