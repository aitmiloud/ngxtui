# NGINX Terminal UI Manager (NgxTUI)

A modern, feature-rich terminal UI for managing NGINX servers, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and Go.

<p align="center">
  <img src="assets/header.png" alt="NGX TUI Logo" width="500"/>
</p>

<p align="center">
  <a href="https://github.com/aitmiloud/ngxtui/actions"><img src="https://img.shields.io/github/actions/workflow/status/aitmiloud/ngxtui/ci.yml?branch=main" alt="Build Status"></a>
  <a href="https://github.com/aitmiloud/ngxtui/releases"><img src="https://img.shields.io/github/v/release/aitmiloud/ngxtui" alt="Latest Release"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License: MIT"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-%3E%3D%201.24-00ADD8" alt="Go Version"></a>
</p>

## Features

- Interactive, keyboard-driven dashboard: Sites, Logs, Stats, and Metrics tabs
- Site management: enable/disable, config test, graceful reload, quick add
- Powerful template system for "Add Site" with 11 pre-configured templates:
  - Static, SPA, Node.js, WordPress, Laravel, Django, Docker/Proxy, WebSocket, Domain Redirect, API Gateway, Blank
- Auto-populated forms, validation, and quick-start guides per template
- NGINX config generation with SSL/TLS, proxy, PHP-FPM, and custom block support
- Real-time metrics: CPU, memory, network, request throughput with live charts
- Access log viewer: color-coded by status code, auto-scroll, quick filtering
- Statistics: site distribution and performance summaries
- Modular architecture: testable, production-ready, easy to extend
- Modern TUI styling with a clean dark theme

## Installation & Build

Prerequisites:
- Go 1.24 or higher
- NGINX installed and running
- A Unix-like OS (Linux recommended)

Build from source:

```bash
git clone https://github.com/aitmiloud/ngxtui.git
cd ngxtui
make build
# or
go build -o bin/ngxtui ./cmd/ngxtui
```

## Usage

Run the application (requires elevated privileges to manage NGINX configs):

```bash
sudo ./bin/ngxtui        # after building locally
# or
sudo make run            # build and run in one step
```

Notes:
- Sudo is required to read/write NGINX configuration files and reload the service.
- NgxTUI does not modify your configs without explicit actions from you.

### Keyboard Controls

- `←/→` or `h/l`: Switch between tabs
- `↑/↓` or `k/j`: Navigate items
- `Enter`: Select/Execute action
- `Esc`: Go back
- `a`: Add site
- `r`: Refresh sites
- `q`: Quit application

## Tabs Overview

### Sites Tab
- View all configured NGINX sites
- Enable/disable sites
- Test configuration
- Reload NGINX server
- Add new site configuration

### Logs Tab
- Real-time access log viewing
- Color-coded by status codes
- Auto-scroll support
- Quick status filtering

### Stats Tab
- Total sites overview
- Active sites count
- Request rate statistics
- Uptime monitoring

### Metrics Tab
- Real-time CPU usage
- Memory utilization
- Network traffic
- Request rate trends

## Template System

NgxTUI ships with a comprehensive template system for quick, safe site provisioning.

- 11 ready-to-use templates: Static, SPA, Node.js, WordPress, Laravel, Django, Docker/Proxy, WebSocket, Domain Redirect, API Gateway, Blank
- Auto-populates form fields and generates production-grade NGINX configs
- Includes SSL/TLS, proxying, PHP-FPM, and custom configuration blocks
- Built-in validation, previews, and quick-start guidance

## Requirements

- Go 1.24 or higher
- NGINX installed and configured
- Unix-like operating system (Linux, macOS)
- Terminal with true color support

## Development

This project follows a modular architecture with clear separation of concerns. See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed documentation.

### Project Structure
```
cmd/ngxtui/          # Application entry point
internal/
  ├── app/           # Bubble Tea application logic
  ├── model/         # Data models and types
  ├── nginx/         # NGINX service layer
  ├── styles/        # UI styling
  └── ui/            # View rendering
```

### Building
```bash
make build          # Build the application
make run            # Build and run (requires sudo)
make test           # Run tests
make lint           # Run linter
make fmt            # Format code
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the architecture patterns (see `ARCHITECTURE.md`)
4. Add tests for new functionality
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss): Style definitions
- [bubbles](https://github.com/charmbracelet/bubbles): UI components

## Troubleshooting

- Permission denied / operation not permitted
  - Run via `sudo` as shown above.
  - Ensure your user is allowed to reload NGINX (e.g., via `sudoers`).

- NGINX reload fails
  - Check syntax with `sudo nginx -t`.
  - Inspect error logs (e.g., `/var/log/nginx/error.log`).

- No sites listed / paths differ
  - Verify your NGINX `sites-available` and `sites-enabled` directories.
  - Customize paths in your system or adjust service configuration as needed.

- Colors or graphics look wrong
  - Use a terminal with true-color support.
  - Ensure `TERM` is set appropriately (e.g., `xterm-256color`).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for their amazing terminal UI libraries
- NGINX community for inspiration and documentation