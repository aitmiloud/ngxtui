# NGINX Terminal UI Manager (NgxTUI)

A modern, feature-rich terminal user interface for managing NGINX servers, built with [bubbletea](https://github.com/charmbracelet/bubbletea) and Go.

<p align="center">
  <img src="assets/header.png" alt="NGX TUI Logo" width="500"/>
</p>

## Features

- **Interactive Dashboard**: Multiple views including Sites, Logs, Stats, and Metrics
- **Site Management**: Enable/disable sites, test configuration, and reload NGINX
- **Real-time Metrics**: Monitor CPU, Memory, Network, and Request metrics with live charts
- **Access Log Viewer**: Color-coded log viewing with status code highlighting
- **Statistics**: Visual representation of site distribution and performance metrics
- **Professional UI**: Clean, modern dark theme with consistent styling
- **Keyboard Driven**: Efficient keyboard shortcuts for all operations

## Installation

```bash
go install github.com/aitmiloud/ngxtui@latest
```

Or build from source:

```bash
git clone https://github.com/aitmiloud/ngxtui.git
cd ngxtui
make build
# or
go build -o bin/ngxtui ./cmd/ngxtui
```

## Usage

Run the application:

```bash
sudo ngxtui
```

> Note: Sudo privileges are required for managing NGINX configuration files.

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for their amazing terminal UI libraries
- NGINX community for inspiration and documentation