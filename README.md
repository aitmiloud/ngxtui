# NGINX Terminal UI Manager (NgxTUI)

A modern, feature-rich terminal user interface for managing NGINX servers, built with [bubbletea](https://github.com/charmbracelet/bubbletea) and Go.

![NGINX Terminal UI Manager](https://placeholder-for-screenshot.png)

## Features

- 📊 **Interactive Dashboard**: Multiple views including Sites, Logs, Stats, and Metrics
- 🖥️ **Site Management**: Enable/disable sites, test configuration, and reload NGINX
- 📈 **Real-time Metrics**: Monitor CPU, Memory, Network, and Request metrics with live charts
- 📝 **Access Log Viewer**: Color-coded log viewing with status code highlighting
- 📊 **Statistics**: Visual representation of site distribution and performance metrics
- 🎨 **Professional UI**: Clean, modern dark theme with consistent styling
- ⌨️ **Keyboard Driven**: Efficient keyboard shortcuts for all operations

## Installation

```bash
go install github.com/aitmiloud/ngxtui@latest
```

Or build from source:

```bash
git clone https://github.com/aitmiloud/ngxtui.git
cd ngxtui
go build
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
- `r`: Refresh sites
- `q`: Quit application

## Tabs Overview

### Sites Tab
- View all configured NGINX sites
- Enable/disable sites
- Test configuration
- Reload NGINX server

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

To contribute to the project:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Dependencies

- [bubbletea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss): Style definitions
- [bubbles](https://github.com/charmbracelet/bubbles): UI components

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for their amazing terminal UI libraries
- NGINX community for inspiration and documentation

## Author

Mohamed AIT MILOUD - [@aitmiloud](https://github.com/aitmiloud)

---

Made with ❤️ using [Go](https://golang.org) and [Charm](https://charm.sh/)