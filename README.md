# Prabogo Installer

This is a command-line tool for creating new projects based on the Prabogo framework.

## Installation

```sh
go install github.com/prabogo/prabogo-install@latest
```

If you encounter a module path error, you can install from local source:

```sh
# Clone the repository
git clone https://github.com/prabogo/prabogo-install.git
cd prabogo-install

# Install from local directory
go install .
```

## Usage

```sh
prabogo-install my-project-name
```

This will:
1. Clone the Prabogo repository
2. Remove the .git directory to start fresh
3. Replace the module name in go.mod
4. Set up everything for your new project

## Customization Options

Future versions may include flags for customizing your new project:

```sh
prabogo-install my-project-name --skip-compose --minimal
```