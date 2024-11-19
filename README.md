<h1 align="center">
  <img src="./logo.png" alt="nmkill logo" width="40" height="40" style="vertical-align: middle; margin-right: 10px;">
  nmkill
</h1>

nmkill is a fast, efficient CLI tool for finding and removing `node_modules` directories, inspired by npkill but written in Go.

## Features

- Quickly scans directories for `node_modules` folders
- Interactive TUI for easy navigation and selection
- Displays folder sizes for informed decision-making
- Fast performance thanks to Go's concurrency

## Installation

Install nmkill with a single command:

```bash
curl -fsSL https://nmkill.sattwik.com/install.sh | bash
```

This script will download the latest version of nmkill and install it on your system.
## Usage
Run GoNpKill in your terminal:
```bash
nmkill
```

- Use arrow keys to navigate the list of node_modules directories
- Press Enter to select a directory for deletion
- Press q or Ctrl+C to exit

## Building from Source
If you prefer to build from source:
- Ensure you have Go installed (version 1.23 or later)
- Clone this repository
- Run `go build -o nmkill`

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
- Inspired by the original npkill project
- Built with Bubble Tea, Bubbles, and Lipgloss
